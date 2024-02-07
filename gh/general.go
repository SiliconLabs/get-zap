/*
Copyright © 2024 Silicon Labs
*/
package gh

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type GithubConfiguration struct {
	Owner   string
	Repo    string
	Release string
	Token   string
	Asset   string
}

func CreateGithubClient(cfg *GithubConfiguration) *github.Client {
	var client *github.Client
	if cfg.Token == "" {
		fmt.Println("You do not have GET_ZAP_TOKEN set. This will limit the number of requests you can make to the github API.")
		client = github.NewClient(nil)
	} else {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}
	return client
}

func DefaultAction(cfg *GithubConfiguration) {
	fmt.Printf("Downloading release '%v' of repo '%v/%v' for the platform '%v/%v'...\n", cfg.Release, cfg.Owner, cfg.Repo, runtime.GOOS, runtime.GOARCH)
	DownloadAssets(cfg, true)
}

func DetermineAssetPlatform(assetName string) (os string, arch string) {
	os = ""
	arch = ""
	if strings.Contains(assetName, "-windows") || strings.Contains(assetName, "-win") {
		os = "windows"
	} else if strings.Contains(assetName, "-darwin") || strings.Contains(assetName, "-mac") || strings.Contains(assetName, "-osx") {
		os = "darwin"
	} else if strings.Contains(assetName, "-linux") {
		os = "linux"
	}
	if strings.Contains(assetName, "-amd64") || strings.Contains(assetName, "-x64") || strings.Contains(assetName, "-x86_64") {
		arch = "amd64"
	} else if strings.Contains(assetName, "-arm64") || strings.Contains(assetName, "-aarch64") {
		arch = "arm64"
	}
	return
}

func IsLocalAsset(assetOs string, assetArch string) bool {
	if assetOs == "" {
		return true
	}
	if assetOs != runtime.GOOS {
		return false
	}
	if assetArch != "" && assetArch != runtime.GOARCH {
		return false
	}
	return true
}

// If local only is true, then only assets matching the local platform will be downloaded
func DownloadAssets(cfg *GithubConfiguration, localOnly bool) {

	client := CreateGithubClient(cfg)
	var release *github.RepositoryRelease
	// Get latest release
	if cfg.Release == "latest" {
		r, _, err := client.Repositories.GetLatestRelease(context.Background(), cfg.Owner, cfg.Repo)
		cobra.CheckErr(err)
		release = r
	} else if cfg.Release == "all" {
		fmt.Println("Downloading assets for all releases is not supported. Please use 'latest' or specific release.")
		return
	} else {
		release = findRelease(client, cfg.Owner, cfg.Repo, cfg.Release)
		if release == nil {
			fmt.Printf("Could not find release '%v'\n", cfg.Release)
			return
		}
	}
	fmt.Printf("Downloading assets for release '%v' of repo '%v/%v':\n", release.GetTagName(), cfg.Owner, cfg.Repo)
	printRelease(client, cfg.Owner, cfg.Repo, release)
	assets, _, err := client.Repositories.ListReleaseAssets(context.Background(), cfg.Owner, cfg.Repo, release.GetID(), &github.ListOptions{})
	cobra.CheckErr(err)
	for _, asset := range assets {

		if localOnly {
			assetOs, assetArch := DetermineAssetPlatform(asset.GetName())
			if !IsLocalAsset(assetOs, assetArch) {
				fmt.Printf("Skipping asset '%v' [os='%v', arch='%v'] as it does not match the local platform.\n", asset.GetName(), assetOs, assetArch)
				continue
			}
		}

		rc, redirect, err := client.Repositories.DownloadReleaseAsset(context.Background(), cfg.Owner, cfg.Repo, asset.GetID())
		cobra.CheckErr(err)
		if rc != nil {
			fmt.Printf("Not redirected.\n")
			buff := make([]byte, 10*1024)
			n, err := rc.Read(buff)
			fmt.Printf("Read %d bytes\n", n)
			if n > 0 {
				fmt.Printf("Read %d bytes\n", n)
			} else {
				if err == io.EOF {
					rc.Close()
					fmt.Println("Eof")
				} else {
					cobra.CheckErr(err)
				}
			}
		} else {
			DownloadFileFromUrl(redirect, asset.GetName(), DefaultSecurityOptions())
		}
	}
}

func findRelease(client *github.Client, owner string, repo string, tag string) *github.RepositoryRelease {
	allReleases, _, err := client.Repositories.ListReleases(context.Background(), owner, repo, &github.ListOptions{})
	cobra.CheckErr(err)
	for _, release := range allReleases {
		if release.GetTagName() == tag {
			return release
		}
	}
	return nil
}

func printRelease(client *github.Client, owner string, repo string, release *github.RepositoryRelease) {
	fmt.Printf("  %v  [Published: %v]\n", release.GetTagName(), release.GetCreatedAt())
	assets, _, err := client.Repositories.ListReleaseAssets(context.Background(), owner, repo, release.GetID(), &github.ListOptions{})
	cobra.CheckErr(err)
	for _, asset := range assets {
		fmt.Printf("    %v [%v bytes]\n", asset.GetName(), asset.GetSize())
	}
}

func ListGithub(cfg *GithubConfiguration) {
	client := CreateGithubClient(cfg)
	if cfg.Release == "all" {
		fmt.Printf("Listing all releases of repo '%v/%v':\n", cfg.Owner, cfg.Repo)
		allReleases, _, err := client.Repositories.ListReleases(context.Background(), cfg.Owner, cfg.Repo, &github.ListOptions{})
		cobra.CheckErr(err)
		for _, release := range allReleases {
			fmt.Printf("  %v [Published: %v]\n", release.GetTagName(), release.GetPublishedAt())
		}
	} else if cfg.Release == "latest" {
		// Get latest release
		fmt.Printf("Viewing latest release of repo '%v/%v':\n", cfg.Owner, cfg.Repo)
		release, _, err := client.Repositories.GetLatestRelease(context.Background(), cfg.Owner, cfg.Repo)
		cobra.CheckErr(err)
		printRelease(client, cfg.Owner, cfg.Repo, release)
	} else {
		// Get specific release
		fmt.Printf("Viewing release '%v' of repo '%v/%v':\n", cfg.Release, cfg.Owner, cfg.Repo)
		rel := findRelease(client, cfg.Owner, cfg.Repo, cfg.Release)
		if rel == nil {
			fmt.Printf("Could not find a release with tag '%v'\n", cfg.Release)
		} else {
			printRelease(client, cfg.Owner, cfg.Repo, rel)
		}
	}

}

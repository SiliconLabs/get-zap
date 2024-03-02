/*
Copyright © 2024 Silicon Labs
*/
package gh

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

type DownloadOptions struct {
	skipCertCheck  bool
	proxyUrl       *url.URL
	allowHttp      bool
	showPercentage bool
}

func (dso *DownloadOptions) SetProxy(proxyS string) error {
	url, err := url.Parse(proxyS)
	if err != nil {
		return err
	}
	dso.proxyUrl = url
	return nil
}

func (dso *DownloadOptions) SetSkipCertCheck(skipCertCheck bool) {
	dso.skipCertCheck = skipCertCheck
}

func (dso *DownloadOptions) SetAllowHttp(allowHttp bool) {
	dso.allowHttp = allowHttp
}

func (dso *DownloadOptions) SetShowPercentage(showPercentage bool) {
	dso.showPercentage = showPercentage
}

// Returns the default security options.
func DefaultSecurityOptions() *DownloadOptions {
	s := DownloadOptions{
		skipCertCheck:  false,
		proxyUrl:       nil,
		allowHttp:      false,
		showPercentage: true,
	}
	return &s
}

// If local only is true, then only assets matching the local platform will be downloaded
func DownloadAssets(cfg *GithubConfiguration, destinationDirectory string, localOnly bool, suffixOnly string) {

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

		if suffixOnly != "" && !strings.HasSuffix(asset.GetName(), suffixOnly) {
			fmt.Printf("Skipping asset '%v' as it does not have the suffix '%v'.\n", asset.GetName(), suffixOnly)
			continue
		}

		rc, redirect, err := client.Repositories.DownloadReleaseAsset(context.Background(), cfg.Owner, cfg.Repo, asset.GetID())
		cobra.CheckErr(err)
		err = os.MkdirAll(release.GetName(), 0775)
		cobra.CheckErr(err)
		if rc != nil {
			err = downloadFileFromReadCloser(rc, release.GetName(), asset.GetName())
			cobra.CheckErr(err)
		} else {
			err = downloadFileFromUrl(redirect, release.GetName(), asset.GetName(), DefaultSecurityOptions())
			cobra.CheckErr(err)
		}
	}
}

func downloadFileFromReadCloser(rc io.ReadCloser, destinationDirectory string, destinationPath string) error {
	defer rc.Close()
	output, err := os.Create(destinationDirectory + "/" + destinationPath)
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(output, rc)
	return err
}

// This function downloads a file from a given URL and puts it into the
// destination path.
func downloadFileFromUrl(urlAsString string, destinationDirectory string, destinationPath string, sec *DownloadOptions) error {

	tlsConfig := &tls.Config{}
	if sec.skipCertCheck {
		tlsConfig.InsecureSkipVerify = sec.skipCertCheck
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	if sec.proxyUrl != nil {
		tr.Proxy = http.ProxyURL(sec.proxyUrl)
	}

	client := &http.Client{Transport: tr}

	u, err := url.Parse(urlAsString)
	if err != nil {
		return err
	}

	if !sec.allowHttp && u.Scheme == "http" {
		return fmt.Errorf("only secure encrypted HTTPS protocol is allowed, downloads via HTTP are blocked: %v", urlAsString)
	}

	// Security alert: Let's do an actual get now
	response, err := client.Get(urlAsString)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %v", response.StatusCode)
	}

	defer response.Body.Close()

	len := response.ContentLength
	fmt.Printf("Downloading %v bytes to %v ...\n", len, destinationPath)

	output, err := os.Create(destinationDirectory + "/" + destinationPath)
	if err != nil {
		return err
	}
	defer output.Close()

	var chunk int64
	if len > 100 {
		chunk = len / 100
	} else {
		chunk = len
	}
	cnt := 0
	var totalDownloaded int64 = 0
	for {
		written, err := io.CopyN(output, response.Body, chunk)
		totalDownloaded += written
		percentage := (100 * totalDownloaded) / len

		if err == io.EOF {
			fmt.Printf("%v%%: Downloaded %v out of %v bytes. Done!\n", percentage, totalDownloaded, len)
			break
			// done.
		} else if err != nil {
			// real error.
			return err
		} else {
			if sec.showPercentage {
				fmt.Printf("%v%%: Downloaded %v out of %v bytes...\r", percentage, totalDownloaded, len)
			}
			cnt++
		}
	}
	return nil
}

/*
Copyright © 2024 Silicon Labs
*/
package jf

import (
	"fmt"

	"github.com/jfrog/jfrog-client-go/artifactory"
	rtAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/spf13/cobra"
)

type ArtifactoryConfiguration struct {
	Url    string
	ApiKey string
	User   string
	Repo   string
	Path   string
}

func (cfg *ArtifactoryConfiguration) IsValid() bool {
	return (cfg.Url != "" && cfg.ApiKey != "" && cfg.User != "")
}

func (cfg *ArtifactoryConfiguration) CreateDetails() *auth.ServiceDetails {
	if !cfg.IsValid() {
		cobra.CheckErr(fmt.Errorf("invalid artifactory configuration, you need to provide url, api key and user either via command line, environment variables, or configuration file"))
	}
	rtDetails := rtAuth.NewArtifactoryDetails()
	rtDetails.SetUrl(cfg.Url)
	rtDetails.SetApiKey(cfg.ApiKey)
	rtDetails.SetUser(cfg.User)
	return &rtDetails
}

func ArtifactoryDelete(cfg *ArtifactoryConfiguration, pattern string) {
	rtDetails := cfg.CreateDetails()

	s, err := config.NewConfigBuilder().SetServiceDetails(*rtDetails).Build()
	cobra.CheckErr(err)

	m, err := artifactory.New(s)
	cobra.CheckErr(err)

	params := services.NewDeleteParams()
	params.Pattern = cfg.Repo + "/" + pattern
	fmt.Printf("Deleting files from %v/%v: %v\n", cfg.Url, cfg.Repo, params.Pattern)

	pathsToDelete, err := m.GetPathsToDelete(params)
	cobra.CheckErr(err)
	defer pathsToDelete.Close()
	cnt, err := m.DeleteFiles(pathsToDelete)
	cobra.CheckErr(err)
	fmt.Printf("Deleted files: %v\n", cnt)
}

func ArtifactoryPattern(release string, isLocal bool) string {
	return release + "/**"
}

func ArtifactoryDownload(cfg *ArtifactoryConfiguration, pattern string) int {

	rtDetails := cfg.CreateDetails()

	s, err := config.NewConfigBuilder().SetServiceDetails(*rtDetails).Build()
	cobra.CheckErr(err)

	m, err := artifactory.New(s)
	cobra.CheckErr(err)

	params := services.NewDownloadParams()
	params.Pattern = cfg.Repo + "/" + pattern
	fmt.Printf("Downloading files from %v/%v: %v\n", cfg.Url, cfg.Repo, params.Pattern)
	success, failures, err := m.DownloadFiles(params)
	cobra.CheckErr(err)

	fmt.Printf("Downloaded files: success %v, failure %v\n", success, failures)
	return success
}

func ArtifactoryUpload(cfg *ArtifactoryConfiguration, pattern string) {
	rtDetails := cfg.CreateDetails()

	s, err := config.NewConfigBuilder().SetServiceDetails(*rtDetails).Build()
	cobra.CheckErr(err)

	m, err := artifactory.New(s)
	cobra.CheckErr(err)

	params := services.NewUploadParams()
	params.Pattern = pattern
	fmt.Printf("Uploading files to %v/%v: %v\n", cfg.Url, cfg.Repo, params.Pattern)
	params.Target = cfg.Repo + "/"

	success, failures, err := m.UploadFiles(params)
	cobra.CheckErr(err)
	fmt.Printf("Uploaded files: success %v, failure %v\n", success, failures)
}

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

func (cfg *ArtifactoryConfiguration) Validate() {
	if !(cfg.Url != "" && cfg.ApiKey != "" && cfg.User != "") {
		cobra.CheckErr(fmt.Errorf("Invalid artifactory configuration. You need to provide url, api key and user either via command line, environment variables, or configuration file."))
	}
}

func (cfg *ArtifactoryConfiguration) CreateDetails() *auth.ServiceDetails {
	cfg.Validate()
	rtDetails := rtAuth.NewArtifactoryDetails()
	rtDetails.SetUrl(cfg.Url)
	rtDetails.SetApiKey(cfg.ApiKey)
	rtDetails.SetUser(cfg.User)
	return &rtDetails
}

func ArtifactoryDownload(cfg *ArtifactoryConfiguration) {

	rtDetails := cfg.CreateDetails()

	s, err := config.NewConfigBuilder().SetServiceDetails(*rtDetails).Build()
	cobra.CheckErr(err)

	m, err := artifactory.New(s)
	cobra.CheckErr(err)

	params := services.NewDownloadParams()
	params.Pattern = cfg.Repo + "/" + cfg.Path
	fmt.Println("Downloading files from", cfg.Url, "with pattern", params.Pattern)
	success, failures, err := m.DownloadFiles(params)
	cobra.CheckErr(err)

	fmt.Printf("Download files: success %v, failure %v\n", success, failures)
}

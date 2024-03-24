/*
Copyright © 2024 Silicon Labs
*/
package cmd

import (
	"fmt"
	"os"
	"silabs/get-zap/gh"
	"silabs/get-zap/jf"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const ownerArg = "ghOwner"
const repoArg = "ghRepo"
const githubTokenArg = "ghToken"
const releaseArg = "ghRelease"
const assetArg = "ghAsset"
const rtUrl = "rtUrl"
const rtApiKey = "rtApiKey"
const rtUser = "rtUser"
const rtRepo = "rtRepo"
const rtPath = "rtPath"

const useRt = "useRt"
const useGh = "useGh"
const localRoot = "localRoot"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "get-zap",
	Short: "Application to retrieve artifacts from github.",
	Long:  `This application by default retrieves zap artifacts, with the right arguments, it can be used to retrieve assets from any public github repo.`,
	Run: func(cmd *cobra.Command, args []string) {
		Fetch(ReadGithubConfiguration(), ReadArtifactoryConfiguration(), viper.GetBool(useGh), viper.GetBool(useRt))
	},
}

func ReadArtifactoryConfiguration() *jf.ArtifactoryConfiguration {
	return &jf.ArtifactoryConfiguration{
		Url:    viper.GetString(rtUrl),
		ApiKey: viper.GetString(rtApiKey),
		User:   viper.GetString(rtUser),
		Repo:   viper.GetString(rtRepo),
		Path:   viper.GetString(rtPath),
	}
}

func ReadGithubConfiguration() *gh.GithubConfiguration {
	return &gh.GithubConfiguration{
		Owner:   viper.GetString(ownerArg),
		Repo:    viper.GetString(repoArg),
		Token:   viper.GetString(githubTokenArg),
		Release: viper.GetString(releaseArg),
		Asset:   viper.GetString(assetArg),
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initViper)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is $HOME/.%s.json)", configKey))
	rootCmd.PersistentFlags().String(ownerArg, "project-chip", "Owner of the github repository.")
	rootCmd.PersistentFlags().String(repoArg, "zap", "Name of the github repository.")
	rootCmd.PersistentFlags().StringP(githubTokenArg, "t", "", "Github token to use for authentication.")
	rootCmd.PersistentFlags().StringP(releaseArg, "r", "latest", "Release to download. Specify a name, or 'all' or 'latest' for all releases.")
	rootCmd.PersistentFlags().String(localRoot, ".", "Local root directory to download assets to. All operations are limited to within this directory.")
	rootCmd.PersistentFlags().StringP(assetArg, "a", "local", "Asset to download. Specify a name, or 'all' or 'local' for matching the platform.")
	rootCmd.PersistentFlags().String(rtUrl, "", "Artifactory URL.")
	rootCmd.PersistentFlags().String(rtApiKey, "", "Artifactory API Key.")
	rootCmd.PersistentFlags().String(rtUser, "", "Artifactory user.")
	rootCmd.PersistentFlags().String(rtRepo, "", "Artifactory repository.")
	rootCmd.PersistentFlags().String(rtPath, "", "Artifactory path within the repo.")
	rootCmd.PersistentFlags().Bool(useRt, true, "Use Artifactory.")
	rootCmd.PersistentFlags().Bool(useGh, true, "Use GitHub.")
}

// This function initializes viper, which is used to read configuration from environment variables and config files.
const configKey = "get_zap"

var configFile = ""

func initViper() {

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("." + configKey) // The config file is the configKey, prepended with a dot
		viper.AddConfigPath("$HOME")         // Look in user home directory
		viper.SetConfigType("json")          // And the file is in JSON format.
	}

	viper.ReadInConfig()
	viper.SetEnvPrefix(configKey)
	viper.AutomaticEnv()

	// We bind all the flags, so that variables set by viper, but marked as required, are considered properly set
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && viper.IsSet(f.Name) {
			rootCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.Flags())
}

// EO Viper configuration.

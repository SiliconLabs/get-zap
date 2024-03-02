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

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "get-zap",
	Short: "Application to retrieve artifacts from github.",
	Long:  `This application by default retrieves zap artifacts, with the right arguments, it can be used to retrieve assets from any public github repo.`,
	Run: func(cmd *cobra.Command, args []string) {
		gh.DefaultAction(ReadGithubConfiguration(), ReadArtifactoryConfiguration())
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.get-zap.yaml)")
	rootCmd.PersistentFlags().String(ownerArg, "project-chip", "Owner of the github repository.")
	rootCmd.PersistentFlags().String(repoArg, "zap", "Name of the github repository.")
	rootCmd.PersistentFlags().StringP(githubTokenArg, "t", "", "Github token to use for authentication.")
	rootCmd.PersistentFlags().StringP(releaseArg, "r", "latest", "Release to download. Specify a name, or 'all' or 'latest' for all releases.")
	rootCmd.PersistentFlags().StringP(assetArg, "a", "local", "Asset to download. Specify a name, or 'all' or 'local' for matching the platform.")
	rootCmd.PersistentFlags().String(rtUrl, "", "Artifactory URL.")
	rootCmd.PersistentFlags().String(rtApiKey, "", "Artifactory API Key.")
	rootCmd.PersistentFlags().String(rtUser, "", "Artifactory user.")
	rootCmd.PersistentFlags().String(rtRepo, "", "Artifactory repository.")
	rootCmd.PersistentFlags().String(rtPath, "", "Artifactory path within the repo.")

	viper.BindPFlag(ownerArg, rootCmd.PersistentFlags().Lookup(ownerArg))
	viper.BindPFlag(repoArg, rootCmd.PersistentFlags().Lookup(repoArg))
	viper.BindPFlag(githubTokenArg, rootCmd.PersistentFlags().Lookup(githubTokenArg))
	viper.BindPFlag(releaseArg, rootCmd.PersistentFlags().Lookup(releaseArg))
	viper.BindPFlag(assetArg, rootCmd.PersistentFlags().Lookup(assetArg))
	viper.BindPFlag(rtUrl, rootCmd.PersistentFlags().Lookup(rtUrl))
	viper.BindPFlag(rtApiKey, rootCmd.PersistentFlags().Lookup(rtApiKey))
	viper.BindPFlag(rtUser, rootCmd.PersistentFlags().Lookup(rtUser))
	viper.BindPFlag(rtRepo, rootCmd.PersistentFlags().Lookup(rtRepo))
	viper.BindPFlag(rtPath, rootCmd.PersistentFlags().Lookup(rtPath))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".get-zap" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".get-zap")
	}

	viper.SetEnvPrefix("get_zap") // will be uppercased automatically
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

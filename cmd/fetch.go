/*
Copyright © 2024 Silicon Labs
*/
package cmd

import (
	"fmt"
	"runtime"
	"silabs/get-zap/gh"
	"silabs/get-zap/jf"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Retrieves an artifact either from Artifactory cache, or from Github, wherever it can be found.",
	Long: `This is the default operation, if you don't pass any commands. The action performed is:
- first the existence of release artifact is tested on Artifactory. If it's found, it's downloaded from there.
- if it's not found on Artifactory, it's downloaded from Github. If it's found on Github, it's downloaded from there.
- after the artifact is found on Github, it is uploaded to Artifactory, so that it can be found there next time
	
Note: command line arguments can modify this flow.`,
	Run: func(cmd *cobra.Command, args []string) {
		Fetch(ReadGithubConfiguration(), ReadArtifactoryConfiguration(), viper.GetBool(useGh), viper.GetBool(useRt))
	},
}

// This is what gets executed if no toplevel commands are passed.
func Fetch(ghCfg *gh.GithubConfiguration, rtCfg *jf.ArtifactoryConfiguration, useGh bool, useRt bool) {
	if !useGh && !useRt {
		fmt.Println("Neither Artifactory nor Github are enabled, nothing to do.")
		return
	}

	if !useGh {
		// We only check artifactory, if we don't find it, we're done.
		if ghCfg.Release == "latest" || ghCfg.Release == "all" {
			fmt.Printf("Artifactory does not cache 'latest' or 'all' releases. When using --useGh=false, please specify a specific release.\n")
			return
		}
		jf.ArtifactoryDownload(rtCfg, ghCfg.Release+"/**")
		return
	}

	if !useRt {
		// We only attempt to download from github, if we don't find it, we're done.
		fmt.Printf("Downloading release '%v' of repo '%v/%v' for the platform '%v/%v'...\n", ghCfg.Release, ghCfg.Owner, ghCfg.Repo, runtime.GOOS, runtime.GOARCH)
		gh.DownloadAssets(ghCfg, ".", true, ".zip")
		return
	}

	// If we get here, we're going to do the following: first we attempt to download the assset from artifactory. If we can't find it, we will download it
	// from github. If we do find it, we will then upload it to artifactory for the next time someone tries to download this same thing.
	fmt.Printf("Full cycle not yet implemented.\n")

}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

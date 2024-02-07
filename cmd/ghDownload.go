/*
Copyright © 2024 Silicon Labs
*/
package cmd

import (
	"silabs/get-zap/gh"

	"github.com/spf13/cobra"
)

// ghDownloadCmd represents the ghDownload command
var ghDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads assets from Github",
	Long:  `This command can be used to download assets from Github.`,
	Run: func(cmd *cobra.Command, args []string) {
		gh.DownloadAssets(ReadGithubConfiguration(), false)
	},
}

func init() {
	ghCmd.AddCommand(ghDownloadCmd)
}

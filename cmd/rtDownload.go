/*
Copyright © 2024 Silicon Labs
*/
package cmd

import (
	"silabs/get-zap/jf"

	"github.com/spf13/cobra"
)

var rtDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads a file from Artifactory",
	Long:  `Performs a download of one or more files from Artifactory.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("file")
		cobra.CheckErr(err)
		jf.ArtifactoryDownload(ReadArtifactoryConfiguration(), file)
	},
}

func init() {
	rtCmd.AddCommand(rtDownloadCmd)
	rtDownloadCmd.Flags().StringP("file", "f", "", "File to upload")
	rtDownloadCmd.MarkFlagRequired("file")
}

/*
Copyright © 2024 Silicon Labs
*/
package cmd

import (
	"silabs/get-zap/jf"

	"github.com/spf13/cobra"
)

var rtUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads a file to artifactory",
	Long:  `Performs an upload of specified files to artifactory.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("file")
		cobra.CheckErr(err)
		jf.ArtifactoryUpload(ReadArtifactoryConfiguration(), file)
	},
}

func init() {
	rtCmd.AddCommand(rtUploadCmd)
	rtUploadCmd.Flags().StringP("file", "f", "", "File to upload")
	rtUploadCmd.MarkFlagRequired("file")
}

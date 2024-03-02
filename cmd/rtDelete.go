/*
Copyright © 2024 Silicon Labs
*/
package cmd

import (
	"silabs/get-zap/jf"

	"github.com/spf13/cobra"
)

var rtDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a file on artifactory",
	Long:  `Performs a file delete of a given pattern.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("file")
		cobra.CheckErr(err)
		jf.ArtifactoryDelete(ReadArtifactoryConfiguration(), file)
	},
}

func init() {
	rtCmd.AddCommand(rtDeleteCmd)
	rtDeleteCmd.Flags().StringP("file", "f", "", "File to upload")
	rtDeleteCmd.MarkFlagRequired("file")
}

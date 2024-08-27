/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vincent-vinf/image-tool/pkg/image"
	"github.com/vincent-vinf/image-tool/pkg/output"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",

	RunE: func(cmd *cobra.Command, args []string) error {
		imageFiles, err := output.ReadImageFromDir(outputDir)
		if err != nil {
			return fmt.Errorf("could not read the images directory: %w", err)
		}
		for _, i := range imageFiles {
			tarPath, err := gunzipImage(i)
			if err != nil {
				return err
			}
			err = image.CheckImageTar(cmd.Context(), tarPath)
			if err != nil {
				return fmt.Errorf("could not check image tar: %w", err)
			}
			logger.Infof("image %s is ok", i.URL)
		}
		logger.Infof("all image files are ok")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

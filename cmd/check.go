/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

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
		errMap := make(map[string]error)
		for _, i := range imageFiles {
			tarPath, err := gunzipImage(i)
			if err != nil {
				return err
			}

			var os, arch string
			if platform != "" {
				parts := strings.Split(platform, "/")
				if len(parts) < 2 {
					return fmt.Errorf("invalid platform format: %s", platform)
				}
				os = parts[0]
				arch = parts[1]
			}
			err = image.CheckImageTar(cmd.Context(), tarPath, os, arch)
			if err != nil {
				errMap[i.URL] = err
			} else {
				logger.Infof("image %s is ok", i.URL)
			}
		}
		if len(errMap) > 0 {
			for imageURL, err := range errMap {
				logger.Errorf("failed to check image %s, reason: %v", imageURL, err)
			}

			return fmt.Errorf("failed to check %d images", len(errMap))
		}
		logger.Infof("all image files are ok")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

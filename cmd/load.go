package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"image-tool/pkg/input"
)

var (
	registry string
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load image from image tar to registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if registry == "" {
			return fmt.Errorf("--registry is required")
		}
		var (
			images []string
			err    error
		)

		if imageListPath != "" {
			images, err = input.ReadImagesFile(imageListPath)
			if err != nil {
				return fmt.Errorf("could not read the images file: %w", err)
			}
			logger.Infof("found %d images for file %s", len(images), imageListPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.PersistentFlags().StringVarP(&registry, "registry", "r", "", "example: harbor.qkd.cn:8443/library")
}

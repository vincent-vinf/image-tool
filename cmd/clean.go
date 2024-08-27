package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vincent-vinf/image-tool/pkg/image"
	"github.com/vincent-vinf/image-tool/pkg/output"
	"github.com/vincent-vinf/image-tool/pkg/zip"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "delete all compressed files or all uncompressed files",
	RunE: func(cmd *cobra.Command, args []string) error {
		imageFiles, err := output.ReadImageFromDir(outputDir)
		if err != nil {
			return fmt.Errorf("could not read the images directory: %w", err)
		}
		logger.Infof("found %d image file in %s", len(imageFiles), outputDir)
		if err != nil {
			return err
		}
		err = cleanTar(imageFiles)
		if err != nil {
			return err
		}
		logger.Infof("successfully cleaned up")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func cleanTar(imageFiles []*image.TarImage) error {
	for _, i := range imageFiles {
		for _, f := range i.Files {
			if f.Compress == zip.Uncompressed {
				logger.Infof("%s will be remove", f.Path)
				err := os.Remove(f.Path)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vincent-vinf/image-tool/pkg/output"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean",
	RunE: func(cmd *cobra.Command, args []string) error {
		imageFiles, err := output.ReadImageFromDir(outputDir)
		if err != nil {
			return fmt.Errorf("could not read the images directory: %w", err)
		}
		logger.Infof("found %d image file in %s", len(imageFiles), outputDir)
		//for _, f := range imageFiles {
		//}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

//func cleanOnlyKeep(images []*image.TarImage, compression zip.Compression) error {
//	for _, i := range images {
//		for _, f := range i.Files {
//
//		}
//	}
//}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vincent-vinf/image-tool/pkg/image"
	"github.com/vincent-vinf/image-tool/pkg/input"
	"github.com/vincent-vinf/image-tool/pkg/output"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			logger.Warnf("args will be ignored")
		}
		images, err := input.ReadImagesFile(imageListPath)
		if err != nil {
			return fmt.Errorf("could not read the images file: %w", err)
		}
		logger.Infof("found %d images from file %s", len(images), imageListPath)

		imageFiles, err := output.ReadImageFromDir(outputDir)
		if err != nil {
			return fmt.Errorf("could not read the images directory: %w", err)
		}
		logger.Infof("found %d image file in %s", len(imageFiles), outputDir)

		diff := image.GetDiff(images, imageFiles)
		logger.Infof("will pull %d image\n%s", len(diff.Add), strings.Join(diff.Add, "\n"))
		var delImages []string
		for _, i := range diff.Del {
			delImages = append(delImages, i.URL)
		}
		logger.Infof("will delete %d image\n%s", len(diff.Del), strings.Join(delImages, "\n"))

		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create the output directory: %w", err)
		}

		errMap := make(map[string]error)
		for _, i := range diff.Add {
			err := saveImage(cmd.Context(), i)
			if err != nil {
				errMap[i] = err
			}
		}
		for _, i := range diff.Del {
			for _, f := range i.Files {
				err := os.Remove(f.Path)
				if err != nil {
					errMap[i.URL] = err
				}
			}
		}

		if len(errMap) > 0 {
			for i, err := range errMap {
				logger.Errorf("failed to sync %s, reason: %v", i, err)
			}
			return fmt.Errorf("failed to sync the %d images", len(errMap))
		}

		logger.Infof("successfully sync %d images, ", len(images))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

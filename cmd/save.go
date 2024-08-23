package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"image-tool/pkg/image"
	"image-tool/pkg/input"
	"image-tool/pkg/output"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Download the images listed in image.txt to the directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		images, err := input.ReadImagesFile(imageListPath)
		if err != nil {
			return fmt.Errorf("could not read the images file: %w", err)
		}
		logger.Infof("found %d images for file %s", len(images), imageListPath)

		imageFiles, err := output.ReadImageFromDir(outputDir)
		if err != nil {
			return fmt.Errorf("could not read the images directory: %w", err)
		}
		logger.Infof("found %d image file in %s", len(imageFiles), outputDir)
		diff := image.GetDiff(images, imageFiles)
		logger.Infof("will save %d image\n%s", len(diff.Add), strings.Join(diff.Add, "\n"))

		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create the output directory: %w", err)
		}

		for _, i := range images {
			var exist bool
			for _, file := range imageFiles {
				if i == file.URL {
					exist = true
					break
				}
			}
			if exist {
				logger.Infof("image %s already exist", i)
				continue
			}

			err := PullSave(i, path.Join(outputDir, image.ConvertToFilename(i)+".tar"), "")
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
	// saveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func PullSave(imageURL, dst, platform string) error {
	logger.Infof("start save image: %s", imageURL)
	err := imageHandler.Pull(imageURL, platform)
	if err != nil {
		return fmt.Errorf("fail to pull image: %w", err)
	}
	err = imageHandler.Save(imageURL, dst)
	if err != nil {
		return fmt.Errorf("fail to save image: %w", err)
	}
	logger.Infof("successfully saved image: %s to %s\n", imageURL, dst)

	return nil
}

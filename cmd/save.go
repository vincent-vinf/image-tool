package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"image-tool/pkg/image"
	"image-tool/pkg/input"
	"image-tool/pkg/output"
	"image-tool/pkg/utils"
	"image-tool/pkg/zip"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Download the images listed in image.txt to the directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		argImages, err := readImagesFromArgs(args)
		if err != nil {
			return err
		}
		logger.Infof("found %d images from args", len(argImages))

		fileImages, err := input.ReadImagesFile(imageListPath)
		if err != nil {
			return fmt.Errorf("could not read the images file: %w", err)
		}
		logger.Infof("found %d images from file %s", len(fileImages), imageListPath)

		images := append(fileImages, argImages...)

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

		errMap := make(map[string]error)
		for _, i := range diff.Add {
			err := saveImage(cmd.Context(), i)
			if err != nil {
				errMap[i] = err
			}
		}

		if len(errMap) > 0 {
			for i, err := range errMap {
				logger.Errorf("failed to save %s, reason: %v", i, err)
			}
			return fmt.Errorf("failed to save the %d images", len(errMap))
		}
		logger.Infof("successfully save %d images, ", len(diff.Add))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}

func readImagesFromArgs(args []string) ([]string, error) {
	var images []string
	for _, arg := range args {
		if !utils.ValidateImageName(arg) {
			return nil, fmt.Errorf("invalid image name %s", arg)
		}
		images = append(images, arg)
	}

	return images, nil
}

func saveImage(ctx context.Context, imageURL string) error {
	tarPath := path.Join(outputDir, image.ConvertToFilename(imageURL)+".tar")
	err := image.PullImageToTar(ctx, imageURL, platform, username, password, tarPath)
	if err != nil {
		return err
	}
	if autoZip {
		err = compressImage(imageURL, tarPath)
		if err != nil {
			return err
		}
		err = os.Remove(tarPath)
		if err != nil {
			return fmt.Errorf("could not remove the tar file: %w", err)
		}
	}

	return nil
}

func compressImage(imageURL string, tarFile string) error {
	zipTarPath := path.Join(outputDir, image.ConvertToFilename(imageURL)+".tgz")
	logger.Infof("zip image %s tar to %s", imageURL, zipTarPath)
	return zip.Compress(tarFile, zipTarPath, zip.Gzip)
}

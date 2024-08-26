package cmd

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"image-tool/pkg/image"
	"image-tool/pkg/input"
	"image-tool/pkg/output"
	"image-tool/pkg/zip"
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

		argImages, err := readImagesFromArgs(args)
		if err != nil {
			return err
		}
		logger.Infof("found %d images from args", len(argImages))

		var fileImages []string
		if imageListPath != "" {
			fileImages, err = input.ReadImagesFile(imageListPath)
			if err != nil {
				return fmt.Errorf("could not read the images file: %w", err)
			}
			logger.Infof("found %d images for file %s", len(fileImages), imageListPath)
		}
		images := append(fileImages, argImages...)

		imageFiles, err := output.ReadImageFromDir(outputDir)
		if err != nil {
			return fmt.Errorf("could not read the images directory: %w", err)
		}
		logger.Infof("found %d image file in %s", len(imageFiles), outputDir)

		errMap := make(map[string]error)
		if len(images) == 0 {
			logger.Infof("will load all images in %s", outputDir)
			for _, image := range imageFiles {
				err := loadImage(cmd.Context(), image, registry)
				if err != nil {
					errMap[image.URL] = err
				}
			}
		} else {
			diff := image.GetDiff(images, imageFiles)
			for _, image := range diff.Exists {
				err := loadImage(cmd.Context(), image, registry)
				if err != nil {
					errMap[image.URL] = err
				}
			}
			for _, image := range diff.Add {
				errMap[image] = fmt.Errorf("image tar file does not exist")
			}
		}

		if len(errMap) > 0 {
			for imageURL, err := range errMap {
				logger.Errorf("failed to load %s, reason: %v", imageURL, err)
			}

			return fmt.Errorf("failed to load %d images", len(errMap))
		}

		logger.Infof("successfully loaded %d image files", max(len(images), len(imageFiles)))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.PersistentFlags().StringVarP(&registry, "registry", "r", "", "example: harbor.qkd.cn:8443/library")
}

func loadImage(ctx context.Context, i *image.TarImage, newRegistry string) error {
	// ignore the existence of multiple files with the same compression format
	m := make(map[zip.Compression]image.File)
	for _, file := range i.Files {
		m[file.Compress] = file
	}

	var tarPath string
	if f, ok := m[zip.Uncompressed]; !ok {
		// no uncompressed image file exists, try decompressing it
		if f, ok := m[zip.Gzip]; ok {
			tarPath = path.Join(filepath.Dir(f.Path), image.ConvertToFilename(i.URL)+".tar")
			err := zip.Decompress(f.Path, tarPath, f.Compress)
			if err != nil {
				return fmt.Errorf("could not decompress image: %s, %w", f.Path, err)
			}
		} else {
			return fmt.Errorf("no uncompressed or gzip-compressed image file exists %s", i.URL)
		}
	} else {
		tarPath = f.Path
	}
	sp := strings.Split(i.URL, "/")
	newURL := fmt.Sprintf("%s/%s", newRegistry, sp[len(sp)-1])
	err := image.PushImageToRegistry(ctx, tarPath, newURL, platform, username, password)
	if err != nil {
		return fmt.Errorf("could not load the image: %w", err)
	}
	return nil
}

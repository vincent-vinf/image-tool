package cmd

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vincent-vinf/image-tool/pkg/image"
	"github.com/vincent-vinf/image-tool/pkg/input"
	"github.com/vincent-vinf/image-tool/pkg/output"
	"github.com/vincent-vinf/image-tool/pkg/zip"
)

var (
	registry string
	toDaemon bool
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load image from image tar to registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !toDaemon && registry == "" {
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
			logger.Infof("found %d images from file %s", len(fileImages), imageListPath)
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
				err := loadImage(cmd.Context(), image, registry, toDaemon)
				if err != nil {
					errMap[image.URL] = err
				}
			}
		} else {
			diff := image.GetDiff(images, imageFiles)
			for _, image := range diff.Exists {
				err := loadImage(cmd.Context(), image, registry, toDaemon)
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

		logger.Infof("image file loaded successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.PersistentFlags().StringVarP(&registry, "registry", "r", "", "example: harbor.qkd.cn:8443/library")
	loadCmd.PersistentFlags().BoolVar(&toDaemon, "daemon", false, "just load image to local daemon")
}

func loadImage(ctx context.Context, i *image.TarImage, newRegistry string, toDaemon bool) error {
	tarPath, err := gunzipImage(i)
	if err != nil {
		return err
	}

	newURL := i.URL

	if newRegistry != "" {
		sp := strings.Split(i.URL, "/")
		newURL = fmt.Sprintf("%s/%s", strings.TrimRight(newRegistry, "/"), sp[len(sp)-1])
	}

	if toDaemon {
		err = image.LoadImageToDocker(ctx, tarPath, newURL)
	} else {
		err = image.PushImageToRegistry(ctx, tarPath, newURL, platform, username, password)
	}

	if err != nil {
		return fmt.Errorf("could not load the image: %w", err)
	}

	return nil
}

func gunzipImage(i *image.TarImage) (string, error) {
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
				return "", fmt.Errorf("could not decompress image: %s, %w", f.Path, err)
			}
		} else {
			return "", fmt.Errorf("no uncompressed or gzip-compressed image file exists %s", i.URL)
		}
	} else {
		tarPath = f.Path
	}

	return tarPath, nil
}

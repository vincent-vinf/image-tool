package output

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"image-tool/pkg/image"
	"image-tool/pkg/utils"
	"image-tool/pkg/zip"
)

var (
	logger = utils.GetLogger()
)

func ReadImageFromDir(dir string) ([]*image.TarImage, error) {
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	imageMap := make(map[string]*image.TarImage)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		i, err := filenameToImageURL(path)
		if err != nil {
			logger.Warnf("mismatched formats to read file %s: %s", path, err)
			return nil
		}
		compression, err := zip.DetectFileCompression(path)
		if err != nil {
			logger.Warnf("failed to detect compression type %s", path)
			return nil
		}
		file := image.File{
			Path:     path,
			Compress: compression,
		}

		if m, ok := imageMap[i]; ok {
			m.Files = append(m.Files, file)
		} else {
			imageMap[i] = &image.TarImage{
				URL:   i,
				Files: []image.File{file},
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	var images []*image.TarImage
	for _, i := range imageMap {
		images = append(images, i)
	}

	return images, nil
}

func filenameToImageURL(file string) (string, error) {
	name := filepath.Base(file)
	name = strings.TrimSuffix(name, ".tar")
	name = strings.TrimSuffix(name, ".tar.gz")
	name = strings.TrimSuffix(name, ".tgz")

	imageName := image.ConvertToImageName(name)
	if !utils.ValidateImageName(imageName) {
		return "", fmt.Errorf("illegal mirror name %s", imageName)
	}

	return imageName, nil
}

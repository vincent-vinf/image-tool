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

func ReadImageFromDir(dir string) ([]image.File, error) {
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var images []image.File
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		i, err := fileToImage(path)
		if err != nil {
			logger.Warnf("mismatched formats to read file %s: %s", path, err)
		}
		images = append(images, i)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return images, nil
}

func fileToImage(file string) (image.File, error) {
	name := filepath.Base(file)
	name = strings.TrimSuffix(name, ".tar")
	name = strings.TrimSuffix(name, ".tar.gz")
	name = strings.TrimSuffix(name, ".tgz")

	imageName := image.ConvertToImageName(name)
	if !utils.ValidateImageName(imageName) {
		return image.File{}, fmt.Errorf("illegal mirror name %s", imageName)
	}
	compression, err := zip.DetectFileCompression(file)
	if err != nil {
		return image.File{}, fmt.Errorf("failed to detect compression type %s", file)
	}
	i := image.File{
		URL:      imageName,
		Path:     file,
		Compress: compression,
	}

	return i, nil
}

package image

import (
	"strings"

	"image-tool/pkg/zip"
)

// ConvertToFilename converts a Docker image name to a valid filename
// It replaces '/' with '#' and '@' with '+', but keeps '.' unchanged.
func ConvertToFilename(image string) string {
	filename := strings.ReplaceAll(image, "/", "#")
	filename = strings.ReplaceAll(filename, ":", "+")
	return filename
}

func ConvertToImageName(filename string) string {
	image := strings.ReplaceAll(filename, "#", "/")
	image = strings.ReplaceAll(image, "+", ":")
	return image
}

type File struct {
	URL      string
	Path     string
	Compress zip.Compression
}

type Diff struct {
	Add []string
	Del []File
}

func GetDiff(targetImages []string, existImages []File) Diff {
	var d Diff
	for _, image := range targetImages {
		var exist bool
		for _, file := range existImages {
			if image == file.URL {
				exist = true
				d.Del = append(d.Del)
				break
			}
		}
		if !exist {
			d.Add = append(d.Add, image)
		}
	}

	return d
}

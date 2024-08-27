package image

import (
	"strings"

	"github.com/vincent-vinf/image-tool/pkg/zip"
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

type TarImage struct {
	URL   string
	Files []File
}

type File struct {
	Path     string
	Compress zip.Compression
}

type Diff struct {
	Exists []*TarImage
	Add    []string
	Del    []*TarImage
}

func GetDiff(targetImages []string, existImages []*TarImage) Diff {
	var d Diff
	// 对于image.txt文件中的镜像，判断是否已被保存到本地
	for _, image := range targetImages {
		var exist bool
		for _, i := range existImages {
			if image == i.URL {
				exist = true
				break
			}
		}
		if !exist {
			d.Add = append(d.Add, image)
		}
	}

	// 对于本地到tar包，判断是否在image.txt中
	for _, i := range existImages {
		var exist bool
		for _, image := range targetImages {
			if i.URL == image {
				exist = true
				break
			}
		}
		if exist {
			d.Exists = append(d.Exists, i)
		} else {
			d.Del = append(d.Del, i)
		}
	}

	return d
}

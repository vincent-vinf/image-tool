package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vincent-vinf/image-tool/pkg/utils"
)

func ReadImagesFile(file string) ([]string, error) {
	var images []string

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !utils.ValidateImageName(line) {
			return nil, fmt.Errorf("invalid image: %s", line)
		}
		images = append(images, line)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

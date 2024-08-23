package image

import (
	"fmt"
	"os"
	"os/exec"

	"image-tool/pkg/utils"
)

var (
	logger = utils.GetLogger()
)

type Handler interface {
	Pull(imageURL string, platform string) error
	Save(imageURL string, dst string) error
}

type DockerHandler struct {
}

func (h DockerHandler) Pull(imageURL string, platform string) error {
	var cmd *exec.Cmd
	if platform != "" {
		cmd = exec.Command("docker", "pull", fmt.Sprintf("--platform=%s", platform), imageURL)
	} else {
		cmd = exec.Command("docker", "pull", imageURL)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to pull image: %s, error: %v", imageURL, err)
	}

	return nil
}

// Save 保存 Docker 镜像为 tar 文件
func (h DockerHandler) Save(imageURL string, dst string) error {
	cmd := exec.Command("docker", "save", "-o", dst, imageURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to save image: %s to destination: %s, error: %v, output: %s", imageURL, dst, err, string(output))
	}

	return nil
}

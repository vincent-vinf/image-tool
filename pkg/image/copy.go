package image

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"image-tool/pkg/utils"
)

var (
	logger = utils.GetLogger()
)

func PullImageToTar(ctx context.Context, srcImage, platform, username, passwd, dstTar string) error {
	src, err := NewRegistryImageNode(registryImageKey(srcImage), platform, username, passwd)
	if err != nil {
		return err
	}
	dst := NewImageNode(tarImageKey(dstTar))

	return copyImage(ctx, src, dst)
}

func LoadImageToDocker(ctx context.Context, srcTar string, dstDaemon string) error {
	src := NewImageNode(tarImageKey(srcTar))
	dst := NewImageNode(daemonImageKey(dstDaemon))

	return copyImage(ctx, src, dst)
}

func PushImageToRegistry(ctx context.Context, srcTar string, registry string, platform, dstUsername, dstPasswd string) error {
	src := NewImageNode(tarImageKey(srcTar))
	dst, err := NewRegistryImageNode(registryImageKey(registry), platform, dstUsername, dstPasswd)
	if err != nil {
		return err
	}

	return copyImage(ctx, src, dst)
}

func CopyBetweenRegistry(ctx context.Context,
	srcImage, platform, srcUsername, srcPasswd string,
	dstImage, dstUsername, dstPasswd string,
	authFilePath string,
) error {
	src, err := NewRegistryImageNode(registryImageKey(srcImage), platform, srcUsername, srcPasswd)
	if err != nil {
		return err
	}
	dst, err := NewRegistryImageNode(registryImageKey(dstImage), platform, dstUsername, dstPasswd)
	if err != nil {
		return err
	}

	return copyImage(ctx, src, dst)
}

func registryImageKey(s string) string {
	return fmt.Sprintf("docker://%s", s)
}

func tarImageKey(s string) string {
	return fmt.Sprintf("docker-archive:%s", s)
}

func daemonImageKey(s string) string {
	return fmt.Sprintf("docker-daemon:%s", s)
}

func NewImageNode(imageKey string) ImageNode {
	return ImageNode{
		ImageKey: imageKey,
	}
}

func NewRegistryImageNode(imageKey, platform, username, passwd string) (ImageNode, error) {
	src := ImageNode{
		ImageKey: imageKey,
	}
	if platform != "" {
		parts := strings.Split(platform, "/")

		if len(parts) < 2 {
			return ImageNode{}, fmt.Errorf("invalid platform format: %s", platform)
		}
		src.Platform = &Platform{
			Arch: parts[1],
			OS:   parts[0],
		}
	}
	if username != "" && passwd != "" {
		src.DockerAuth = &types.DockerAuthConfig{
			Username: username,
			Password: passwd,
		}
	}

	return src, nil
}

type ImageNode struct {
	ImageKey     string
	Platform     *Platform
	DockerAuth   *types.DockerAuthConfig
	AuthFilePath string
}

type Platform struct {
	Arch string
	OS   string
}

func (n ImageNode) ToSystemContext() *types.SystemContext {
	c := &types.SystemContext{
		OCIInsecureSkipTLSVerify:    true,
		DockerInsecureSkipTLSVerify: types.NewOptionalBool(true),
	}
	if n.Platform != nil {
		c.ArchitectureChoice = n.Platform.Arch
		c.OSChoice = n.Platform.OS
	}
	if n.DockerAuth != nil {
		c.DockerAuthConfig = n.DockerAuth
	}
	c.AuthFilePath = n.AuthFilePath

	return c
}

func copyImage(ctx context.Context, src, dst ImageNode) error {
	logger.Infof("copy image from %s to %s", src.ImageKey, dst.ImageKey)
	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return err
	}
	defer policyContext.Destroy()

	// 解析源镜像和目标 tar 包路径
	srcRef, err := alltransports.ParseImageName(src.ImageKey)
	if err != nil {
		return fmt.Errorf("invalid source image name %s: %v", src.ImageKey, err)
	}
	destRef, err := alltransports.ParseImageName(dst.ImageKey)
	if err != nil {
		return fmt.Errorf("invalid destination tar file %s: %v", dst.ImageKey, err)
	}

	// 执行镜像拷贝
	_, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
		ReportWriter:       os.Stdout,
		SourceCtx:          src.ToSystemContext(),
		DestinationCtx:     dst.ToSystemContext(),
		ImageListSelection: copy.CopySystemImage,
	})
	if err != nil {
		return fmt.Errorf("error copying image: %v", err)
	}

	return nil
}

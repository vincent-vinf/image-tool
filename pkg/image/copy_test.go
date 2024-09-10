package image

import (
	"context"
	"testing"
)

func TestCheck(t *testing.T) {
	src := ImageNode{ImageKey: "docker-archive:/Users/vincent/Downloads/test-images/registry.cn-hangzhou.aliyuncs.com#qkd-system#workflow-controller+v3.5.3.tar"}
	dst := ImageNode{ImageKey: "docker-daemon:registry.cn-hangzhou.aliyuncs.com/qkd-system/workflow-controller:v3.5.3"}
	err := copyImage(context.Background(), src, dst)
	if err != nil {
		t.Fatal(err)
	}
}

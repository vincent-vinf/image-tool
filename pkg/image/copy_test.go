package image

import (
	"context"
	"testing"
)

func TestCheck(t *testing.T) {
	err := check(context.Background(), "/Users/vincent/Downloads/test-images/registry.cn-hangzhou.aliyuncs.com#qkd-system#tensorboard+v0.1.0.tar")
	if err != nil {
		t.Fatal(err)
	}
}

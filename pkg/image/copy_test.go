package image

import (
	"context"
	"testing"
)

func TestCopy(t *testing.T) {
	src, err := NewRegistryImageNode("docker://docker.io/calico/node:master", "linux/arm64", "", "")
	if err != nil {
		t.Fatal(err)
	}
	dst := ImageNode{ImageKey: "docker-archive:/Users/vincent/Downloads/docker.io#calico#node+v3.20.6.tar"}
	err = copyImage(context.Background(), src, dst)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheck(t *testing.T) {
	err := CheckImageTar(context.Background(), "/Users/vincent/Downloads/docker.io#calico#node+v3.20.6.tar", "linux", "amd64")
	if err != nil {
		t.Fatal(err)
	}
}

package zip

import "testing"

func TestZip(t *testing.T) {
	err := Compress("/Users/vincent/Downloads/registry.cn-hangzhou.aliyuncs.com_qkd-system_ai-stack-gateway_2.5.0_0802.tar", "/Users/vincent/Downloads/t.gz", Gzip)
	if err != nil {
		t.Fatal(err)
	}
	err = Decompress("/Users/vincent/Downloads/t.gz", "/Users/vincent/Downloads/t.tar", Gzip)
	if err != nil {
		t.Fatal(err)
	}
}

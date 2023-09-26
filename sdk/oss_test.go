package sdk

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	client := New("127.0.0.1:8000")
	url, err := client.GenerateDownloadURL("/dir/test", WithGenerateDownloadURLOptionsExt("png"), WithGenerateDownloadURLOptionsExpire("100"), WithGenerateDownloadURLOptionsFilename("test-01"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)
}

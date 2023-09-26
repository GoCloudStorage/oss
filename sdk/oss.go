package sdk

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	addr string
}

func New(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) GenerateDownloadURL(key string, opFn ...generateDownloadURLOptionFunc) (string, error) {
	var (
		config generateDownloadURLOption
		mp     = make(map[string]interface{})
	)
	for _, fn := range opFn {
		fn(&config)
	}
	result, err := url.JoinPath("http://", c.addr, "/url/generate_download_url")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", result, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("OSS-Key", key)
	req.Header.Set("OSS-Filename", config.Filename)
	req.Header.Set("OSS-Ext", config.Ext)
	req.Header.Set("OSS-Expire", config.Expire)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &mp); err != nil {
		return "", err
	}
	return mp["data"].(map[string]interface{})["url"].(string), nil
}

type generateDownloadURLOption struct {
	Expire   string
	Key      string
	Filename string
	Ext      string
}
type generateDownloadURLOptionFunc func(*generateDownloadURLOption)

func WithGenerateDownloadURLOptionsExpire(Expire string) func(*generateDownloadURLOption) {
	return func(option *generateDownloadURLOption) {
		option.Expire = Expire
	}
}

func WithGenerateDownloadURLOptionsFilename(Filename string) func(option *generateDownloadURLOption) {
	return func(option *generateDownloadURLOption) {
		option.Filename = Filename
	}
}

func WithGenerateDownloadURLOptionsExt(Ext string) func(option *generateDownloadURLOption) {
	return func(option *generateDownloadURLOption) {
		option.Ext = Ext
	}
}

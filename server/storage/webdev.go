package storage

import (
	"github.com/studio-b12/gowebdav"
	"github.com/tongruirenye/OrgICSX5/server/config"
)

type WebDevClient struct {
	url      string
	user     string
	password string
	client   *gowebdav.Client
}

var AppStorage *WebDevClient

func InitStorage() {
	if AppStorage != nil {
		return
	}

	AppStorage = NewWebDevClient(config.AppConfig.WebDevRoot,
		config.AppConfig.WebDevUser,
		config.AppConfig.WebDevPass)
}

func NewWebDevClient(url, user, password string) *WebDevClient {
	c := &WebDevClient{
		url:      url,
		user:     user,
		password: password,
	}

	c.client = gowebdav.NewClient(url, user, password)
	return c
}

func (c *WebDevClient) ListFileList(path string) ([]string, error) {
	files, err := c.client.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var f []string
	for _, file := range files {
		f = append(f, file.Name())
	}

	return f, nil
}

func (c *WebDevClient) ReadFile(file string) ([]byte, error) {
	return c.client.Read(file)
}

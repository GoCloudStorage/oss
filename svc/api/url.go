package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"oss/pkg/response"
	"oss/svc/service"
	"strconv"
	"time"
)

func GenerateDownloadUrl(c *fiber.Ctx) error {
	var (
		key      string
		filename string
		ext      string
		expire   time.Duration
	)

	key = c.Get("OSS-Key", "")
	if key == "" {
		return response.Resp400(c, nil, "OSS-Key not is nil")
	}

	if !service.IsExistByKey(key) {
		return response.Resp400(c, nil, "object not found")
	}

	filename = c.Get("OSS-Filename", "default")
	ext = c.Get("OSS-Ext", "txt")

	oe := c.Get("OSS-Expire", "")
	if oe == "" {
		expire = time.Minute * 30
	} else {
		t, err := strconv.Atoi(oe)
		if err != nil {
			return response.Resp400(c, nil, "OSS-Expire format incorrect, base on second")
		}
		expire = time.Second * time.Duration(t)
	}

	code, downloadToken, err := service.GenerateDownloadURL(key, filename, ext, expire)
	if err != nil {
		return response.Resp500(c, nil)
	}
	return response.Resp200(c, map[string]interface{}{
		"url": fmt.Sprintf("http://localhost:8000/pre_sign/%s?token=%s", code, downloadToken),
	}, "success")
}

func GenerateUploadUrl(c *fiber.Ctx) error {
	panic("not impl")
}

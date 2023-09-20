package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"oss/pkg/db/redis"
	"oss/pkg/random"
	"oss/pkg/response"
	"oss/pkg/token"
	"oss/svc/model"
	"strconv"
	"time"
)

func GenerateDownloadUrl(c *fiber.Ctx) error {
	var (
		key      string
		filename string
		ext      string
		expire   time.Duration
		object   model.Object
	)

	key = c.Get("OSS-Key", "")
	if key == "" {
		return response.Resp400(c, nil, "OSS-Key not is nil")
	}

	if !object.IsExistByKey(key) {
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

	// 设置 code: key映射
	code := random.GenerateRandomString(128)
	if err := redis.SetEx(context.Background(), getDownloadCode(code), key, expire); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("can`t set code:key mapping, err: %v", err))
	}

	// 生成token
	downloadToken, err := token.GenerateDownloadToken(key, filename, ext, expire)
	if err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("generate token fail, err: %v", err))
	}
	return response.Resp200(c, map[string]interface{}{
		"url": fmt.Sprintf("http://localhost:8000/pre_sign/%s?token=%s", code, downloadToken),
	}, "success")
}

func GenerateUploadUrl(c *fiber.Ctx) error {
	panic("not impl")
}

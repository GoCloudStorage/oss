package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"oss/pkg/db/redis"
	"oss/pkg/response"
	"oss/pkg/token"
	"oss/svc/model"
)

func Download(c *fiber.Ctx) error {
	var (
		object model.Object
		code   string
	)
	code = c.Params("code")
	t := c.Query("token")

	downloadToken, err := token.ParseDownloadToken(t)
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}

	key, err := redis.Get(context.Background(), getDownloadCode(code))
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}

	if key != downloadToken.Key {
		fmt.Println(key, downloadToken.Key)
		return response.Resp403(c, nil)
	}

	if err := object.GetObjectByKey(key); err != nil {
		return response.Resp500(c, nil, "not found object")
	}

	// 设置响应头
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fmt.Sprintf("%s.%s", downloadToken.Filename, downloadToken.Ext)))
	c.Set("Content-Type", "application/"+downloadToken.Ext)

	return c.SendFile(object.Path)
}

func Upload(c *fiber.Ctx) error {
	panic("not impl")
}

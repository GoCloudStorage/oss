package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"oss/pkg/response"
	"oss/pkg/token"
	"oss/svc/service"
)

func Download(c *fiber.Ctx) error {
	var (
		code string
	)
	code = c.Params("code")
	t := c.Query("token")

	downloadToken, err := token.ParseDownloadToken(t)
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}

	ds := service.DownloadService{}
	key, err := ds.GetKeyByCode(code)
	if err != nil {
		return response.Resp400(c, nil)
	}
	if key != downloadToken.Key {
		return response.Resp403(c, nil)
	}
	path, err := ds.GetPath(key)
	if err != nil {
		return response.Resp500(c, nil)
	}

	// 设置响应头
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fmt.Sprintf("%s.%s", downloadToken.Filename, downloadToken.Ext)))
	c.Set("Content-Type", "application/"+downloadToken.Ext)

	return c.SendFile(path)
}

func Upload(c *fiber.Ctx) error {
	panic("not impl")
}

package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"oss/pkg/response"
	"oss/pkg/storage"
	"oss/svc/model"
)

func AbortMultipartUpload(c *fiber.Ctx) error {
	var (
		object model.Object
	)
	key := c.Get("OSS-Key", "")
	if key == "" {
		return response.Resp400(c, nil, "OSS-Key not is nil")
	}

	if err := object.DeleteByKey(key); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("delete object fail, err: %v", err))
	}

	if err := storage.Client.Remove(key); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("remove object fail, err: %v", err))
	}

	return response.Resp200(c, nil, "success")
}

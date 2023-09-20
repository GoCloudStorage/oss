package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"oss/pkg/response"
	"oss/pkg/storage"
	"oss/svc/model"
)

func CompleteMultipartUpload(c *fiber.Ctx) error {
	var (
		object model.Object
	)
	key := c.Get("OSS-Key", "")
	if key == "" {
		return response.Resp400(c, nil, "OSS-Key not is nil")
	}

	if !object.IsExistByKey(key) {
		return response.Resp400(c, nil, fmt.Sprintf("object not found, %s", key))
	}

	if err := object.GetObjectByKey(key); err != nil {
		return response.Resp500(c, nil, err.Error())
	}

	if object.Size != object.TotalSize {
		msg := fmt.Sprintf("object not complete, size: %v, total size: %v", object.Size, object.TotalSize)
		return response.Resp400(c, nil, msg)
	}

	if !object.IsComplete {
		path, err := storage.Client.MergeChunk(key, object.TotalSize)
		if err != nil {
			return response.Resp500(c, nil, err.Error())
		}
		object.Path = path
		object.IsComplete = true
	}

	if err := object.Update(); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("save object record fail, err: %v", err))
	}

	return response.Resp200(c, nil, "success")
}

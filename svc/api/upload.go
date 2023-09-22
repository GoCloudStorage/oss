package api

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"oss/pkg/response"
	"oss/pkg/storage"
	"oss/svc/service"
)

// 分块上传
func UploadPart(c *fiber.Ctx) error {

	// 解析请求头
	uploadReq, err := parasUploadHeader(c)
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}

	object, err := service.UploadPart(bytes.NewReader(c.Body()), uploadReq)
	if err != nil {
		return response.Resp500(c, nil)
	}

	return response.Resp200(c, object, "success")
}

// PutObject 上传一个文件
func PutObject(c *fiber.Ctx) error {

	// 解析请求头
	uploadReq, err := parasUploadCommonHeader(c)
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}

	object, err := service.UploadPart(bytes.NewReader(c.Body()), uploadReq)
	if err != nil {
		return response.Resp500(c, nil)
	}

	// merge object
	if object.Size == uploadReq.ContentRange.Total {
		path, err := storage.Client.MergeChunk(object.Key, object.TotalSize)
		if err != nil {
			return response.Resp500(c, nil, fmt.Sprintf("merge chunk failed, err: %v", err))
		}
		object.Path = path
		object.IsComplete = true
	}

	if err = object.Update(); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("save object record failed, err: %v", err))
	}

	return response.Resp200(c, object, "success")
}

// CompleteMultipartUpload 合并分块
func CompleteMultipartUpload(c *fiber.Ctx) error {

	key := c.Get("OSS-Key", "")
	if key == "" {
		return response.Resp400(c, nil, "OSS-Key not is nil")
	}

	flag := service.IsExistByKey(key)
	if !flag {
		return response.Resp400(c, nil)
	}

	err := service.MergeMultipartUpload(key)
	if err != nil {
		return response.Resp500(c, nil)
	}

	return response.Resp200(c, nil, "success")
}

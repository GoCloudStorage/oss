package controller

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"oss/pkg/response"
	"oss/pkg/storage"
	"oss/svc/model"
)

func UploadPart(c *fiber.Ctx) error {
	var (
		object model.Object
		err    error
	)

	// 解析请求头
	uploadReq, err := parasUploadHeader(c)
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}

	// 获取 object record
	if object.IsExistByKey(uploadReq.Key) {
		err := object.GetObjectByKey(uploadReq.Key)
		if err != nil {
			return response.Resp500(c, nil, fmt.Sprintf("failed get object, err: %v", err))
		}
	} else {
		object = model.Object{
			Key:       uploadReq.Key,
			MD5:       uploadReq.MD5,
			Type:      "Multipart",
			TotalSize: uploadReq.TotalSize,
		}
		if err := object.Create(); err != nil {
			return response.Resp500(c, nil, fmt.Sprintf("failed create object record, err: %v", err))
		}
	}

	// save object data
	f := bytes.NewReader(c.Body())
	if err = storage.Client.SaveChunk(object.Key, 0, f, int64(uploadReq.ContentRange.start)); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("failed save chunk, err: %v", err))
	}
	object.Size += f.Len()

	if err = object.Update(); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("save object record failed, err: %v", err))
	}

	return response.Resp200(c, object, "success")
}
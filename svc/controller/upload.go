package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	response "oss/pkg/reponse"
	"oss/pkg/storage"
	"oss/svc/model"
	"strconv"
	"strings"
)

func PutObject(c *fiber.Ctx) error {
	var (
		object model.Object
		err    error
		cr     contentRange // Content-Range 参数
	)

	// 获取Content-Range
	r := c.Get("Content-Range", "nil")
	if r == "nil" { // 没有断点续传, 覆盖上传
		return response.Resp400(c, nil, "Content-Range not set")
	}
	cr, err = convertContentRange(r)
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}

	// 获取object md5
	md5 := c.FormValue("md5", "nil")
	if md5 == "nil" {
		return response.Resp400(c, nil, "md5 form not found")
	}
	object.MD5 = md5

	// 获取key
	key := c.FormValue("key", "")
	if key == "" {
		return response.Resp400(c, nil, "key form not is nil")
	}
	object.Key = key

	// 获取object
	fh, err := c.FormFile("object")
	if err != nil {
		return response.Resp400(c, nil, "failed to translate object")
	}

	// 获取 object record
	if object.IsExistByKey(key) {
		err := object.GetObjectByKey(key)
		if err != nil {
			return response.Resp500(c, nil, fmt.Sprintf("failed get object, err: %v", err))
		}
	} else {
		object.TotalSize = cr.total
		object.Type = "Normal"
		if err := object.Create(); err != nil {
			return response.Resp500(c, nil, fmt.Sprintf("failed create object record, err: %v", err))
		}
	}
	logrus.Info(cr, object)
	// 校验Content-Range是否正确
	if cr.start < object.Size || cr.end > object.TotalSize || cr.total > object.TotalSize {
		return response.Resp400(c, nil, fmt.Sprintf("Content-Range range is incorrect"))
	}

	f, err := fh.Open()
	if err != nil {
		return response.Resp500(c, nil, "open file failed")
	}
	defer f.Close()

	if err = storage.Client.SaveChunk(key, 0, f, int64(object.Size)); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("failed save chunk, err: %v", err))
	}
	object.Size += int(fh.Size)

	if err = storage.Client.MergeChunk(key, 1, int(object.Size)); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("merge chunk failed, err: %v", err))
	}
	if err = object.Update(); err != nil {
		return response.Resp500(c, nil, fmt.Sprintf("save object record failed, err: %v", err))
	}
	return response.Resp200(c, object, "success")
}

func UploadPart(c *fiber.Ctx) error {
	panic("not impl")
}

func AbortMultipartUpload(c *fiber.Ctx) error {
	panic("not impl")
}

func CompleteMultipartUpload(c *fiber.Ctx) error {
	panic("not impl")
}

type contentRange struct {
	start int
	end   int
	total int
}

func convertContentRange(d string) (res contentRange, err error) {
	t := strings.Split(d, " ")
	if len(t) != 2 {
		return res, fmt.Errorf("Content-Range format incorrect")
	}
	t = strings.Split(t[1], "-")
	if len(t) != 2 {
		return res, fmt.Errorf("Content-Range format incorrect")
	}
	res.start, err = strconv.Atoi(t[0])
	if err != nil {
		return res, fmt.Errorf("start convert to int64 incorrect, err: %v", err)
	}
	t = strings.Split(t[1], "/")
	if len(t) != 2 {
		return res, fmt.Errorf("Content-Range format incorrect")
	}
	res.end, err = strconv.Atoi(t[0])
	if err != nil {
		return res, fmt.Errorf("end convert to int64 incorrect, err: %v", err)
	}
	res.total, err = strconv.Atoi(t[1])
	if err != nil {
		return res, fmt.Errorf("total convert to int64 incorrect, err: %v", err)
	}
	return res, nil
}

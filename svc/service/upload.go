package service

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"oss/pkg/response"
	"oss/pkg/storage"
	"oss/svc/model"
)

func UploadPart(f *bytes.Reader, uploadReq *model.UploadReq) (*model.Object, error) {
	var (
		object *model.Object
		err    error
	)
	// 获取 object record
	err, find := object.GetObjectByKey(uploadReq.Key)
	if err != nil {
		logrus.Errorf("failed get object, err: %v", err)
		return nil, err
	}
	if !find {
		object = &model.Object{
			Key:       uploadReq.Key,
			MD5:       uploadReq.MD5,
			Type:      "Multipart",
			TotalSize: uploadReq.TotalSize,
		}
		if err = object.Create(); err != nil {
			logrus.Errorf("failed create object record, err: %v", err)
			return nil, err
		}
	}

	// save object data
	object.Size += f.Len()
	if err = storage.Client.SaveChunk(object.Key, uploadReq.ChunkNumber, f, int64(uploadReq.ContentRange.Start)); err != nil {
		logrus.Errorf("failed save chunk, err: %v", err)
		return nil, err
	}

	if err = object.Update(); err != nil {
		logrus.Errorf("save object record failed, err: %v", err)
		return nil, err
	}
	return object, err
}

func IsExistByKey(key string) bool {
	var (
		object *model.Object
	)
	if !object.IsExistByKey(key) {
		logrus.Warnf("object not found, %s", key)
		return false
	}
	return true
}

func MergeMultipartUpload(key string) error {
	var (
		object *model.Object
	)

	if err, _ := object.GetObjectByKey(key); err != nil {
		logrus.Error("GetObjectByKey err:", err)
		return err
	}

	if object.Size != object.TotalSize {
		logrus.Errorf("object not complete, size: %v, total size: %v", object.Size, object.TotalSize)
		return errors.New(response.MSG400)
	}

	if !object.IsComplete {
		path, err := storage.Client.MergeChunk(key, object.TotalSize)
		if err != nil {
			logrus.Error("merge chunk err:", err)
			return err
		}
		object.Path = path
		object.IsComplete = true
	}

	if err := object.Update(); err != nil {
		logrus.Errorf("save object record fail, err: %v", err)
		return err
	}
	return nil
}

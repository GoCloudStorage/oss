package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"oss/pkg/db/redis"
	"oss/svc/model"
)

type DownloadService struct {
}

func (s *DownloadService) GetPath(key string) (string, error) {
	object := new(model.Object)
	if err, _ := object.GetObjectByKey(key); err != nil {
		return "", err
	}
	return object.Path, nil
}

func (s *DownloadService) GetKeyByCode(code string) (string, error) {
	key, err := redis.Get(context.Background(), getDownloadCode(code))
	if err != nil {
		logrus.Warnf("not found code %s", code)
		return "", err
	}
	return key, nil
}

func getDownloadCode(key string) string {
	return "storage:download:key:" + key
}

package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"oss/pkg/db/redis"
	"oss/pkg/random"
	"oss/pkg/token"
	"oss/svc/model"
	"time"
)

func GetPath(key string) (string, error) {
	object := new(model.Object)
	if err, _ := object.GetObjectByKey(key); err != nil {
		return "", err
	}
	return object.Path, nil
}

func GetKeyByCode(code string) (string, error) {
	key, err := redis.Get(context.Background(), getDownloadCode(code))
	if err != nil {
		logrus.Warnf("not found code %s", code)
		return "", err
	}
	return key, nil
}

func GenerateDownloadURL(key, filename, ext string, expire time.Duration) (string, string, error) {
	// 设置 code: key映射
	code := random.GenerateRandomString(128)
	if err := redis.SetEx(context.Background(), getDownloadCode(code), key, expire); err != nil {
		logrus.Errorf("can`t set code:key mapping, err: %v", err)
		return "", "", err
	}

	// 生成token
	downloadToken, err := token.GenerateDownloadToken(key, filename, ext, expire)
	if err != nil {
		logrus.Errorf("generate token fail, err: %v", err)
		return "", "", err
	}
	return code, downloadToken, nil
}

func getDownloadCode(key string) string {
	return "storage:download:key:" + key
}

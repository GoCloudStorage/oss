package model

import (
	"errors"
	"gorm.io/gorm"
	"oss/pkg/db/pg"
)

type Object struct {
	gorm.Model `gorm:""`
	Key        string `gorm:"not null; unique"` // 对象唯一key
	MD5        string `gorm:"not null"`         // 对象md5
	Type       string `gorm:"not null"`         // 存储类型: Normal, Multipart
	Path       string `gorm:""`                 // 存储路径
	Size       int    `gorm:""`                 // 当前存储大小
	TotalSize  int    `gorm:""`                 // 总文件大小
	IsComplete bool   `gorm:""`                 // 是否完整
	VersionID  uint   `gorm:""`                 // 对象版本号
}

func (o *Object) IsExistByKey(key string) bool {

	tx := pg.Client.Where("key = ?", key).First(o)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// GetObjectByKey 获取对象，返回错误和是否找到
func (o *Object) GetObjectByKey(key string) (error, bool) {
	err := pg.Client.Where("key = ?", key).First(o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	}
	return err, true
}

func (o *Object) Create() error {
	return pg.Client.Create(o).Error
}

func (o *Object) Update() error {
	return pg.Client.Updates(o).Error
}

func (o *Object) DeleteByKey(key string) error {
	return pg.Client.Where("key = ?", key).Unscoped().Delete(o).Error
}

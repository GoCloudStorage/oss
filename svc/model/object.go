package model

import (
	"gorm.io/gorm"
	"oss/pkg/db/pg"
)

type Object struct {
	gorm.Model `gorm:"gorm_._model"`
	Key        string `gorm:"key; not null; unique"` // 对象唯一key
	MD5        string `gorm:"md_5; not null"`        // 对象md5
	Type       string `gorm:"type; not null"`        // 存储类型: Normal, Multipart
	Path       string `gorm:"path"`                  // 存储路径
	Size       int    `gorm:"size"`                  // 当前存储大小
	TotalSize  int    `gorm:"total_size"`            // 总文件大小
	VersionID  uint   `gorm:"version_id"`            // 对象版本号
}

func (o *Object) IsExistByKey(key string) bool {
	var cnt int64
	pg.Client.Model(o).Where("key = ?", key).Count(&cnt)
	if cnt > 0 {
		return true
	}
	return false
}

func (o *Object) GetObjectByKey(key string) error {
	return pg.Client.Where("key = ?", key).First(o).Error
}

func (o *Object) Create() error {
	return pg.Client.Create(o).Error
}

func (o *Object) Update() error {
	return pg.Client.Updates(o).Error
}

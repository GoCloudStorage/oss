package model

import (
	"gorm.io/gorm"
	"oss/pkg/db/pg"
)

type Object struct {
	gorm.Model
	Key       string // 对象唯一key
	MD5       string // 对象md5
	Type      string // 存储类型: Normal, Multipart
	Path      string // 存储路径
	Size      int    // 当前存储大小
	TotalSize int    // 总文件大小
	VersionID uint   // 对象版本号
}

func (o *Object) IsExistByKey(key string) bool {
	var cnt int64
	pg.Client.Where("key = ?", key).Count(&cnt)
	if cnt > 0 {
		return true
	}
	return false
}

func (o *Object) GetObjectByKey(key string) error {
	panic("not impl")
}

func (o *Object) Create() error {
	return pg.Client.Create(o).Error
}

func (o *Object) Update() error {
	return pg.Client.Updates(o).Error
}

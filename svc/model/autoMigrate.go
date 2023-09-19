package model

import "oss/pkg/db/pg"

func Migrate() {
	pg.Client.AutoMigrate(&Object{})
}

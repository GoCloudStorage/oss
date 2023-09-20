package main

import (
	"github.com/gofiber/fiber/v2"
	"oss/pkg/db/pg"
	"oss/pkg/db/redis"
	"oss/pkg/storage"
	"oss/svc/controller"
	"oss/svc/middleware"
	"oss/svc/model"
)

func main() {
	redis.Init("162.14.115.114:6379", "12345678", 0)
	pg.Init("162.14.115.114", "cill", "12345678", "test", "5432")
	model.Migrate()
	// storage init
	storage.Init("/home/qiao/Desktop/buckets")
	// web api init
	app := fiber.New()
	upload := app.Group("/upload")
	upload.Use(middleware.Auth)
	{
		upload.Put("/put_object", controller.PutObject)
		upload.Put("/upload_part", controller.UploadPart)
		upload.Post("/abort_multipart_upload", controller.AbortMultipartUpload)
		upload.Post("/complete_multipart_upload", controller.CompleteMultipartUpload)
	}
	preSign := app.Group("/pre_sign")
	{
		preSign.Get("/:code", controller.Download)
		preSign.Put("/:code", controller.Upload)
	}
	url := app.Group("/url")
	upload.Use(middleware.Auth)
	{
		url.Post("/generate_download_url", controller.GenerateDownloadUrl)
		url.Post("/generate_upload_url", controller.GenerateUploadUrl)
	}

	app.Listen(":8000")
}

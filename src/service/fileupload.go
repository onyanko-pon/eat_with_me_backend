package service

import (
	"io"
	"os"

	"github.com/onyanko-pon/eat_with_me_backend/src/infra_handler"
	"github.com/onyanko-pon/eat_with_me_backend/src/infra_handler/aws"
)

type FileService struct {
	UploadHandler infra_handler.FileUploaderInterface
}

func NewFileService() (*FileService, error) {
	if os.Getenv("GO_ENV") == "production" {
		accessKey, _ := aws.NewAccessKeyFromEnv()
		bucket, _ := aws.NewS3Bucket(aws.AP_NORTHEAST_1, aws.USERICON_BUCKET_NAME, accessKey)
		handler, _ := aws.NewS3Handler(bucket)
		return &FileService{UploadHandler: handler}, nil
	}

	handler := infra_handler.LocalFileUploader{
		RootDir: os.Getenv("ROOT_DIR") + "/images",
	}
	return &FileService{UploadHandler: handler}, nil
}

func (service FileService) Upload(file aws.File) (string, error) {
	url, err := service.UploadHandler.Upload(file)
	return url, err
}

func (service FileService) UploadUserIcon(data io.Reader, filename string) (string, error) {
	file := aws.File{
		Filepath: "usericons",
		Filename: filename,
		Data:     data,
	}
	return service.UploadHandler.Upload(file)
}

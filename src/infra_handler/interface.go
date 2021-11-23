package infra_handler

import (
	myaws "github.com/onyanko-pon/eat_with_me_backend/src/infra_handler/aws"
)

type FileUploaderInterface interface {
	Upload(s3File myaws.File) (url string, err error)
}

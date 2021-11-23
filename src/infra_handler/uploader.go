package infra_handler

import (
	"bytes"
	"fmt"
	"io"
	"os"

	myaws "github.com/onyanko-pon/eat_with_me_backend/src/infra_handler/aws"
)

type LocalFileUploader struct {
	RootDir string
}

func (uploader LocalFileUploader) Upload(file myaws.File) (string, error) {
	key, _ := file.GenKey()
	filepath := uploader.RootDir + "/" + key
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println("error ", filepath, uploader.RootDir)
		return "", err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, file.Data)
	ret := buf.Bytes()

	_, err = f.Write(ret)
	if err != nil {
		return "", err
	}

	return "file://" + filepath, nil
}

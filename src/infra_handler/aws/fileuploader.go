package aws

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Handler struct {
	Bucket   *S3Bucket
	Uploader *s3manager.Uploader
}

type File struct {
	Filepath string
	Filename string
	Data     io.Reader
}

func (file File) GenExtension() string {
	return filepath.Ext(file.Filename)
}

func (file File) GenKey() (string, error) {
	return file.Filepath + "/" + file.Filename, nil
}

func (file File) GenContentType() (string, error) {
	var contentType string

	switch file.GenExtension() {
	case ".jpg":
		contentType = "image/jpeg"
	case ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".png":
		contentType = "image/png"
	case ".ping":
		contentType = "image/png"
	default:
		return "", errors.New("this extension is invalid")
	}
	return contentType, nil
}

func NewS3Handler(bucket *S3Bucket) (*S3Handler, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
				AccessKeyID:     bucket.AccessKey.AccessKeyID,
				SecretAccessKey: bucket.AccessKey.SecretAccessKey,
			}),
			Region: aws.String(bucket.Region),
		},
	}))
	return &S3Handler{
		Bucket:   bucket,
		Uploader: s3manager.NewUploader(sess),
	}, nil
}

func (h *S3Handler) Upload(s3File File) (url string, err error) {

	contentType, _ := s3File.GenContentType()
	key, _ := s3File.GenKey()

	result, err := h.Uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Body:        s3File.Data,
		Bucket:      aws.String(h.Bucket.BucketName),
		ContentType: aws.String(contentType),
		Key:         aws.String(key),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	return result.Location, nil
}

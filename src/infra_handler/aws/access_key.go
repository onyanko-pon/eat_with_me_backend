package aws

import "os"

type AccessKey struct {
	AccessKeyID     string
	SecretAccessKey string
}

func NewAccessKeyFromEnv() (*AccessKey, error) {
	return &AccessKey{
		AccessKeyID:     os.Getenv("AWS_IAM_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_IAM_SECRET_ACCESS_KEY"),
	}, nil
}

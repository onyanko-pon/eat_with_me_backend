package aws

// ap-northeast-1
const AP_NORTHEAST_1 = "ap-northeast-1"
const USERICON_BUCKET_NAME = "eat-with"

type S3Bucket struct {
	Region     string
	BucketName string
	AccessKey  *AccessKey
}

func NewS3Bucket(region string, bucketName string, accessKey *AccessKey) (*S3Bucket, error) {
	bucket := &S3Bucket{
		Region:     region,
		BucketName: bucketName,
		AccessKey:  accessKey,
	}

	return bucket, nil
}

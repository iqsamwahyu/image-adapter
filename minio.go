package image_adapter

import "github.com/minio/minio-go/v7"

type MinioProvider struct {
	minio *minio.Client
}

// Get implements Provider.
func (m MinioProvider) Get(url string) string {
	return url
}

// Upload implements Provider.
func (m MinioProvider) Upload(url string) (string, error) {
	return url, nil
}

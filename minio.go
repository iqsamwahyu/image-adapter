package image_adapter

import "github.com/minio/minio-go/v7"

type minioProvider struct {
	minio *minio.Client
}

// Get implements Provider.
func (m minioProvider) Get(url string) string {
	return url
}

// Upload implements Provider.
func (m minioProvider) Upload(url string) (string, error) {
	return url, nil
}

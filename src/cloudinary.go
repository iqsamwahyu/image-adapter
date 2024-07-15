package image_adapter

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryProvider struct {
	cld            *cloudinary.Cloudinary
	AllowedFormats []string
}

// Get implements Provider.
func (c CloudinaryProvider) Get(bucket, url string) string {
	return url
}

// Upload implements Provider.
func (c CloudinaryProvider) Upload(bucket, fileName, url string) (string, error) {
	var (
		b bool = true
	)

	// remove all extensions from file name to cloudinary, because it will be ex:"image.jpg.jpg"
	// cloudinary has its own special feature that can automatically convert image file just changing the extension
	f := fileName
	if dotIndex := strings.Index(fileName, "."); dotIndex != -1 {
		f = fileName[:dotIndex]
	}

	uploadResult, err := c.cld.Upload.Upload(context.Background(), url, uploader.UploadParams{
		PublicID:       f,
		Folder:         bucket,
		Type:           "private",
		Invalidate:     &b,
		AllowedFormats: c.AllowedFormats,
	})

	if err != nil || uploadResult.Error.Message != "" {
		if err != nil {
			return fileName, err
		} else {
			return fileName, errors.New(uploadResult.Error.Message)
		}
	}

	return fileName, nil
}

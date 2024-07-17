package image_adapter

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryProvider struct {
	cld            *cloudinary.Cloudinary
	allowedFormats []string
}

// Get implements Provider.
func (c cloudinaryProvider) Get(bucket, fileName, transformation string) (string, error) {
	return c.makeStringURL(bucket, fileName, transformation)
}

// Upload implements Provider.
func (c cloudinaryProvider) Upload(bucket, fileName, url string) (string, error) {
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
		AllowedFormats: c.allowedFormats,
	})

	if err != nil {
		return "", err
	}

	if uploadResult.Error.Message != "" {
		return "", errors.New(uploadResult.Error.Message)
	}

	return c.makeStringURL(bucket, fileName, "")
}

func (c cloudinaryProvider) makeStringURL(bucket, fileName, transformation string) (string, error) {
	f := bucket + "/" + fileName

	image, err := c.cld.Image(f)
	if err != nil {
		return "", err
	}
	image.DeliveryType = "private"
	image.Config.URL.Analytics = false
	image.Config.URL.SignURL = true
	if transformation != "" {
		image.Config.URL.SignURL = false
		image.Transformation = transformation
	}

	imageStr, err := image.String()
	if err != nil {
		return "", err
	}

	return imageStr, nil
}

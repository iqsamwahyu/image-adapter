package image_adapter

import (
	"context"
	"errors"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

type Provider interface {
	Get(bucket, url string) string
	Upload(bucket, fileName, url string) (string, error)
}

type ImageAdapter struct {
	opt        Option
	main       string
	cloudinary Provider
	minio      Provider
	// s3         Provider
}

type Option struct {
	AllowedExtensions []string
	// IsPublic          bool
}

type UploadParam struct {
	Bucket   string
	Url      string
	FileName string
}

var optionsDefault = Option{
	AllowedExtensions: []string{"jpg", "png", "jpeg"},
	// IsPublic:          false,
}

func New(opt ...Option) *ImageAdapter {
	i := new(ImageAdapter)
	i.opt = optionsDefault

	for _, o := range opt {
		i.opt = o
	}

	return i
}

func (i *ImageAdapter) WithCloudinary(connectionURL string) *ImageAdapter {
	// initiate cloudinary client
	cld, err := cloudinary.NewFromURL(connectionURL)
	if err != nil {
		log.Panic(err)
	}

	// ping check to cloudinary connection
	ping, err := cld.Admin.Ping(context.Background())
	if ping.Status != "ok" {
		if err == nil {
			log.Panic("couldn't connect to cloudinary")
		}
		log.Panic(err)
	}

	if i.main == "" {
		i.main = "cloudinary"
	}

	i.cloudinary = CloudinaryProvider{
		cld:            cld,
		AllowedFormats: i.opt.AllowedExtensions,
	}

	return i
}

func (i *ImageAdapter) Upload(bucketName, fileName, url string) (string, error) {
	f, err := i.makeFileName(fileName)
	if err != nil {
		return f, err
	}

	if url == "" {
		return f, errors.New("url image is empty")
	}

	if i.main == "" {
		return f, errors.New("main adapter is not set")
	}

	return i.cloudinary.Upload(bucketName, f, url)
	// TODO: Upload to each provider asynchronously using go routine and return the main one

}

// make standardized file name
func (i *ImageAdapter) makeFileName(fileName string) (string, error) {
	// TODO: do validating name
	if fileName == "" {
		return "", errors.New("file name is empty")

	}
	return fileName, nil
}

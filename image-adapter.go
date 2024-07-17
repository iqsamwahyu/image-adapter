package image_adapter

import (
	"context"
	"errors"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

type provider interface {
	Get(bucket, fileName, transformation string) (string, error)
	Upload(bucket, fileName, url string) (string, error)
}

type imageAdapter struct {
	opt        Option
	main       string
	cloudinary provider
	minio      provider
	// s3         provider
}

type Option struct {
	AllowedExtensions []string
	// IsPublic          bool
}

var optionsDefault = Option{
	AllowedExtensions: []string{"jpg", "png", "jpeg"},
	// IsPublic:          false,
}

func New(opt ...Option) *imageAdapter {
	i := new(imageAdapter)
	i.opt = optionsDefault

	for _, o := range opt {
		i.opt = o
	}

	return i
}

func (i *imageAdapter) WithCloudinary(connectionURL string) *imageAdapter {
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

	i.cloudinary = cloudinaryProvider{
		cld:            cld,
		allowedFormats: i.opt.AllowedExtensions,
	}

	return i
}

func (i *imageAdapter) Upload(bucketName, fileName, url string) (string, error) {
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

func (i *imageAdapter) Get(bucketName, fileName, transformation string) (string, error) {
	if bucketName == "" {
		return "", errors.New("bucket name is empty")
	}

	if fileName == "" {
		return "", errors.New("file name is empty")
	}

	switch i.main {
	case "cloudinary":
		return i.cloudinary.Get(bucketName, fileName, transformation)
	default:
		return "", errors.New("main adapter is not set")
	}
}

// make standardized file name
func (i *imageAdapter) makeFileName(fileName string) (string, error) {
	// TODO: do validating name
	if fileName == "" {
		return "", errors.New("file name is empty")

	}
	return fileName, nil
}

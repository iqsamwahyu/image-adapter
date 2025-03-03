package image_adapter

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

var (
	cloudinaryConnectionURL = "<put your cloudinary connection url here>"
)

func TestNew(t *testing.T) {
	type args struct {
		opt []Option
	}
	type testCase struct {
		name     string
		args     args
		expected *ImageAdapter
	}

	tests := []testCase{
		{
			name: "init with empty options",
			args: args{
				opt: []Option{},
			},
			expected: &ImageAdapter{
				opt: optionsDefault,
			},
		},
		{
			name: "init with options",
			args: args{
				opt: []Option{
					{
						AllowedExtensions: []string{"jpg", "png", "jpeg"},
					},
				},
			},
			expected: &ImageAdapter{
				opt: optionsDefault,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.opt...); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("New() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestInitWithCloudinary(t *testing.T) {
	type args struct {
		connectionURL string
	}

	type testCase struct {
		name     string
		args     args
		expected bool
	}

	tests := []testCase{
		{
			name: "empty connection url cloudinary",
			args: args{
				connectionURL: "",
			},
			expected: false,
		},
		{
			name: "invalid connection url cloudinary",
			args: args{
				connectionURL: "cloudinary://421413421423:QRsasdsadas@drasdsadsad",
			},
			expected: false,
		},
		{
			name: "valid connection url cloudinary",
			args: args{
				connectionURL: cloudinaryConnectionURL,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			defer func() {
				t.Logf("with cloudinary: %s, expected: %v", test.args.connectionURL, test.expected)
				r := recover()
				// catch panic
				if r != nil {

					if test.expected {
						t.Errorf("expected WithCloudinary(%s) should panic, but it did not", test.args.connectionURL)
					}
				} else {
					if !test.expected {
						t.Errorf("expected WithCloudinary(%s) should not panic, but it did", test.args.connectionURL)
					}
				}
			}()

			i := New()
			i.WithCloudinary(test.args.connectionURL)

		})

	}
}

func TestUploadOnCloudinary(t *testing.T) {
	type args struct {
		bucketName string
		fileName   string
		url        string
	}

	type testCase struct {
		name          string
		args          args
		expectedError bool
	}

	tests := []testCase{
		{
			name: "[1] upload jpg to cloudinary",
			args: args{
				bucketName: "test-bucket",
				fileName:   "test-file-1.jpg",
				url:        "https://letsenhance.io/static/8f5e523ee6b2479e26ecc91b9c25261e/1015f/MainAfter.jpg",
			},
			expectedError: false,
		},
		{
			name: "[2] upload jpg to cloudinary",
			args: args{
				bucketName: "test-bucket",
				fileName:   "test-file-2.jpg",
				url:        "https://upload.wikimedia.org/wikipedia/commons/b/b2/JPEG_compression_Example.jpg",
			},
			expectedError: false,
		},
		{
			name: "[3] upload png to cloudinary",
			args: args{
				bucketName: "test-bucket",
				fileName:   "test-file-3.png",
				url:        "https://images.rawpixel.com/image_png_800/czNmcy1wcml2YXRlL3Jhd3BpeGVsX2ltYWdlcy93ZWJzaXRlX2NvbnRlbnQvam9iNjgwLTE2Ni1wLWwxZGJ1cTN2LnBuZw.png",
			},
			expectedError: false,
		},
		{
			name: "[4] upload broken image to cloudinary",
			args: args{
				bucketName: "test-bucket",
				fileName:   "test-file-4.png",
				url:        "https://imagddds.rsdsdawpixel.com/image_png_800sadsa/dsdsdadadwqeqdasdsadsadwqdqw.png",
			},
			expectedError: true,
		},
		{
			name: "[5] upload empty file name to cloudinary",
			args: args{
				bucketName: "test-bucket",
				fileName:   "",
				url:        "https://imagddds.rsdsdawpixel.com/image_png_800sadsa/dsdsdadadwqeqdasdsadsadwqdqw.png",
			},
			expectedError: true,
		},
		{
			name: "[6] upload empty url image to cloudinary",
			args: args{
				bucketName: "test-bucket",
				fileName:   "test-file-5.png",
				url:        "",
			},
			expectedError: true,
		},
		{
			name: "[7] upload pdf to cloudinary",
			args: args{
				bucketName: "test-bucket",
				fileName:   "test-file-7.pdf",
				url:        "https://www.electrolux.co.id/globalassets/5-support/manuals/id-id/ewf8005eqwa_user-manual-id-id.pdf",
			},
			expectedError: true,
		},
	}
	i := New()
	i.WithCloudinary(cloudinaryConnectionURL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := fmt.Sprintf("^https://res.cloudinary.com/[a-zA-Z0-9]+/image/private/[a-zA-Z0-9-]+/v[0-9]+/%s/%s$", tt.args.bucketName, tt.args.fileName)
			URLRegex := regexp.MustCompile(p)
			got, err := i.Upload(tt.args.bucketName, tt.args.fileName, tt.args.url)
			t.Log(got, err)
			if (!URLRegex.MatchString(got) && err == nil) || (err != nil && !tt.expectedError) {
				t.Errorf("Upload() = %v, expected fileName %v, error: %v", got, p, err)
			}
		})
	}
}

func TestGetOnCloudinary(t *testing.T) {
	type args struct {
		bucketName     string
		fileName       string
		transformation string
	}

	type testCase struct {
		name          string
		args          args
		expectedError bool
	}

	tests := []testCase{
		{
			name: "[1] get jpg from cloudinary",
			args: args{
				bucketName:     "test-bucket",
				fileName:       "test-file-1.jpg",
				transformation: "",
			},
			expectedError: false,
		},
		{
			name: "[2] get jpg from cloudinary",
			args: args{
				bucketName:     "test-bucket",
				fileName:       "test-file-2.jpeg",
				transformation: "",
			},
			expectedError: false,
		},
		{
			name: "[3] get png from cloudinary",
			args: args{
				bucketName:     "test-bucket",
				fileName:       "test-file-3.png",
				transformation: "",
			},
			expectedError: false,
		},
		{
			name: "[4] get broken image from cloudinary",
			args: args{
				bucketName:     "test-bucket",
				fileName:       "test-file-4.png",
				transformation: "",
			},
			expectedError: false,
		},
		{
			name: "[5] get empty file name from cloudinary",
			args: args{
				bucketName:     "test-bucket",
				fileName:       "",
				transformation: "",
			},
			expectedError: true,
		},
		{
			name: "[6] get empty bucket name from cloudinary",
			args: args{
				bucketName:     "",
				fileName:       "test-file-5.png",
				transformation: "",
			},
			expectedError: true,
		},
		{
			name: "[7] get jpg from cloudinary with transformation",
			args: args{
				bucketName:     "test-bucket",
				fileName:       "test-file-1.jpg",
				transformation: "t_media_lib_thumb",
			},
			expectedError: false,
		},
		{
			name: "[8] get jpg from cloudinary with transformation",
			args: args{
				bucketName:     "test-bucket",
				fileName:       "test-file-2.jpeg",
				transformation: "t_media_lib_thumb",
			},
			expectedError: false,
		},
		{
			name: "[9] get png from cloudinary with transformation",
			args: args{
				bucketName:     "test-bucket",
				fileName:       "test-file-3.png",
				transformation: "t_media_lib_thumb",
			},
			expectedError: false,
		},
	}

	i := New()
	i.WithCloudinary(cloudinaryConnectionURL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.Get(tt.args.bucketName, tt.args.fileName, tt.args.transformation)
			if err != nil && !tt.expectedError {
				t.Errorf("Get() = %v, error: %v", got, err)
			}
			t.Logf("Get() = %v, error: %v", got, err)
		})
	}
}

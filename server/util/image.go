package util

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

func CreateThumbnail(data []byte, ext string) ([]byte, error) {
	if ext == ".jpg" {
		return CreateThumbnailJPG(data)
	} else if ext == ".png" {
		return CreateThumbnailPNG(data)
	} else {
		return nil, fmt.Errorf("Cannot Create Thumbnail form: %s", ext)
	}
}

func LimitSize(reader io.Reader, ext string, maxWidth uint, maxHeight uint) (io.Reader, error) {
	var img image.Image
	var err error
	if strings.ToLower(ext) == ".jpg" {
		img, err = jpeg.Decode(reader)
	} else if strings.ToLower(ext) == ".png" {
		img, err = png.Decode(reader)
	} else {
		return nil, fmt.Errorf("Cannot Create Thumbnail form: %s", ext)
	}

	m := resize.Thumbnail(maxWidth, maxHeight, img, resize.Bicubic)
	if err != nil {
		return nil, err
	}

	buf := *new(bytes.Buffer)
	/*if ext == ".jpg" {
		jpeg.Encode(&buf, m, nil)
	}else if ext == ".png" {
		png.Encode(&buf, m)
	}*/
	jpeg.Encode(&buf, m, nil)
	return &buf, nil
}

func CreateThumbnailPNG(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return createTN(img)
}

func CreateThumbnailJPG(data []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return createTN(img)
}

func createTN(img image.Image) ([]byte, error) {
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Thumbnail(240, 300, img, resize.Bicubic)
	buf := *new(bytes.Buffer)
	jpeg.Encode(&buf, m, nil)

	return buf.Bytes(), nil
}

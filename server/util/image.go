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

//CreateThumbnail takes a jpeg or png and creates a jpg encoded thumbnail
func CreateThumbnail(data []byte, ext string) ([]byte, error) {
	if ext == ".jpg" {
		return createThumbnailJPG(data)
	} else if ext == ".png" {
		return createThumbnailPNG(data)
	} else {
		return nil, fmt.Errorf("cannot create thumbnail form: %s", ext)
	}
}

//LimitSize resizes image to a maximum size in w/h keeping the aspect ratio
//Output format is always JPEG to support compression
func LimitSize(reader io.Reader, ext string, maxWidth uint, maxHeight uint) (io.Reader, error) {
	var img image.Image
	var err error
	if strings.ToLower(ext) == ".jpg" {
		img, err = jpeg.Decode(reader)
	} else if strings.ToLower(ext) == ".png" {
		img, err = png.Decode(reader)
	} else {
		return nil, fmt.Errorf("cannot ceate thumbnail form: %s", ext)
	}

	m := resize.Thumbnail(maxWidth, maxHeight, img, resize.Bicubic)
	if err != nil {
		return nil, err
	}

	buf := *new(bytes.Buffer)
	ence := jpeg.Encode(&buf, m, nil)
	if ence != nil {
		panic(ence)
	}
	return &buf, nil
}

func createThumbnailPNG(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return createTN(img)
}

func createThumbnailJPG(data []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return createTN(img)
}

func createTN(img image.Image) ([]byte, error) {
	//preserve aspect ratio
	m := resize.Thumbnail(240, 300, img, resize.Bicubic)
	buf := *new(bytes.Buffer)
	ence := jpeg.Encode(&buf, m, nil)
	if ence != nil {
		panic(ence)
	}
	return buf.Bytes(), nil
}

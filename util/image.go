package util

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
)

func CreateThumbnail(data []byte, ext string)([]byte, error){
	if ext == ".jpg" {
		return CreateThumbnailJPG(data)
	}else if ext == ".png" {
		return CreateThumbnailPNG(data)
	} else {
		return nil, fmt.Errorf("Cannot Create Thumbnail form: %s", ext)
	}
}

func CreateThumbnailPNG(data []byte) ([]byte, error){
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil{
		return nil, err
	}
	return createTN(img)
}


func CreateThumbnailJPG(data []byte) ([]byte, error){
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil{
		return nil, err
	}
	return createTN(img)
}

func createTN(img image.Image) ([]byte, error) {
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Thumbnail(240, 300, img, resize.Lanczos3)
	buf := *new(bytes.Buffer)
	jpeg.Encode(&buf, m, nil)

	return buf.Bytes(),nil
}
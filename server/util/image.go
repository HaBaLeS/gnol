package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

// LimitSize resizes image to a maximum size in w/h keeping the aspect ratio
// Output format is always JPEG to support compression
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

	if err != nil {
		return nil, err
	}
	m := Thumbnail(maxWidth, maxHeight, img)

	buf := *new(bytes.Buffer)
	ence := jpeg.Encode(&buf, m, nil)
	if ence != nil {
		panic(ence)
	}
	return &buf, nil
}

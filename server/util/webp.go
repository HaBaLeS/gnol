package util

import (
	"golang.org/x/image/webp"
	"image/png"
	"os"
)

func Webp2Png(in, out string) error {

	inf, err := os.Open(in)
	defer inf.Close()
	if err != nil {
		return err
	}

	wpi, err := webp.Decode(inf)
	if err != nil {
		return err
	}

	outf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outf.Close()
	return png.Encode(outf, wpi)

}

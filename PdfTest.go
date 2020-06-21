package main

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main(){

	doc, err := fitz.New("/home/falko/comics/PickmansModelArchive.pdf")
	if err != nil {
		panic(err)
	}

	defer doc.Close()

	tmpDir, err := ioutil.TempDir(os.TempDir(), "fitz")
	if err != nil {
		panic(err)
	}

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		//FIXME add shrinking to max size here

		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join(tmpDir, fmt.Sprintf("test%03d.jpg", n)))
		if err != nil {
			panic(err)
		}



		err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
		if err != nil {
			panic(err)
		}

		f.Close()
	}

}

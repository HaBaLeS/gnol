package main

import (
	"compress/flate"
	"encoding/json"
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/mholt/archiver/v3"
	"github.com/nfnt/resize"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"
)

func (s *Session) convert (args []string, options map[string]string) int {

	if !s.processOptionsAndValidate(args, options) {
		return -1
	}

	doc, err := fitz.New(s.InputFile)
	if err != nil {
		panic(err)
	}
	defer doc.Close()
	s.fillMetaData(doc)


	zip := archiver.NewZip()
	zip.CompressionLevel = flate.NoCompression //Disable compression


	if outZip, err :=  os.Create(args[1]) ;err != nil {
		s.Error("Could not create output file\n %v",err)
		return -1
	} else {
		defer outZip.Close()
		if err := zip.Create(outZip); err != nil {
			s.Error("Could not create output file\n %v",err)
			return -1
		}
	}
	defer zip.Close()

	if s.MetaData.CoverPage > doc.NumPage() {
		s.Error("Coverpage %d higher than number of pages %d", s.MetaData.CoverPage , doc.NumPage())
		return -1
	}

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		img, e1 := doc.Image(n)
		if e1 != nil {
			s.Error("Could not create extract image: \n %v",err)
			return -1
		}

		if n+1 == s.MetaData.CoverPage {
			s.SetCoverImage(img)
		}

		pagename := fmt.Sprintf("page%03d.jpg", n)

		tp := filepath.Join(s.TempDir, pagename)
		f, e2 := os.Create(tp)
		if e2 != nil {
			panic(e2)
		}

		//Some thoughts on image for web https://cloudinary.com/blog/top_10_mistakes_in_handling_website_images_and_how_to_solve_them
		m := resize.Thumbnail(2560, 1440, img, resize.MitchellNetravali)
		e3 := jpeg.Encode(f, m, &jpeg.Options{Quality: jpeg.DefaultQuality}) //fixme add profiles for quality
		if e3 != nil {
			panic(e3)
		}
		f.Sync()
		f.Seek(0,0)

		info, _ := os.Stat(tp)

		zip.Write(archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   info,
				CustomName: pagename,
			},
			ReadCloser: f,
		})
		f.Close()
		s.Log("Writing page: %s (%dx%d) -> (%dx%d)\n", pagename, img.Bounds().Dx(),img.Bounds().Dy(),m.Bounds().Dx(), m.Bounds().Dy() )
	}



	meta, merr := os.Create(path.Join(s.TempDir,"gnol.json"))
	if merr != nil {
		panic(merr)
	}
	defer meta.Close()
	enc := json.NewEncoder(meta)
	encErr := enc.Encode(s.MetaData)
	if encErr != nil {
		panic(encErr)
	}
	meta.Sync()
	meta.Seek(0,0)

	info, _ := os.Stat(path.Join(s.TempDir,"gnol.json"))
	zip.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo:   info,
			CustomName: "gnol.json",
		},
		ReadCloser: meta,
	})

	s.cleanup()

	fi, _ := os.Stat(args[1])
	fi.Size()

	return 0
}
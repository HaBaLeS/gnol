package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/HaBaLeS/gnol/lib"
)

func (s *Session) Convert(args []string, options map[string]string) int {
	if !s.processOptionsAndValidate(args, options) {
		return -1
	}
	s.Log("Loading PDF %s", s.InputFile)

	conv, err := lib.NewPdfToImages(s.InputFile, s.TempDir, 85, false)
	if err != nil {
		panic(err)
	}
	defer conv.Close()

	if s.MetaData.CoverPage > conv.NumPages {
		s.Error("--coverpage %d is higher than total number of pages (%d) in pdf", s.MetaData.CoverPage, s.MetaData.NumPages)
		return -1
	}
	s.Log("Using %s as temporary directory", s.TempDir)

	var pdfTime time.Duration
	var resizeTime time.Duration
	var packageTime time.Duration

	for n := 0; n < conv.NumPages; n++ {
		img, e1 := conv.NextImage()
		if e1 != nil {
			s.Log("Image Extraction Failed")
			panic(e1)
		}

		if n+1 == s.MetaData.CoverPage {
			s.Log("Coverpage as Page %d", s.MetaData.CoverPage)
			s.SetCoverImage(img)
		}

		//s2 := time.Now()
		//Some thoughts on image for web https://cloudinary.com/blog/top_10_mistakes_in_handling_website_images_and_how_to_solve_them
		//pageImage := util.Thumbnail(2560, 1440, img)
		//s.storeAsJpg(n, pageImage)
		//resizeTime += time.Since(s2)
		//s.Log("Writing page: %s (%dx%d) -> (%dx%d)", pagename, img.Bounds().Dx(), img.Bounds().Dy(), m.Bounds().Dx(), m.Bounds().Dy())
		fmt.Printf("%d ", n)
		if n%10 == 0 {
			fmt.Print("\n")
		}
	}

	s.Log("Write gnol.json with Metadata and Cover image")
	meta, merr := os.Create(path.Join(s.TempDir, "gnol.json"))
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
	meta.Seek(0, 0)

	s3 := time.Now()
	if err := s.zipFilesTempDir(); err != nil {
		s.Error("Error writing CBZ %v", err)
	}
	packageTime += time.Since(s3)
	s.cleanup()

	s.Log("Pdf Images: %s, Resize: %s,  Repack %s", pdfTime, resizeTime, packageTime)
	s.Log("Done with Pdf2Cbz")
	return 0
}

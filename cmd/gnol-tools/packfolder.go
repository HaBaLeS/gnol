package main

import (
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
)

var allowedTypes = map[string]struct{}{
	".png":  {},
	".jpg":  {},
	".jpeg": {},
}

func (s *Session) packfolder(args []string, options map[string]string) int {
	if !s.processOptionsAndValidate(args, options) {
		return -1
	}
	s.packInternal() //FIXME return error instead of int

	if !s.DryRun && s.DirectUpload {
		fmt.Printf("Directly Uploading %s\n", s.OutputFile)
		s.InputFile = s.OutputFile
		s.uploadInternal()
		fmt.Printf("Deleting CBZ %s\n", s.OutputFile)
		os.Remove(s.OutputFile)
	}
	return 1
}

func (s *Session) packInternal() int {
	if fi, err := os.Stat(s.InputFile); err != nil || !fi.IsDir() {
		s.Error("File does not exist or is not a directory: %s", s.InputFile)
		return -1
	}

	if _, err := os.Stat(path.Join(s.InputFile, "gnol.json")); err == nil {
		f, err := os.Open(path.Join(s.InputFile, "gnol.json"))
		if err == nil {
			dec := json.NewDecoder(f)
			dex := dec.Decode(s.MetaData)
			if dex != nil {
				s.Error("Error decoding existing json", dex)
			}
		}
	}

	s.Log("Creating: %s", s.OutputFile)

	files := make([]string, 0)
	filepath.Walk(s.InputFile, func(p string, fi fs.FileInfo, err error) error {
		ext := path.Ext(fi.Name())
		if _, ok := allowedTypes[ext]; !ok {
			s.Log("Skipping file: '%s' extension '%s' not supported", fi.Name(), ext)
			return nil
		}
		files = append(files, fi.Name())
		return nil
	})
	sort.Strings(files)

	if s.ListOrder {
		for i, v := range files {
			fmt.Printf("[%d]\t%s\n", i, v)
		}
	}

	if s.DryRun {
		fmt.Printf("DryRun - Exiting before any file are created")
		return 0
	}

	for idx, v := range files {
		img, err := s.LoadImage(path.Join(s.InputFile, v))
		if err != nil {
			s.Error("Could not open Image %s", v)
			continue
		}
		oz := img.Bounds()
		img = resize.Thumbnail(2560, 1440, img, resize.MitchellNetravali)
		s.Log("Resized: (%d,%d) -> (%d,%d)", oz.Dx(), oz.Dy(), img.Bounds().Dx(), img.Bounds().Dy()) //TODO make resizing optional
		err = s.StoreAsJpg(idx, img)
		if err != nil {
			s.Error("Could not store resized Image %s", v)
			continue
		}
		if idx+1 == s.MetaData.CoverPage {
			s.SetCoverImage(img)
		}
		s.MetaData.NumPages++
	}

	s.WriteMetadataJson()

	if err := s.ZipFilesInWorkFolder(); err != nil {
		s.Error("Error writing CBZ %v", err)
	}

	s.cleanup()
	if s.HasErrors {
		s.Error("Job finished with Errors")
		return -1
	}
	return 0
}

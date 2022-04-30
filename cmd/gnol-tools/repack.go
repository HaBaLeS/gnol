package main

import (
	"fmt"
	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v3"
	"io"
	"os"
	"path"
	"sort"
)

func (s *Session) repack(args []string, options map[string]string) int {
	if !s.processOptionsAndValidate(args, options) {
		return -1
	}

	fmt.Printf("Name: %s\n", s.MetaData.Name)
	fmt.Printf("OutFile: %s\n", s.OutputFile)

	f, err := os.Open(s.InputFile)
	if err != nil {
		panic(err)
	}

	eIface, err := archiver.ByHeader(f)
	if err != nil {
		panic(err)
	}

	e, ok := eIface.(archiver.Walker)
	if !ok {
		panic(fmt.Errorf("format specified by source filename is not an extractor format: (%T)", eIface))
	}

	filter := make([]string, 0) //All files to filter against
	err = e.Walk(s.InputFile, func(f archiver.File) error {
		h := f.Header.(zip.FileHeader)
		ident := h.Name
		if _, ok := allowedTypes[path.Ext(ident)]; ok {
			filter = append(filter, ident)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	files := make(map[string]string, 0)

	sort.Strings(filter)
	fmt.Printf("Limiting from %d to %d\n\n", s.From, s.To)
	for i, v := range filter {
		if i >= s.From && i <= s.To {
			if s.Verbose {
				fmt.Printf("[%d]\t%s\n", i, v)
				files[v] = v
			}
		}
	}

	//Exit before stuff gets written to disk
	if s.DryRun {
		fmt.Printf("DryRun - Exiting before any file are created")
		return 0
	}

	workdir, err := os.MkdirTemp("", "repack")
	if err != nil {
		panic(err)
	}

	err = e.Walk(s.InputFile, func(f archiver.File) error {
		h := f.Header.(zip.FileHeader)
		ident := h.Name
		if _, ok := files[ident]; ok {
			out, err := os.Create(path.Join(workdir, f.Name()))
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(out, f.ReadCloser)
			if err != nil {
				panic(err)
			}
			out.Close()
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	s.InputFile = workdir
	s.packInternal()

	if s.DirectUpload {
		fmt.Printf("Directly Uploading %s\n", s.OutputFile)
		s.InputFile = s.OutputFile
		s.uploadInternal()
		fmt.Printf("Deleting CBZ %s\n", s.OutputFile)
		os.Remove(s.OutputFile)
	}

	fmt.Printf("Deleting Workdir %s\n", workdir)
	os.RemoveAll(workdir)

	return -1
}

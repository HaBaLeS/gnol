package session

import (
	"context"
	"fmt"
	"github.com/mholt/archives"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
)

func (s *Session) Repack(args []string, options map[string]string) int {
	if !s.processOptionsAndValidate(args, options) {
		return -1
	}

	fmt.Printf("Name: %s\n", s.MetaData.Name)
	fmt.Printf("OutFile: %s\n", s.OutputFile)

	zfs, err := archives.FileSystem(context.Background(), s.InputFile, nil)
	if err != nil {
		panic(err)
	}

	filter := make([]string, 0) //All files to filter against
	extractError := fs.WalkDir(zfs, ".", func(dirPath string, d fs.DirEntry, err error) error {
		if _, ok := allowedTypes[dirPath]; ok {
			filter = append(filter, dirPath)
		}
		return nil
	})
	if extractError != nil {
		panic(extractError)
	}

	files := make(map[string]string, 0)

	sort.Strings(filter)
	fmt.Printf("Limiting from %d to %d\n\n", s.From, s.To)
	for i, v := range filter {
		if i >= s.From && i <= s.To {
			if s.Verbose {
				fmt.Printf("[%d]\t%s\n", i, v)
			}
			files[v] = v
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

	zfs2, err := archives.FileSystem(context.Background(), s.InputFile, nil)
	if err != nil {
		panic(err)
	}

	extractError2 := fs.WalkDir(zfs2, ".", func(dirPath string, d fs.DirEntry, err error) error {
		ident := dirPath
		if _, ok := files[ident]; ok {
			of := path.Base(d.Name())
			out, err := os.Create(path.Join(workdir, of))
			if err != nil {
				panic(err)
			}
			in, err := os.Open(dirPath)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(out, in)
			if err != nil {
				panic(err)
			}
			out.Close()
		}
		return nil
	})
	if extractError2 != nil {
		panic(extractError2)
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

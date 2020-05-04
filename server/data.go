package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var BASE_PATH = "/home/falko/comics/nsx"
var cl = &ComicList{}
type ComicList struct {
	Comics []Metadata
}

func GetComiList() *ComicList{
	if cl.Comics == nil {
		err := filepath.Walk(BASE_PATH, investigateStructure)
		if err != nil {
			panic(err)
		}
	}
	return cl
}


func investigateStructure (path string, info os.FileInfo, err error) error{
	if strings.HasPrefix(info.Name(),"."){
		fmt.Printf("skipping: %s\n", info.Name())
		return filepath.SkipDir
	}

	if info.IsDir() {
		//fmt.Printf("Path: %s\n", path)
		return nil
	}

	usp, me := NewMetadata(path)
	if usp != nil{
		//unsupported filetype
		return nil
	}

	lr := me.Load()

	force := false
	if lr != nil || force {
		//fmt.Println(err)
		err2 := me.Update()
		if err2 != nil{
			fmt.Printf("Unsupported File: %s\n %v\n", path, err2 )
		}
		me.Save()
	}

	cl.Comics = append(cl.Comics, *me)

	return nil
}



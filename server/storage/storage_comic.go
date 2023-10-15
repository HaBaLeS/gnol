package storage

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/mholt/archiver/v3"
	"io"
	"os"
	"path"
)

func GetPageImage(config *util.ToolConfig, filepath string, comicIdent string, pageNum int) (string, error) {

	comicDir := path.Join(config.TempDirectory, comicIdent)
	fileID := fmt.Sprintf("page-%d", pageNum)
	filename := path.Join(comicDir, fileID+".gnol")

	//TODO add a config parameter to enforce jpeg instead of preserving the original
	cnt := 0
	comic, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	arc, err := archiver.ByHeader(comic)
	if err != nil {
		panic(err)
	}
	wk, ok := arc.(archiver.Walker)
	if !ok {
		panic("Cannot Cast")
	}
	extractError := wk.Walk(filepath, func(f archiver.File) error {
		if !isImageFile(f.Name()) {
			return nil
		}
		if cnt == pageNum {

			//Create dir if not exists
			if _, err := os.Stat(comicDir); os.IsNotExist(err) {
				os.Mkdir(comicDir, os.ModePerm)
			}
			out, cerr := os.Create(filename)
			if cerr != nil {
				panic(cerr)
			}

			ext := path.Ext(f.Name())
			newImg, convErr := util.LimitSize(f.ReadCloser, ext, 2560, 1440)
			if convErr != nil {
				return convErr
			}
			io.Copy(out, newImg)
		}
		cnt++
		return nil
	})

	if extractError != nil {
		return "", extractError
	}

	return filename, nil

}

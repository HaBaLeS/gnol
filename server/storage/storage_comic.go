package storage

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/HaBaLeS/gnol/server/cache"
	"github.com/HaBaLeS/gnol/server/database/dao"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/mholt/archives"
)

type FileStorage struct {
	Cache  *cache.ImageCache
	Config *util.ToolConfig
	Dao    *dao.DAO
}

func (f *FileStorage) FetchImageData(comicId string, num int) ([]byte, error) {

	filePath := f.Dao.ComicFilenameForId(comicId)

	//get file from cache
	var err error
	file, hit := f.Cache.GetFileFromCache(filePath, num)
	if !hit {
		file, err = GetPageImage(f.Config, filePath, comicId, num)
		if err != nil {
			return nil, err
		}
		f.Cache.AddFileToCache(file)
	}

	//as a image-provider module not the cache directly
	img, oerr := os.Open(file)
	if oerr != nil {
		return nil, oerr
	} else {
		defer img.Close()
	}

	data, rerr := ioutil.ReadAll(img)
	if rerr != nil {
		return nil, rerr
	}

	return data, nil
}

func GetPageImage(config *util.ToolConfig, filepath string, comicIdent string, pageNum int) (string, error) {

	comicDir := path.Join(config.TempDirectory, comicIdent)
	fileID := fmt.Sprintf("page-%d", pageNum)
	filename := path.Join(comicDir, fileID+".gnol")

	//TODO add a config parameter to enforce jpeg instead of preserving the original
	cnt := 0

	fsys, err := archives.FileSystem(context.Background(), filepath, nil)
	if err != nil {
		panic(err)
	}

	extractError := fs.WalkDir(fsys, ".", func(dirPath string, d fs.DirEntry, err error) error {
		if !isImageFile(d.Name()) {
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
			ext := path.Ext(d.Name())

			file, err := fsys.Open(dirPath)
			if err != nil {
				return err
			}
			newImg, convErr := util.LimitSize(file, ext, 2560, 1440)
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

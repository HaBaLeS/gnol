package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/boltdb/bolt"
	"github.com/mholt/archiver/v3"
	"path"
	"os"
	"io"
)


type ComicList struct {
	Comics []Metadata
}


type ComicStorage struct {
	bs *BoltStorage
	config *util.ToolConfig
}

func newComicStore(bs *BoltStorage, cfg *util.ToolConfig) *ComicStorage {
	return &ComicStorage{
		bs: bs,
		config: cfg,
	}
}

func (cs *ComicStorage) GetMetadata(id string) (*Metadata, error) {
	return nil, nil
}


func (cs *ComicStorage) GetComiList() *ComicList {
	//panic("Implement me")
	return &ComicList{}
}

func (cs *ComicStorage) GetPageImage(comicID string, pageNum int) (string, error) {

	me, notfound := cs.GetMetadata(comicID)
	if notfound != nil {
		return "", fmt.Errorf("Unknown ComicID: %s", comicID)
	}

	comicDir := path.Join(cs.config.TempDirectory, comicID)
	fileID := fmt.Sprintf("%s-%d", comicID, pageNum)
	filename := path.Join(comicDir, fileID+".gnol")

	//TODO add a config parameter to enforce jpeg instead of preserving the original
	cnt := 0
	extractError := me.arc.Walk(me.FilePath, func(f archiver.File) error {
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


func (cs *ComicStorage) LoadComicMetadata(me *Metadata) error {
	return cs.bs.ReadRaw(func(tx *bolt.Tx) error {
		j := tx.Bucket([]byte("meta")).Get([]byte(me.Id))
		if j == nil {
			return fmt.Errorf("Entity with ID: %s not found", me.Id)
		}
		dec := json.NewDecoder(bytes.NewReader(j))
		return dec.Decode(me)
	})
}

func (cs *ComicStorage) SaveComicMeta(me *Metadata) error {
	return cs.bs.Write(me)
}

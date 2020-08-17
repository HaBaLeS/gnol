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


type MetadataList struct {
	Comics []*Metadata
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

func (cs *ComicStorage) GetMetadata(mid []byte) (*Metadata) {
	me := &Metadata{}
	err := cs.bs.ReadRaw(func(tx *bolt.Tx) error {
		b := tx.Bucket(META_BUCKET)
		k, v := b.Cursor().Seek(mid)
		if k != nil && bytes.Equal(k, mid) {
			der := loadFromJson(me, v)
			return der
		} else {
			return fmt.Errorf("metadata with Id %s not found", mid)
		}
	})
	if err != nil {
		fmt.Printf("Error Loading Metadata. %s\n", err)
		return nil
	}
	return me
}

func (cs *ComicStorage) GetPageImage(comicID string, pageNum int) (string, error) {

	me := cs.GetMetadata([]byte(comicID))
	if me == nil {
		return "", fmt.Errorf("Unknown ComicID: %s", comicID)
	}

	comicDir := path.Join(cs.config.TempDirectory, comicID)
	fileID := fmt.Sprintf("%s-%d", comicID, pageNum)
	filename := path.Join(comicDir, fileID+".gnol")

	//TODO add a config parameter to enforce jpeg instead of preserving the original
	cnt := 0
	extractError := me.arc().Walk(me.FilePath, func(f archiver.File) error {
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
		j := tx.Bucket([]byte("meta")).Get(me.IdBytes())
		if j == nil {
			return fmt.Errorf("Entity with ID: %s not found", me.IdBytes())
		}
		dec := json.NewDecoder(bytes.NewReader(j))
		return dec.Decode(me)
	})
}

func (cs *ComicStorage) SaveComicMeta(me *Metadata) error {
	return cs.bs.Write(me)
}

func (cs *ComicStorage) MetadataForList(list []string) *MetadataList {
	ml := &MetadataList{
		Comics: make([]*Metadata,len(list)),
	}
	for i,v := range list {
		ml.Comics[i] = cs.GetMetadata([]byte(v))
	}
	return ml
}

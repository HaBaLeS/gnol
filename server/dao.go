package server

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type DAOHandler struct {
	metaDB    map[string]*Metadata
	comicList *ComicList
	session   Session
}

type ComicList struct {
	Comics []Metadata
}

func NewDAO(session Session) *DAOHandler {
	return &DAOHandler{
		metaDB:    make(map[string]*Metadata),
		comicList: &ComicList{},
		session:   session,
	}
}

func (dao *DAOHandler) getMetadata(id string) (*Metadata, error) {
	m, ok := dao.metaDB[id]
	if !ok {
		return nil, fmt.Errorf("Could find Metadata for: %s", id)
	}
	return m, nil
}

func (dao *DAOHandler) GetComiList() *ComicList {
	return dao.comicList
}

func (dao *DAOHandler) Warmup() {
	dao.session.logger.Info("Reading Data Directory, warmup results")
	err := filepath.Walk(dao.session.config.DataDirectory, dao.investigateStructure)
	if err != nil {
		panic(err)
	}
}

func (dao *DAOHandler) getPageImage(id string, imageNum string) ([]byte, error) {
	//	r := rand.Intn(len(dao.comicList.Comics)-1)
	//	return base64.StdEncoding.DecodeString(dao.comicList.Comics[r].CoverImageBase64)
	m := dao.metaDB[id]
	cnt := 0
	var data []byte
	wr := m.arc.Walk(m.FilePath, func(f archiver.File) error {
		if !isImageFile(f.Name()) {
			return nil
		}
		num, err := strconv.Atoi(imageNum)
		if err != nil {
			return err
		}

		if cnt == num {
			data, err = ioutil.ReadAll(f.ReadCloser)
			if err != nil {
				return err
			}
		}
		cnt++
		return nil
	})

	if wr != nil {
		return nil, wr
	}
	return data, nil
}

func (dao *DAOHandler) investigateStructure(path string, info os.FileInfo, err error) error {
	/*if strings.HasPrefix(info.Name(), ".") && !info.IsDir() {
		return filepath.SkipDir
	}*/

	if info.IsDir() {
		//fmt.Printf("Path: %s\n", path)
		return nil
	}

	usp, me := NewMetadata(path)
	if usp != nil {
		//unsupported filetype
		return nil
	}

	lr := me.Load()

	force := false
	if lr != nil || force {
		//fmt.Println(err)
		err2 := me.Update()
		if err2 != nil {
			fmt.Printf("Unsupported File: %s\n %v\n", path, err2)
		}
		me.Save()
	}

	dao.comicList.Comics = append(dao.comicList.Comics, *me)
	dao.metaDB[me.Id] = me

	return nil
}

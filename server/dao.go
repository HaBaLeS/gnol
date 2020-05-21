package server

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var BASE_PATH = "/home/falko/comics/"

type DAOHandler struct {
	metaCache map[string]*Metadata
	comicList *ComicList
}

type ComicList struct {
	Comics []Metadata
}

func NewDAO() *DAOHandler{
	return &DAOHandler{
		metaCache: make(map[string]*Metadata),
		comicList: &ComicList{},
	}
}

func (dao *DAOHandler) getMetadata(id string) (*Metadata, error){
	m, ok := dao.metaCache[id]
	if !ok {
		return nil, fmt.Errorf("Could find Metadata for: %s", id)
	}
	return m,nil
}


func (dao *DAOHandler) GetComiList() *ComicList{
	return dao.comicList
}

func (dao *DAOHandler) Warmup() {
	err := filepath.Walk(BASE_PATH, dao.investigateStructure)
	if err != nil {
		panic(err)
	}
}

func (dao *DAOHandler) getPageImage(id string, imageNum string) ([]byte, error) {
//	r := rand.Intn(len(dao.comicList.Comics)-1)
//	return base64.StdEncoding.DecodeString(dao.comicList.Comics[r].CoverImageBase64)
	m := dao.metaCache[id]
	cnt := 0
	var data []byte
	wr := m.arc.Walk(m.FilePath, func(f archiver.File) error {
		if !isImageFile(f.Name()){
			return nil
		}
		num,  err :=strconv.Atoi(imageNum)
		if err != nil {
			return  err
		}
		
		if cnt == num  {
			data , err = ioutil.ReadAll(f.ReadCloser)
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
	return  data, nil
}

func (dao *DAOHandler) investigateStructure (path string, info os.FileInfo, err error) error{
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

	dao.comicList.Comics = append(dao.comicList.Comics, *me)
	dao.metaCache[me.Id] = me

	return nil
}
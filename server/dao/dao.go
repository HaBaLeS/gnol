package dao

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/HaBaLeS/go-logger"
	"github.com/mholt/archiver/v3"
	"io"
	"os"
	"path"
	"path/filepath"
)

type DAOHandler struct {
	metaDB    map[string]*Metadata
	comicList *ComicList
	logger    *logger.Logger
	config    *util.ToolConfig
}

type ComicList struct {
	Comics []Metadata
}

func NewDAO(logger *logger.Logger, config *util.ToolConfig) *DAOHandler {
	return &DAOHandler{
		metaDB:    make(map[string]*Metadata),
		comicList: &ComicList{},
		logger:    logger,
		config:    config,
	}
}

func (dao *DAOHandler) GetMetadata(id string) (*Metadata, error) {
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
	dao.logger.Info("Reading Data Directory, warmup results")
	err := filepath.Walk(dao.config.DataDirectory, dao.investigateStructure)
	if err != nil {
		panic(err)
	}
}

func (dao *DAOHandler) GetPageImage(comicID string, pageNum int) (string, error) {

	me, notfound := dao.GetMetadata(comicID)
	if notfound != nil {
		return "", fmt.Errorf("Unknown ComicID: %s", comicID)
	}

	comicDir := path.Join(dao.config.TempDirectory, comicID)
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

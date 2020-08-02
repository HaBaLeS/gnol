package dao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/HaBaLeS/go-logger"
	"github.com/boltdb/bolt"
	"github.com/mholt/archiver/v3"
	"github.com/rs/xid"
	"golang.org/x/crypto/argon2"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

type DAOHandler struct {
	metaDB    map[string]*Metadata
	comicList *ComicList
	logger    *logger.Logger
	config    *util.ToolConfig
	Db        *bolt.DB
}

type ComicList struct {
	Comics []Metadata
}

func NewDAO(logger *logger.Logger, config *util.ToolConfig) *DAOHandler {

	db, err := bolt.Open(config.Database, 0600, &bolt.Options{Timeout: 1 * time.Second})
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("meta"))
		tx.CreateBucketIfNotExists([]byte("jobs_open"))
		tx.CreateBucketIfNotExists([]byte("jobs_error"))
		tx.CreateBucketIfNotExists([]byte("jobs_done"))
		tx.CreateBucketIfNotExists(USER_BUCKET)

		return nil
	})
	if err != nil {
		panic(err)
	}

	return &DAOHandler{
		metaDB:    make(map[string]*Metadata),
		comicList: &ComicList{},
		logger:    logger,
		config:    config,
		Db:        db,
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
	//needs to be replaced with load from DB + checkking if file still exists
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

	if info.IsDir() {
		//fmt.Printf("Path: %s\n", path)
		return nil
	}

	usp, me := NewMetadata(path)
	if usp != nil {
		//unsupported filetype
		return nil
	}

	mer := dao.ComicMetadata(me)

	force := false
	if mer != nil || force {
		//fmt.Println(err)
		err2 := me.UpdateMeta()
		if err2 != nil {
			fmt.Printf("Unsupported File: %s\n %v\n", path, err2)
		}
		dao.SaveComicMeta(me)
	}
	dao.AddComicToList(me)

	return nil
}

func (dao *DAOHandler) AddComicToList(me *Metadata) {
	if me.Public {
		dao.comicList.Comics = append(dao.comicList.Comics, *me)
		dao.metaDB[me.Id] = me
	}
}

func (dao *DAOHandler) Write(bucket []byte, data Entity) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		buf := new(bytes.Buffer)
		e := json.NewEncoder(buf)
		ece := e.Encode(data)
		if ece != nil {
			return ece
		}
		b.Put(data.IdBytes(), buf.Bytes())
		return nil
	})
}

func (dao *DAOHandler) Close() {
	dao.Db.Close()
}

func (dao *DAOHandler) ComicMetadata(me *Metadata) error {
	return dao.Db.View(func(tx *bolt.Tx) error {
		j := tx.Bucket([]byte("meta")).Get([]byte(me.Id))
		if j == nil {
			return fmt.Errorf("Entity with ID: %s not found", me.Id)
		}
		dec := json.NewDecoder(bytes.NewReader(j))
		return dec.Decode(me)
	})
}

func (dao *DAOHandler) SaveComicMeta(me *Metadata) error {
	return dao.Write(META_BUCKET, me)
}

func (dao *DAOHandler) CreateUser(name string, pass string) *User {
	//FIXME users can just override each others
	hash, salt := hashPassword(pass)
	u := &User{
		BaseEntity: CreateBaseEntity(),
		Name:       name,
		PwdHash:    hash,
		Salt:       salt,
	}
	u.Id = name + u.Id
	dao.Write(USER_BUCKET, u)
	return u
}

func (dao *DAOHandler) AuthUser(name string, pass string) (*User, error) {
	u := new(User)
	logError := dao.Db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(USER_BUCKET).Cursor()
		spx := []byte(name)
		for k, v := c.Seek(spx); k != nil && bytes.HasPrefix(k, spx); k, v = c.Next() {
			err := loadFromJson(u, v)
			if err != nil {
				return err
			}
			if u.Name == name {
				return checkPassword(u.Salt, u.PwdHash, pass)
			}
		}
		fmt.Printf("Did not find: %s in DB", name)
		return fmt.Errorf("Login failed")
	})
	return u, logError
}

func loadFromJson(i interface{}, v []byte) error {
	d := json.NewDecoder(bytes.NewReader(v))
	return d.Decode(i)
}

func hashPassword(pass string) ([]byte, []byte) {
	salt := xid.New().Bytes()
	hash := argon2.Key([]byte(pass), salt, 3, 32*1024, 4, 32)
	return hash, salt
}

func checkPassword(salt []byte, dbhash []byte, pass string) error {
	hash := argon2.Key([]byte(pass), salt, 3, 32*1024, 4, 32)
	if bytes.Compare(hash, dbhash) != 0 {
		fmt.Println("Password do not match")
		return fmt.Errorf("Login Error")
	}
	return nil
}

func CreateBaseEntity() BaseEntity {
	return BaseEntity{
		Id: xid.New().String(),
	}
}

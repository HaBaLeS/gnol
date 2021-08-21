package storage

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var META_BUCKET = []byte("meta")

type Metadata struct {
	*BaseEntity
	LastUpdate       time.Time
	FilePath         string
	Name             string
	Type             string
	CoverImageBase64 string
	NumPages         int
	UploadUser       string
	Public           bool
	//arc              archiver.Walker
	owners			 []string
}

func NewMetadata(path string) (error, *Metadata) {
	m := &Metadata{
		FilePath:   path,
		UploadUser: "Anon",
		Public:     false,
		BaseEntity: CreateBaseEntity(META_BUCKET),
	}

	fi, err := os.Stat(path)
	if err != nil {
		return err, nil
	}
	ids := []byte(fmt.Sprintf("%s:%d", fi.Name(), fi.Size()))
	m.ChangeId(fmt.Sprintf("CMX-%x", sha1.Sum(ids)))

	return nil, m
}

func (m *Metadata) UpdateMeta() error {
	fmt.Printf("[i] Create Metadata for: %s\n", m.FilePath)
	m.LastUpdate = time.Now()
	fi, err := os.Stat(m.FilePath)
	if err != nil {
		return err
	}
	m.Name = findName(fi.Name())
	m.Type = filepath.Ext(fi.Name())
	e2 := m.extractCoverImage()
	if e2 != nil {
		return e2
	}
	return nil
}

func (m *Metadata) arc() archiver.Walker {
	ext := filepath.Ext(m.FilePath)
	if ext == ".cbz" || ext == ".zip" {
		return archiver.NewZip()
	} else if ext == ".cbr" || ext == ".rar" {
		return archiver.NewRar()
	}
	return NilWalker{}
}

type NilWalker struct{}
func (r NilWalker) Walk(archive string, walkFn archiver.WalkFunc) error{
	return fmt.Errorf("Unsuported format")
}

func (m *Metadata) extractCoverImage() error {
	var names []string
	aerr := m.arc().Walk(m.FilePath, func(f archiver.File) error {
		if isImageFile(f.Name()) {
			names = append(names, strings.ToLower(f.Name()))
		}
		return nil
	})
	if aerr != nil {
		return aerr
	}
	sort.Strings(names)
	m.NumPages = len(names) //FIXME  -- NumPages is wrong and has issues on the page count when viewing
	for _, name := range names {
		if strings.HasPrefix(name, "."){
			//ignore files starting with a .
		} else if strings.HasSuffix(name, "/") {
			fmt.Printf("Dir in archive: %s\n", name)
		} else if isCoverImage(name) {
			aerr := m.arc().Walk(m.FilePath, func(f archiver.File) error {
				if strings.HasPrefix(f.Name(), "."){
					//ignore files starting with a .
				} else if strings.HasSuffix(strings.ToLower(f.Name()), name) {
					input, err := ioutil.ReadAll(f.ReadCloser)
					if err != nil {
						panic(err)
					}
					res, resErr := util.CreateThumbnail(input, filepath.Ext(name))
					//m.CoverImage = res
					m.CoverImageBase64 = base64.StdEncoding.EncodeToString(res)
					if resErr != nil {
						fmt.Printf("File: %s\n", name)
						panic(resErr)
					}
					//ioutil.WriteFile(TN_FOLDER+"_tn.jpg", res, os.ModePerm)
					fmt.Printf("Choosing image: %s as cover\n", f.Name())
				}
				return nil
			})
			if aerr != nil {
				return aerr
			}
			break
		}
	}
	return nil
}



func isCoverImage(f string) bool {
	if !isImageFile(f) {
		return false
	}
	if strings.Contains(strings.ToLower(f), "banner") {
		return false
	}
	return true
}

func isImageFile(f string) bool {
	ext := strings.ToLower(filepath.Ext(f))
	if strings.HasPrefix(f, ".") {
		return false //no . (hidden files)
	}
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func findName(filename string) string {
	out := filename
	//Remove things between ( )
	re := regexp.MustCompile("\\(.*\\)")

	for _, r := range re.FindAll([]byte(filename), -1) {
		out = strings.ReplaceAll(out, string(r), "")
	}

	//clean separators
	out = strings.ReplaceAll(out, "_", " ")
	out = strings.ReplaceAll(out, "-", " ")
	out = strings.ReplaceAll(out, filepath.Ext(filename), "")
	return out
}

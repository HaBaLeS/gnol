package dao

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
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

type Metadata struct {
	Id               string
	LastUpdate       time.Time
	FilePath         string
	Name             string
	metaFile         string
	Type             string
	CoverImageBase64 string
	NumPages         int
	arc              archiver.Walker
}

func NewMetadata(path string) (error, *Metadata) {
	m := &Metadata{
		FilePath: path,
		metaFile: path + ".meta",
	}

	ext := filepath.Ext(path)
	if ext == ".cbz" || ext == ".zip" {
		m.arc = archiver.NewZip()
	} else if ext == ".cbr" || ext == ".rar" {
		m.arc = archiver.NewRar()
	}
	if m.arc == nil {
		return fmt.Errorf("Unsupported File: %s", path), nil
	}

	fi, err := os.Stat(path)
	if err != nil {
		return err, nil
	}
	ids := []byte(fmt.Sprintf("%s:%d", fi.Name(), fi.Size()))
	m.Id = fmt.Sprintf("CMX-%x", sha1.Sum(ids))

	return nil, m
}

func (m *Metadata) Save() error {
	mf, err := os.Create(m.metaFile)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(mf)
	enc.SetIndent("", "\t")
	enc.Encode(m)
	mf.Close()
	return nil
}

func (m *Metadata) Load() error {
	mf, err := os.Open(m.metaFile)
	if err != nil {
		return err
	}
	jd := json.NewDecoder(mf)
	err = jd.Decode(m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Metadata) Update() error {
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

func (m *Metadata) extractCoverImage() error {
	var names []string
	aerr := m.arc.Walk(m.FilePath, func(f archiver.File) error {
		isImageFile(f.Name())
		names = append(names, strings.ToLower(f.Name()))
		return nil
	})
	if aerr != nil {
		return aerr
	}
	sort.Strings(names)
	m.NumPages = len(names)
	for _, name := range names {
		if strings.HasSuffix(name, "/") {
			fmt.Printf("Dir in archive: %s\n", name)
		} else if isCoverImage(name) {
			aerr := m.arc.Walk(m.FilePath, func(f archiver.File) error {
				if strings.HasSuffix(strings.ToLower(f.Name()), name) {
					input, err := ioutil.ReadAll(f.ReadCloser)
					if err != nil {
						panic(err)
					}
					res, resErr := util.CreateThumbnail(input, filepath.Ext(name))
					//m.CoverImage = res
					m.CoverImageBase64 = base64.StdEncoding.EncodeToString(res)
					if resErr != nil {
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

package persistence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/cmd/leech-tool/modules"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Session struct {
	FileName      string
	mainbucket    []byte
	Start         string
	Count         int
	Workdir       string
	Name          string
	Series        string
	NumInSeries   int
	CoverPage     int
	Nsfw          bool
	Author        string
	Plm           modules.PageLeechModule
	NextSelector  string
	ImageSelector string
	StopOnURl     string
}

type LeechJob struct {
	PageNum         int
	Created         time.Time
	DataUrl         *url.URL
	DataContentType string
	LastScan        time.Time
	FirstScan       time.Time
	ImageUrl        string //fixme convert to URL
	ImageLocalPath  string
	session         *Session
	PageData        []byte
}

/*
func (lj *LeechJob) ID() []byte {
	return []byte(fmt.Sprintf("lu_%s", lj.DataUrl.String()))
}
func (lj *LeechJob) PageID() []byte {
	return []byte(fmt.Sprintf("page_%s", lj.DataUrl.String()))
}*/

/*func (lj *LeechJob) WritePageData(data []byte) error {
	return lj.gnolsession.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(lj.gnolsession.mainbucket)
		b.Put(lj.PageID(), data)
		return nil
	})
} */

func (lj *LeechJob) WriteImageData(data []byte) {
	os.MkdirAll(path.Join("work", lj.session.Workdir), os.ModePerm)
	sfx := getSuffixFromContentType(lj.DataContentType)
	fn := path.Join("work", lj.session.Workdir, fmt.Sprintf("%04d.%s", lj.PageNum, sfx))
	err := ioutil.WriteFile(fn, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
	lj.ImageLocalPath = fn
}

/*
func (lj *LeechJob) Save() {
	err := lj.gnolsession.db.Update(func(tx *bolt.Tx) error {
		tx.Bucket(lj.gnolsession.mainbucket).Put(lj.ID(), marshallJson(lj))
		return nil
	})
	if err != nil {
		panic(err)
	}
}*/

func getSuffixFromContentType(contentType string) string {
	fmt.Printf("ContentType: %s\n", contentType)
	mime := strings.Split(contentType, ";")[0]
	switch mime {
	case "text/html":
		return "html"
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return "unknown"
	}
}

/*
func (gnolsession *Session) OpenDataBase(filepath string) error {
	db, err := bolt.Open(filepath, 0666, nil)
	if err != nil {
		return err
	}

	//Enhance Session with Database Tngs
	gnolsession.db = db
	gnolsession.FileName = filepath
	gnolsession.mainbucket = []byte(gnolsession.Workdir)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(gnolsession.mainbucket)
		return err
	})
	return err
}*/

func (session *Session) LeechJobForURL(leechurl string) *LeechJob {
	url, px := url.Parse(leechurl)
	if px != nil {
		panic(px)
	}
	retVal := &LeechJob{
		Created: time.Now(),
		DataUrl: url,
	}
	//id := retVal.ID()
	/*err := gnolsession.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(gnolsession.mainbucket)
		data := b.Get(id)
		if data != nil {
			return unmarshalJson(data, retVal)
		} else {
			b.Put(id, marshallJson(retVal))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}*/
	retVal.session = session
	return retVal
}

/*
func (gnolsession *Session) HTMLDataForJob(job *LeechJob) []byte {
	var data []byte
	err := gnolsession.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket(gnolsession.mainbucket).Get(job.PageID())
		return nil
	})
	if err != nil {
		panic(err)
	}
	return data
}*/

/*func (gnolsession *Session) CloseDB() {
	if err := gnolsession.db.Close(); err != nil {
		panic(err)
	}
}*/

func (session *Session) WriteMetaFile() {
	f, err := os.Create(path.Join("work", session.Workdir, "gnol.json"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err = enc.Encode(session); err != nil {
		panic(err)
	}

}

func marshallJson(val *LeechJob) []byte {
	w := bytes.NewBuffer(make([]byte, 0))
	enc := json.NewEncoder(w)
	err := enc.Encode(val)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

func unmarshalJson(data []byte, obj interface{}) error {
	dec := json.NewDecoder(bytes.NewBuffer(data))
	return dec.Decode(obj)
}

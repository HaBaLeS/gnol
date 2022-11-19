package persistence

import (
	"encoding/json"
	"github.com/HaBaLeS/gnol/cmd/leech-tool/modules"
	"io/ioutil"
	"net/url"
	"os"
	"path"
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

type LastScrape struct {
	LastScrapeTime time.Time
}

func (lj *LeechJob) WriteImageData(data []byte) {
	os.MkdirAll(path.Join("leech-data", lj.session.Workdir), os.ModePerm)
	err := ioutil.WriteFile(lj.ImageLocalPath, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (session *Session) LeechJobForURL(leechurl string) *LeechJob {
	url, px := url.Parse(leechurl)
	if px != nil {
		panic(px)
	}
	retVal := &LeechJob{
		Created: time.Now(),
		DataUrl: url,
	}
	retVal.session = session
	return retVal
}

func (session *Session) WriteMetaFile() {
	f, err := os.Create(path.Join("leech-data", session.Workdir, "gnol.json"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err = enc.Encode(session); err != nil {
		panic(err)
	}
}

func (session *Session) WriteScrapeStatusFile() {
	f, err := os.Create(path.Join("leech-data", session.Workdir, "last-scrape.json"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)

	ls := &LastScrape{
		LastScrapeTime: time.Now(),
	}
	if err = enc.Encode(ls); err != nil {
		panic(err)
	}
}

func (session *Session) LoadScrapeStatusFile() *LastScrape {
	f, err := os.Open(path.Join("leech-data", session.Workdir, "last-scrape.json"))
	if err != nil {
		return nil
	}
	defer f.Close()
	enc := json.NewDecoder(f)

	ls := &LastScrape{}
	if err = enc.Decode(ls); err != nil {
		panic(err)
	}
	return ls
}

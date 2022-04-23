package engine

import (
	"bytes"
	"fmt"
	"github.com/HaBaLeS/gnol/cmd/leech-tool/network"
	"github.com/HaBaLeS/gnol/cmd/leech-tool/persistence"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"os"
	"strings"
)

type Engine struct {
	Session *persistence.Session
}

func (e *Engine) Leech() {
	leechJob := e.Session.LeechJobForURL(e.Session.Start)
	leechJob.PageNum = 0
	next := e.Session.Start
	for next != "" {
		fmt.Printf("Processing Page #%04d %s", e.Session.Count, next)

		n, err := url.Parse(next)
		panicIfErr(err)
		lj := e.Session.LeechJobForURL(n.String())
		lj.PageNum = e.Session.Count
		next = e.processComicPage(lj)
		if strings.HasPrefix(next, "/") {
			next = lj.DataUrl.Scheme + "://" + lj.DataUrl.Host + next
		}
		u, _ := url.Parse(next)
		next = lj.DataUrl.ResolveReference(u).String()
		e.Session.Count++
		if e.Session.StopOnURl != "" && next == e.Session.StopOnURl {
			break
		}
		fmt.Print("\n")
	}
}

func (e *Engine) processComicPage(job *persistence.LeechJob) string {

	data := job.PageData
	if data == nil {
		data = network.DownloadForJob(job)
	} else {
		fmt.Print(" [C] ")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	image := e.Session.Plm.FindCurrentImage(doc)
	u, _ := url.Parse(image)
	job.ImageUrl = job.DataUrl.ResolveReference(u).String()
	fmt.Printf("Found: %s", job.ImageUrl)
	e.processComicImage(job)

	return e.Session.Plm.FindNextPage(doc)
}

func (e *Engine) processComicImage(job *persistence.LeechJob) {
	dl := true
	if job.ImageUrl != "" {
		fi, err := os.Stat(job.ImageLocalPath)
		if err != nil || fi.Size() < 10*1024 {
			dl = true
		} else {
			dl = false
			fmt.Print(" [X] ")
		}
	}
	if dl {
		data, ct := network.DownloadForUrl(job.ImageUrl)
		job.DataContentType = ct
		job.WriteImageData(data)
	}

	//job.Save()
}

func panicIfErr(err interface{}) {
	if err != nil {
		panic(err)
	}
}

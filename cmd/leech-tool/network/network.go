package network

import (
	"fmt"
	"github.com/HaBaLeS/gnol/cmd/leech-tool/persistence"
	"io/ioutil"
	"net/http"
	"time"
)

func DownloadForJob(job *persistence.LeechJob) []byte {
	res, err := http.Get(job.DataUrl.String())
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic(fmt.Errorf("Got Status %d wanted 200", res.StatusCode))
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	job.LastScan = time.Now()
	job.DataContentType = res.Header.Get("content-type")
	job.PageData = data
	fmt.Print(" [D] ")
	return data
}

func DownloadForUrl(url string) ([]byte, string) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic(fmt.Errorf("Got Status %d wanted 200", res.StatusCode))
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Print(" [D] ")
	return data, res.Header.Get("content-type")
}

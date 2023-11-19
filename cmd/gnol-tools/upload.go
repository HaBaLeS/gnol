package main

import (
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/server/dto"
	"github.com/HaBaLeS/gnol/server/router"
	"github.com/HaBaLeS/gnol/server/util"
	"io"
	"net/http"
	"os"
)

func (s *Session) upload(args []string, options map[string]string) int {
	if !s.processOptionsAndValidate(args, options) {
		return -1
	}
	return s.uploadInternal()
}

type UplProgessReader struct {
	io.Reader
	total, sum int64
}

func (rdr *UplProgessReader) Read(p []byte) (n int, err error) {
	n, err = rdr.Reader.Read(p)
	rdr.sum += int64(n)
	fmt.Printf("Status:\033[0K %d/%d\r", rdr.sum, rdr.total)
	return n, err
}

func (s *Session) uploadInternal() int {

	if exist, obj := s.checkIfFileExists(); exist {
		s.Logger.Printf("File exists on Server not uploading it!! %v", obj)
		return 0
	}

	inFile, err := os.Open(s.InputFile)

	pr := &UplProgessReader{
		Reader: inFile,
	}

	if fi, err := inFile.Stat(); err != nil {
		panic(err)
	} else {
		pr.total = fi.Size()
	}

	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	rq, err := http.NewRequest("POST", s.GnolHost, pr)
	if err != nil {
		panic(err)
	}
	rq.Header.Add(router.API_GNOL_TOKEN, s.ApiToken)
	q := rq.URL.Query()
	q.Add(router.API_SERIES_ID, s.SeriesId)
	if s.MetaData.Nsfw {
		q.Add(router.API_NSFW, "srly")
	}
	q.Add(router.API_ODER_NUM, s.OrderNum)
	rq.URL.RawQuery = q.Encode()

	fmt.Printf("Sending: %s\n\n", rq.URL.RequestURI())
	resp, err := client.Do(rq)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Status %s\n", resp.Status)
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Data:\n%s", data)
		return -1
	}

	fmt.Printf("Resp: %s", resp.Body)

	return 0
}

func (s *Session) checkIfFileExists() (bool, *dto.ComicEntry) {
	hash, err := util.HashFile(s.InputFile)
	if err != nil {
		panic(err) //file must exist at that point
	}

	url := fmt.Sprintf("%s/%s/%s", s.GnolHost, "checkhash", hash)
	s.Logger.Printf("SQuery API: %s", url)
	client := http.DefaultClient
	rq, err := http.NewRequest("GET", url, nil)
	rq.Header.Add(router.API_GNOL_TOKEN, s.ApiToken)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(rq)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	dto := &dto.ComicEntry{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(dto); err != nil {
		panic(err)
	}
	return true, dto
}

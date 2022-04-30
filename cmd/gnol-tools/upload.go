package main

import (
	"fmt"
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

func (s *Session) uploadInternal() int {
	host := "gnol.habales.de"
	port := 443
	protocol := "https"
	path := "api/upload"
	secret := "8baf2620-a419-4e97-bd3c-6de387a0d897"

	url := fmt.Sprintf("%s://%s:%d/%s", protocol, host, port, path)
	fmt.Printf("Posting to: %s\n", url)
	uplf, err := os.Open(s.InputFile)
	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	rq, err := http.NewRequest("POST", url, uplf)
	if err != nil {
		panic(err)
	}
	rq.Header.Add("gnol-token", secret)

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

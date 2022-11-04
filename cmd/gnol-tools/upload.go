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
	uplf, err := os.Open(s.InputFile)
	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	rq, err := http.NewRequest("POST", s.GnolHost, uplf)
	if err != nil {
		panic(err)
	}
	rq.Header.Add("gnol-token", s.ApiToken)

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

package main

import (
	"fmt"
	"github.com/HaBaLeS/gnol/cmd/h-leech-tool/engine"
)

func (s *Session) leechlist(args []string, options map[string]string) int {
	fmt.Printf("h-leech-tool \n")
	//https://ilikecomix.com/
	return 0
}

func (s *Session) leechpull(args []string, options map[string]string) int {
	target := args[0]
	s.Verbose = true
	fmt.Printf("leechpull %s\n", target)

	path := engine.Leech(target)
	s.InputFile = path
	s.OutputFile = "tmp.cbz"

	//FIXME set session value correct
	s.packInternal()

	//FIXME Upload
	s.InputFile = s.OutputFile
	if !s.processOptionsAndValidate(args, options) {
		return -1
	}

	s.SeriesId = "49"
	s.MetaData.Nsfw = true
	s.uploadInternal()

	//Fixme cleanup

	return 0

}

func (s *Session) leechupdate(args []string, options map[string]string) int {
	fmt.Printf("leechupdate \n")
	return 0

}

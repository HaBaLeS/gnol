package main

import (
	"flag"
	"playground.dahoam/server"
)

var BASE_PATH = "/home/falko/comics/"

//go:generate go run -tags=dev gen.go

func main() {

	cfgPath := flag.String("c", "default.cfg", "Config File to use")
	flag.Parse()

	s := server.NewServer(*cfgPath)
	s.Start()
}

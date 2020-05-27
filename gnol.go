package main

import (
	"flag"
	"github.com/HaBaLeS/gnol/server"
)

//go:generate go run -tags=dev gen.go

func main() {
	cfgPath := flag.String("c", "default.cfg", "Config File to use")
	flag.Parse()

	s := server.NewServer(*cfgPath)
	s.Start()
}

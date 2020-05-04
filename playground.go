package main

import (
	"playground.dahoam/server"
)

var BASE_PATH = "/home/falko/comics/"


func main(){

	server.GetComiList()

	s := server.NewServer()
	s.Start()
}


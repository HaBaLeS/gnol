package main

import (
	"fmt"
	"playground.dahoam/server"
)

var BASE_PATH = "/home/falko/comics/"


func main(){

	fmt.Print("http://192.168.1.248:6969/comics\n")
	s := server.NewServer()
	s.Start()
}


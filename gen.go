// +build ignore

package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
)

func main() {
	// Assets contains project assets.
	var StaticAssets http.FileSystem = http.Dir("data/")
	var err = vfsgen.Generate(StaticAssets, vfsgen.Options{
		PackageName:  "util",
		Filename:     "server/util/static_assets.go",
		BuildTags:    "!dev",
		VariableName: "StaticAssets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}

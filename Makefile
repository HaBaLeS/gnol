

all: build

build: generate
	go build

generate: get
	go generate

get:
	go get -u github.com/shurcooL/vfsgen

clean:
	go clean


all: build

release: lint build
	echo "RElease loool"

build: generate
	go build

generate: get
	go generate

get:
	go get -u github.com/shurcooL/vfsgen

lint:
	golint  ./...

clean:
	go clean
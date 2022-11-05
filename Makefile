

all: build build-tools build-leech check test

release: check test build

build: prepare generate
	mkdir -p bin
	go build -o bin/gnol
	go build -o bin/gnol-tools ./cmd/gnol-tools
	go build -o bin/leech-tool ./cmd/leech-tool

install:
	echo "installing GO Tools"
	go install ./cmd/gnol-tools
	go install ./cmd/leech-tool


generate:
	go generate

check:
	echo "Skipping checks fix makefile please"
#	golint  ./...
#	go vet ./...

prepare:
	go get github.com/shurcooL/vfsgen@92b8a710ab6cab4c09182a1fcf469157bc938f8f
	go get golang.org/x/tools/...
	go get golang.org/x/lint/golint

test:
	echo "Skipping test fix them please"
	#go test ./...  -cover -coverprofile=c.out
	#go tool cover -html=c.out -o coverage.html


container:
	docker build . -t reg.habales.de/gnol/gnol:0.7.1

push:
	docker push reg.habales.de/gnol/gnol:0.7.1

clean:
	go clean
	rm c.out coverage.html
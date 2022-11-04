

all: build check test

release: check test build

build: prepare generate
	go build

generate:
	go generate

check:
	golint  ./...
	go vet ./...

prepare:
	go get github.com/shurcooL/vfsgen@92b8a710ab6cab4c09182a1fcf469157bc938f8f
	go get golang.org/x/tools/...
	go get golang.org/x/lint/golint

test:
	go test ./...  -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html


container:
	docker build . -t reg.habales.de/gnol/gnol:0.7.1

push:
	docker push reg.habales.de/gnol/gnol:0.7.1

clean:
	go clean
	rm c.out coverage.html
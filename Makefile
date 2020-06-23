

all: build

release: check test build
	echo "RElease loool"

build: generate
	go build

generate: get
	go generate

check: get
	golint  ./...
	go vet ./...

get:
	go get -u github.com/shurcooL/vfsgen
	go get -u golang.org/x/tools/...
	go get -u golang.org/x/lint/golint

test: get
	go test ./...  -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

clean:
	go clean
	rm c.out coverage.html
NOW=$(shell date +'%Y-%m-%d_%T')
VERSION=0.8.1
LD_FLAG=-X main.VersionNum=$(VERSION) -X main.BuildDate=$(NOW)


all: build check test

release: check test build

build: generate
	mkdir -p bin
	go build  -ldflags '$(LD_FLAG)' -v  -o bin/gnol
	go build  -ldflags '$(LD_FLAG)' -v  -o bin/gnol-tools ./cmd/gnol-tools
	go build  -ldflags '$(LD_FLAG)' -v  -o bin/leech-tool ./cmd/leech-tool

install:
	go install -ldflags '$(LD_FLAG)' -v  ./cmd/gnol-tools
	go install -ldflags '$(LD_FLAG)' -v  ./cmd/leech-tool


generate:
	go generate

check:
	echo "Skipping checks fix makefile please"
#	golint  ./...
#	go vet ./...

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
	rm -rf bin
NOW=$(shell date +'%Y-%m-%d_%T')
VERSION=0.17.4
LD_FLAG=-X main.VersionNum=$(VERSION) -X main.BuildDate=$(NOW)


all: build check test

release: check test build

build-server:
	mkdir -p bin
	go build  -ldflags '$(LD_FLAG)' -v  -o bin/gnol

build:
	mkdir -p bin
	go build  -ldflags '$(LD_FLAG)' -v  -o bin/gnol
	go build  -ldflags '$(LD_FLAG)' -v  -o bin/gnol-tools ./cmd/gnol-tools

install:
	go install -ldflags '$(LD_FLAG)' -v  ./cmd/gnol-tools

check:
	echo "Skipping checks fix makefile please"
#	golint  ./...
#	go vet ./...

test:
	echo "Skipping test fix them please"
	#go test ./...  -cover -coverprofile=c.out
	#go tool cover -html=c.out -o coverage.html


container:
	docker build . -t reg.habales.de/gnol/gnol:$(VERSION)
	docker tag reg.habales.de/gnol/gnol:$(VERSION) reg.habales.de/gnol/gnol:latest

push: container
	docker push reg.habales.de/gnol/gnol:$(VERSION)
	docker push reg.habales.de/gnol/gnol:latest

clean:
	go clean
	rm -rf bin

GIT_COMMIT=$(shell git rev-parse HEAD)
KVEXPRESS_VERSION=$(shell ./version)
COMPILE_DATE=$(shell date -u +%Y%m%d.%H%M%S)
BUILD_FLAGS=-X main.CompileDate=$(COMPILE_DATE) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(KVEXPRESS_VERSION)

all: build

deps:
	go get github.com/aryann/difflib
	go get github.com/spf13/cobra
	go get github.com/hashicorp/consul/api
	go get github.com/zorkian/go-datadog-api
	go get github.com/PagerDuty/godspeed

format:
	gofmt -w .

clean:
	rm -f bin/kvexpress || true

build: clean
	go build -ldflags "$(BUILD_FLAGS)" -o bin/kvexpress main.go

gziposx:
	gzip bin/kvexpress
	mv bin/kvexpress.gz bin/kvexpress-$(KVEXPRESS_VERSION)-darwin.gz

linux: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "$(BUILD_FLAGS)" -o bin/kvexpress main.go

gziplinux:
	gzip bin/kvexpress
	mv bin/kvexpress.gz bin/kvexpress-$(KVEXPRESS_VERSION)-linux-amd64.gz

release: clean build gziposx clean linux gziplinux clean

consul:
	consul agent -data-dir `mktemp -d`  -bootstrap -server -bind=127.0.0.1

sorting:
	curl -s https://gist.githubusercontent.com/darron/94447bfab90617f16962/raw/d4cb39471724800ba9e731f99e5844167e93c5df/sorting.txt > sorting

KVEXPRESS_VERSION="1.10"
GIT_COMMIT=$(shell git rev-parse HEAD)
COMPILE_DATE=$(shell date -u +%Y%m%d.%H%M%S)
BUILD_FLAGS=-X main.CompileDate=$(COMPILE_DATE) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(KVEXPRESS_VERSION)
UNAME=$(shell uname -s | tr '[:upper:]' '[:lower:]')

all: build

deps:
	go get -u github.com/progrium/basht
	go get -u github.com/CiscoCloud/consul-cli
	go get -u github.com/DataDog/kvexpress

format:
	gofmt -w .

clean:
	rm -f bin/kvexpress || true

build: clean
	go build -ldflags "$(BUILD_FLAGS)" -o bin/kvexpress main.go

gzip:
	gzip bin/kvexpress
	mv bin/kvexpress.gz bin/kvexpress-$(KVEXPRESS_VERSION)-$(UNAME).gz

release: clean build gzip

consul:
	consul agent -data-dir `mktemp -d` -bootstrap -server -bind=127.0.0.1 1>/dev/null &

consul_kill:
	pkill consul

sorting:
	curl -s https://gist.githubusercontent.com/darron/94447bfab90617f16962/raw/d4cb39471724800ba9e731f99e5844167e93c5df/sorting.txt > sorting

unit:
	cd commands && go test -v -cover

test: unit wercker

wercker_clean:
	bin/kvexpress clean -f sorting
	rm -f output ignored lock-test lock-test.locked url raw_checksum url_exec additional-file decompressed

wercker: consul sorting unit
	basht test/tests.bash

GIT_COMMIT=$(shell git rev-parse HEAD)
KVEXPRESS_VERSION=$(shell ./version)
COMPILE_DATE=$(shell date -u +%Y%m%d.%H%M%S)
BUILD_FLAGS=-X main.CompileDate=$(COMPILE_DATE) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(KVEXPRESS_VERSION)

all: build

deps:
	go get -u github.com/spf13/cobra
	go get -u github.com/hashicorp/consul/api
	go get -u github.com/zorkian/go-datadog-api
	go get -u github.com/PagerDuty/godspeed
	go get -u gopkg.in/yaml.v2
	go get -u github.com/smallfish/simpleyaml
	go get -u github.com/progrium/basht
	go get -u github.com/CiscoCloud/consul-cli

format:
	gofmt -w .

clean:
	rm -f bin/kvexpress || true

build: clean
	go build -ldflags "$(BUILD_FLAGS)" -o bin/kvexpress main.go

gziposx:
	gzip bin/kvexpress
	mv bin/kvexpress.gz bin/kvexpress-$(KVEXPRESS_VERSION)-darwin.gz

release: clean build gziposx

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

wercker: consul sorting
	basht test/tests.bash

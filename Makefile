all: deps build

deps:
	go get github.com/aryann/difflib
	go get github.com/spf13/cobra
	go get github.com/hashicorp/consul/api
	go get github.com/zorkian/go-datadog-api

build:
	rm -f kvexpress || true
	go build -o kvexpress main.go

all: build

deps:
	go get github.com/aryann/difflib
	go get github.com/spf13/cobra
	go get github.com/hashicorp/consul/api
	go get github.com/zorkian/go-datadog-api

format:
	gofmt -w .

clean:
	rm -f kvexpress || true

build: clean
	go build -o kvexpress main.go

linux: clean
	GOOS=linux GOARCH=amd64 go build -o kvexpress main.go

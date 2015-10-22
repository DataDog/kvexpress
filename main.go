package main

import (
	"./commands/"
	"runtime"
	// "crypto/sha256"
	// "github.com/aryann/difflib"
	// consul "github.com/hashicorp/consul/api"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.RootCmd.Execute()
}

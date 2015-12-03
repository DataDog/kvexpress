package main

import (
	"./commands/"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"runtime"
)

var CompileDate = "No date provided."
var GitCommit = "No revision provided."
var Version = "No version provided."

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "kvexpress")
	if e == nil {
		log.SetOutput(logwriter)
	}
	commands.Log(fmt.Sprintf("main: kvexpress version:%s git:%s date:%s", Version, GitCommit, CompileDate), "info")

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("Version: %s\nRevision: %s\nDate: %s\n", Version, GitCommit, CompileDate)
			os.Exit(0)
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.RootCmd.Execute()
}

package main

import (
	"./commands/"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"runtime"
)

var minversion = "No version provided."
var GitCommit = "No revision provided."

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "kvexpress")
	if e == nil {
		log.SetOutput(logwriter)
	}
	log.Print("main: Startup kvexpress version:", minversion, " git:", GitCommit)

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("Version: %s\nRevision: %s\n", minversion, GitCommit)
			os.Exit(0)
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.RootCmd.Execute()
}

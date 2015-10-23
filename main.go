package main

import (
	"./commands/"
	"log"
	"log/syslog"
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.RootCmd.Execute()
}

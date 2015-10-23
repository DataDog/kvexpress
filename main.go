package main

import (
	"./commands/"
	"log"
	"log/syslog"
	"runtime"
)

var minversion = "No Version Provided."

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "kvexpress")
	if e == nil {
		log.SetOutput(logwriter)
	}
	log.Print("Startup kvexpress version:", minversion)
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.RootCmd.Execute()
}

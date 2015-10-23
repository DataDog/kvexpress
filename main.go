package main

import (
	"./commands/"
	"log"
	"log/syslog"
	"runtime"
	// "crypto/sha256"
)

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "kvexpress")
	if e == nil {
		log.SetOutput(logwriter)
	}
	log.Print("Startup")
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.RootCmd.Execute()
}

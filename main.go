// +build linux darwin freebsd

package main

import (
	"fmt"
	"github.com/DataDog/kvexpress/commands"
	"log"
	"log/syslog"
	"os"
	"runtime"
)

// CompileDate tracks when the binary was compiled. It's inserted during a build
// with build flags. Take a look at the Makefile for information.
var CompileDate = "No date provided."

// GitCommit tracks the SHA of the built binary. It's inserted during a build
// with build flags. Take a look at the Makefile for information.
var GitCommit = "No revision provided."

// Version is the version of the built binary. It's inserted during a build
// with build flags. Take a look at the Makefile for information.
var Version = "No version provided."

// GoVersion details the version of Go this was compiled with.
var GoVersion = runtime.Version()

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "kvexpress")
	if e == nil {
		log.SetFlags(log.Lmicroseconds)
		log.SetOutput(logwriter)
	}
	commands.Log(fmt.Sprintf("kvexpress version:%s", Version), "info")

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("%-8s : %s\n%-8s : %s\n%-8s : %s\n%-8s : %s\n", "Version", Version, "Revision", GitCommit, "Date", CompileDate, "Go", GoVersion)
			os.Exit(0)
		}
	}

	commands.RootCmd.Execute()
}

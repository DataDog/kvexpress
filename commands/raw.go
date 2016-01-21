// +build linux darwin freebsd

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "Write a file pulled from any Consul KV data.",
	Long:  `Raw is for writing a file based on any Consul key.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkRawFlags()
		AutoEnable()
	},
	Run: rawRun,
}

func rawRun(cmd *cobra.Command, args []string) {
	start := time.Now()

	c, _ := Connect(ConsulServer, Token)

	// Get the KV data out of Consul.
	KVData := Get(c, RawKeyOutLocation)

	// Is the data long enough?
	longEnough := LengthCheck(KVData, MinFileLength)
	Log(fmt.Sprintf("longEnough='%s'", strconv.FormatBool(longEnough)), "debug")

	// If the data is long enough, write the file.
	if longEnough {
		// Acually write the file.
		WriteFile(KVData, RawFiletoWrite, FilePermissions, Owner)
		StatsdRaw(RawKeyOutLocation)
	} else {
		Log("longEnough='no'", "info")
		os.Exit(0)
	}

	// Run this command after the file is written.
	if PostExec != "" {
		Log(fmt.Sprintf("exec='%s'", PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, RawKeyOutLocation, "complete")
}

func checkRawFlags() {
	Log("Checking cli flags.", "debug")
	if RawKeyOutLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if RawFiletoWrite == "" {
		fmt.Println("Need a file to write in -f")
		os.Exit(1)
	}
	Log("Required cli flags present.", "debug")
}

var (
	// RawKeyOutLocation This Consul key is the location we want to pull data from.
	// This data can be ANY Consul key and doesn't have to be in any particular format or structure.
	// It doesn't have to have a Checksum either - it can be any Consul key at all.
	// Give the complete path - does not use PrefixLocation.
	RawKeyOutLocation string

	// RawFiletoWrite is the location we want to write the data to.
	RawFiletoWrite string
)

func init() {
	RootCmd.AddCommand(rawCmd)
	rawCmd.Flags().StringVarP(&RawKeyOutLocation, "key", "k", "", "Raw key to pull data from")
	rawCmd.Flags().StringVarP(&RawFiletoWrite, "file", "f", "", "where to write the data")
}

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
	Long:  `raw is for writing a file based on any Consul key.`,
	Run:   rawRun,
}

func rawRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	checkRawFlags(Direction)

	c, _ := Connect(ConsulServer, Token, Direction)

	// Get the KV data out of Consul.
	KVData := GetRaw(c, PrefixLocation, RawKeyOutLocation, Direction, DogStatsd)

	// Is the data long enough?
	longEnough := LengthCheck(KVData, MinFileLength, Direction)
	Log(fmt.Sprintf("%s: longEnough='%s'", Direction, strconv.FormatBool(longEnough)), "debug")

	// If the data is long enough, write the file.
	if longEnough {
		// Acually write the file.
		WriteFile(KVData, RawFiletoWrite, FilePermissions, Owner, Direction, DogStatsd)
		if DogStatsd {
			StatsdRaw(RawKeyOutLocation)
		}
	} else {
		Log(fmt.Sprintf("%s: longEnough='no'", Direction), "info")
		os.Exit(0)
	}

	// Run this command after the file is written.
	if PostExec != "" {
		Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, RawKeyOutLocation, "complete", Direction, DogStatsd)
}

func checkRawFlags(direction string) {
	Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if RawKeyOutLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if RawFiletoWrite == "" {
		fmt.Println("Need a file to write in -f")
		os.Exit(1)
	}
	if DogStatsd {
		Log(fmt.Sprintf("%s: Enabling Dogstatsd metrics.", direction), "debug")
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		Log(fmt.Sprintf("%s: Enabling Datadog API.", direction), "debug")
	}
	if Owner == "" {
		Owner = GetCurrentUsername(direction)
	}
	Log(fmt.Sprintf("%s: Required cli flags present.", direction), "debug")
}

var (
	RawKeyOutLocation string
	RawFiletoWrite    string
)

func init() {
	RootCmd.AddCommand(rawCmd)
	rawCmd.Flags().StringVarP(&RawKeyOutLocation, "key", "k", "", "key to pull data from")
	rawCmd.Flags().StringVarP(&RawFiletoWrite, "file", "f", "", "where to write the data")
}

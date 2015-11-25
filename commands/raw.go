package commands

import (
	kvexpress "../kvexpress/"
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
	var Direction = "raw"
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	checkRawFlags(Direction)

	c, _ := kvexpress.Connect(ConsulServer, Token, Direction)

	// Get the KV data out of Consul.
	KVData := kvexpress.GetRaw(c, PrefixLocation, RawKeyOutLocation, Direction, DogStatsd)

	// Is the data long enough?
	longEnough := kvexpress.LengthCheck(KVData, MinFileLength, Direction)
	kvexpress.Log(fmt.Sprintf("%s: longEnough='%s'", Direction, strconv.FormatBool(longEnough)), "debug")

	// If the data is long enough, write the file.
	if longEnough {
		// Acually write the file.
		kvexpress.WriteFile(KVData, RawFiletoWrite, FilePermissions, Owner, Direction, DogStatsd)
		if DogStatsd {
			kvexpress.StatsdRaw(RawKeyOutLocation)
		}
	} else {
		kvexpress.Log(fmt.Sprintf("%s: longEnough='no'", Direction), "info")
		os.Exit(0)
	}

	// Run this command after the file is written.
	if PostExec != "" {
		kvexpress.Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, RawKeyOutLocation, "complete", Direction, DogStatsd)
}

func checkRawFlags(direction string) {
	kvexpress.Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if RawKeyOutLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if RawFiletoWrite == "" {
		fmt.Println("Need a file to write in -f")
		os.Exit(1)
	}
	if DogStatsd {
		kvexpress.Log(fmt.Sprintf("%s: Enabling Dogstatsd metrics.", direction), "debug")
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		kvexpress.Log(fmt.Sprintf("%s: Enabling Datadog API.", direction), "debug")
	}
	if Owner == "" {
		Owner = kvexpress.GetCurrentUsername(direction)
	}
	kvexpress.Log(fmt.Sprintf("%s: Required cli flags present.", direction), "debug")
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

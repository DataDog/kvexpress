package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var outCmd = &cobra.Command{
	Use:   "out",
	Short: "Write a file based on key data.",
	Long:  `out is for writing a file based on a Consul key.`,
	Run:   outRun,
}

func outRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	var Direction = "out"
	if EnvVars {
		ConfigEnvVars(Direction)
	}
	checkOutFlags(Direction)

	KeyData := kvexpress.KeyDataPath(KeyOutLocation, PrefixLocation, Direction)
	KeyChecksum := kvexpress.KeyChecksumPath(KeyOutLocation, PrefixLocation, Direction)
	KeyStop := kvexpress.KeyStopPath(KeyOutLocation, PrefixLocation, Direction)

	c, _ := kvexpress.Connect(ConsulServer, Token, Direction)

	StopKeyData := kvexpress.Get(c, KeyStop, Direction)

	if StopKeyData != "" && IgnoreStop == false {
		kvexpress.Log(fmt.Sprintf("%s: Stop Key is present - stopping. Reason: %s", Direction, StopKeyData), "info")
		kvexpress.RunTime(start, KeyOutLocation, "stop_key", Direction, DogStatsd)
		os.Exit(0)
	} else {
		if IgnoreStop {
			kvexpress.Log(fmt.Sprintf("%s: Ignoring any stop key.", Direction), "info")
		} else {
			kvexpress.Log(fmt.Sprintf("%s: Stop Key is NOT present - continuing.", Direction), "debug")
		}
	}

	// Get the KV data out of Consul.
	KVData := kvexpress.Get(c, KeyData, Direction)

	// Get the Checksum data out of Consul.
	Checksum := kvexpress.Get(c, KeyChecksum, Direction)

	// Is the data long enough?
	longEnough := kvexpress.LengthCheck(KVData, MinFileLength, Direction)
	kvexpress.Log(fmt.Sprintf("%s: longEnough='%s'", Direction, strconv.FormatBool(longEnough)), "debug")

	// Does the checksum match?
	checksumMatch := kvexpress.ChecksumCompare(KVData, Checksum, Direction)
	kvexpress.Log(fmt.Sprintf("%s: checksumMatch='%s'", Direction, strconv.FormatBool(checksumMatch)), "debug")

	// If the data is long enough and the checksum matches, write the file.
	if longEnough && checksumMatch {
		// Does the file already present in FiletoWrite have the same checksum?
		// Is it directory? Does it exist?
		kvexpress.CheckFiletoWrite(FiletoWrite, Checksum, Direction)

		// Acually write the file.
		kvexpress.WriteFile(KVData, FiletoWrite, FilePermissions, Direction)
		if DogStatsd {
			kvexpress.StatsdOut(KeyOutLocation)
		}
	} else {
		kvexpress.Log(fmt.Sprintf("%s: Could not write file.", Direction), "debug")
	}

	// Run this command after the file is written.
	if PostExec != "" {
		kvexpress.Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, KeyOutLocation, "complete", Direction, DogStatsd)
}

func checkOutFlags(direction string) {
	kvexpress.Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if KeyOutLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if FiletoWrite == "" {
		fmt.Println("Need a file to write in -f")
		os.Exit(1)
	}
	if DogStatsd {
		kvexpress.Log(fmt.Sprintf("%s: Enabling Dogstatsd metrics.", direction), "debug")
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		kvexpress.Log(fmt.Sprintf("%s: Enabling Datadog API.", direction), "debug")
	}
	kvexpress.Log(fmt.Sprintf("%s: Required cli flags present.", direction), "debug")
}

var (
	KeyOutLocation string
	FiletoWrite    string
	IgnoreStop     bool
)

func init() {
	RootCmd.AddCommand(outCmd)
	outCmd.Flags().StringVarP(&KeyOutLocation, "key", "k", "", "key to pull data from")
	outCmd.Flags().StringVarP(&FiletoWrite, "file", "f", "", "where to write the data")
	outCmd.Flags().BoolVarP(&IgnoreStop, "ignore_stop", "", false, "ignore stop key")
}

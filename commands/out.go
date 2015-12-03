package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var outCmd = &cobra.Command{
	Use:   "out",
	Short: "Write a file based on kvexpress organized data stored in Consul.",
	Long:  `out is for writing a file based on a Consul kvexpress key and checksum.`,
	Run:   outRun,
}

func outRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	checkOutFlags(Direction)

	KeyData := KeyDataPath(KeyOutLocation, PrefixLocation, Direction)
	KeyChecksum := KeyChecksumPath(KeyOutLocation, PrefixLocation, Direction)
	KeyStop := KeyStopPath(KeyOutLocation, PrefixLocation, Direction)

	c, _ := Connect(ConsulServer, Token, Direction)

	StopKeyData := Get(c, KeyStop, Direction, DogStatsd)

	if StopKeyData != "" && IgnoreStop == false {
		Log(fmt.Sprintf("%s: Stop Key is present - stopping. Reason: %s", Direction, StopKeyData), "info")
		RunTime(start, KeyOutLocation, "stop_key", Direction, DogStatsd)
		os.Exit(0)
	} else {
		if IgnoreStop {
			Log(fmt.Sprintf("%s: Ignoring any stop key.", Direction), "info")
		} else {
			Log(fmt.Sprintf("%s: Stop Key is NOT present - continuing.", Direction), "debug")
		}
	}

	// Get the KV data out of Consul.
	KVData := Get(c, KeyData, Direction, DogStatsd)

	// Decompress here if necessary.
	if Compress {
		KVData = DecompressData(KVData, Direction)
	}

	// Get the Checksum data out of Consul.
	Checksum := Get(c, KeyChecksum, Direction, DogStatsd)

	// Is the data long enough?
	longEnough := LengthCheck(KVData, MinFileLength, Direction)
	Log(fmt.Sprintf("%s: longEnough='%s'", Direction, strconv.FormatBool(longEnough)), "debug")

	// Does the checksum match?
	checksumMatch := ChecksumCompare(KVData, Checksum, Direction)
	Log(fmt.Sprintf("%s: checksumMatch='%s'", Direction, strconv.FormatBool(checksumMatch)), "debug")

	// If the data is long enough and the checksum matches, write the file.
	if longEnough && checksumMatch {
		// Does the file already present in FiletoWrite have the same checksum?
		// Is it directory? Does it exist?
		CheckFiletoWrite(FiletoWrite, Checksum, Direction)

		// Acually write the file.
		WriteFile(KVData, FiletoWrite, FilePermissions, Owner, Direction, DogStatsd)
		if DogStatsd {
			StatsdOut(KeyOutLocation)
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
	RunTime(start, KeyOutLocation, "complete", Direction, DogStatsd)
}

func checkOutFlags(direction string) {
	Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if KeyOutLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if FiletoWrite == "" {
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

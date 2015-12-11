package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
	checkOutFlags()

	KeyData := KeyDataPath(KeyOutLocation)
	KeyChecksum := KeyChecksumPath(KeyOutLocation)
	KeyStop := KeyStopPath(KeyOutLocation)

	c, _ := Connect(ConsulServer, Token)

	StopKeyData := Get(c, KeyStop, DogStatsd)

	if StopKeyData != "" && IgnoreStop == false {
		Log(fmt.Sprintf("Stop Key is present - stopping. Reason: %s", StopKeyData), "info")
		RunTime(start, KeyOutLocation, "stop_key", DogStatsd)
		os.Exit(0)
	} else {
		if IgnoreStop {
			Log("Ignoring any stop key.", "info")
		} else {
			Log("Stop Key is NOT present - continuing.", "debug")
		}
	}

	// Get the KV data out of Consul.
	KVData := Get(c, KeyData, DogStatsd)

	// Decompress here if necessary.
	if Compress {
		KVData = DecompressData(KVData)
	}

	// Get the Checksum data out of Consul.
	Checksum := Get(c, KeyChecksum, DogStatsd)

	// Is the data long enough?
	longEnough := LengthCheck(KVData, MinFileLength)
	Log(fmt.Sprintf("longEnough='%t'", longEnough), "debug")

	// Does the checksum match?
	checksumMatch := ChecksumCompare(KVData, Checksum)
	Log(fmt.Sprintf("checksumMatch='%t'", checksumMatch), "debug")

	// If the data is long enough and the checksum matches, write the file.
	if longEnough && checksumMatch {
		// Does the file already present in FiletoWrite have the same checksum?
		// Is it directory? Does it exist?
		CheckFiletoWrite(FiletoWrite, Checksum)

		// Acually write the file.
		WriteFile(KVData, FiletoWrite, FilePermissions, Owner, DogStatsd)
		if DogStatsd {
			StatsdOut(KeyOutLocation)
		}
	} else {
		Log("longEnough='no'", "info")
		os.Exit(0)
	}

	// Run this command after the file is written.
	if PostExec != "" {
		Log(fmt.Sprintf("exec='%s'", PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, KeyOutLocation, "complete", DogStatsd)
}

func checkOutFlags() {
	Log("Checking cli flags.", "debug")
	if KeyOutLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if FiletoWrite == "" {
		fmt.Println("Need a file to write in -f")
		os.Exit(1)
	}
	if DogStatsd {
		Log("Enabling Dogstatsd metrics.", "debug")
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		Log("Enabling Datadog API.", "debug")
	}
	if Owner == "" {
		Owner = GetCurrentUsername()
	}
	Log("Required cli flags present.", "debug")
}

var (
	// KeyOutLocation This Consul key is the location we want to pull data from.
	// This data MUST be in the standard kvexpress structure of:
	//  /PrefixLocation/KeyOutLocation/data
	//  /PrefixLocation/KeyOutLocation/checksum
	KeyOutLocation string

	// FiletoWrite is the location we want to write the data to.
	FiletoWrite string

	// IgnoreStop is a special command to pull data EVEN if there's a stop key present.
	IgnoreStop bool
)

func init() {
	RootCmd.AddCommand(outCmd)
	outCmd.Flags().StringVarP(&KeyOutLocation, "key", "k", "", "key to pull data from")
	outCmd.Flags().StringVarP(&FiletoWrite, "file", "f", "", "where to write the data")
	outCmd.Flags().BoolVarP(&IgnoreStop, "ignore_stop", "", false, "ignore stop key")
}

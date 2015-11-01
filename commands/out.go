package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

var outCmd = &cobra.Command{
	Use:   "out",
	Short: "Write a file based on key data.",
	Long:  `out is for writing a file based on a Consul key.`,
	Run:   outRun,
}

func outRun(cmd *cobra.Command, args []string) {
	var Direction = "out"
	checkOutFlags(Direction)

	KeyData := kvexpress.KeyDataPath(KeyOutLocation, PrefixLocation, Direction)
	KeyChecksum := kvexpress.KeyChecksumPath(KeyOutLocation, PrefixLocation, Direction)
	KeyStop := kvexpress.KeyStopPath(KeyOutLocation, PrefixLocation, Direction)

	StopKeyData := kvexpress.Get(KeyStop, ConsulServer, Token, Direction)

	if StopKeyData != "" {
		log.Print(Direction, ": Stop Key is present - stopping. Reason: ", StopKeyData)
		os.Exit(1)
	} else {
		log.Print(Direction, ": Stop Key is NOT present - continuing.")
	}

	// Get the KV data out of Consul.
	KVData := kvexpress.Get(KeyData, ConsulServer, Token, Direction)

	// Get the Checksum data out of Consul.
	Checksum := kvexpress.Get(KeyChecksum, ConsulServer, Token, Direction)

	// Is the data long enough?
	longEnough := kvexpress.LengthCheck(KVData, MinFileLength, Direction)
	log.Print(Direction, ": longEnough='", strconv.FormatBool(longEnough), "'")

	// Does the checksum match?
	checksumMatch := kvexpress.ChecksumCompare(KVData, Checksum, Direction)
	log.Print(Direction, ": checksumMatch='", strconv.FormatBool(checksumMatch), "'")

	// If the data is long enough and the checksum matches, write the file.
	if longEnough && checksumMatch {
		kvexpress.WriteFile(KVData, FiletoWrite, FilePermissions, Direction)
		if DogStatsd {
			kvexpress.StatsdOut(KeyOutLocation)
		}
	} else {
		log.Print(Direction, ": Could not write file.")
	}

	// Run this command after the file is written.
	if PostExec != "" {
		log.Print(Direction, ": exec='", PostExec, "'")
		kvexpress.RunCommand(PostExec)
	}
}

func checkOutFlags(direction string) {
	log.Print(direction, ": Checking cli flags.")
	if KeyOutLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if FiletoWrite == "" {
		fmt.Println("Need a file to write in -f")
		os.Exit(1)
	}
	if DogStatsd {
		log.Print(direction, ": Enabling Dogstatsd metrics.")
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		log.Print(direction, ": Enabling Datadog API.")
		if os.Getenv("DATADOG_HOST") != "" {
			DatadogHost = os.Getenv("DATADOG_HOST")
			log.Print(direction, ": Using custom Datadog host: ", DatadogHost)
		}
	}
	log.Print(direction, ": Required cli flags present.")
}

var (
	KeyOutLocation string
	FiletoWrite    string
)

func init() {
	RootCmd.AddCommand(outCmd)
	outCmd.Flags().StringVarP(&KeyOutLocation, "key", "k", "", "key to pull data from")
	outCmd.Flags().StringVarP(&FiletoWrite, "file", "f", "", "where to write the data")
}

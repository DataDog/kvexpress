package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zorkian/go-datadog-api"
	"os"
	"time"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy a Consul key to another location.",
	Long:  `copy is for copying already existing keys.`,
	Run:   copyRun,
}

func copyRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	var dog = new(datadog.Client)
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	checkCopyFlags()

	// Set the source key locations.
	KeyData := KeyDataPath(KeyFrom, PrefixLocation)
	KeyChecksum := KeyChecksumPath(KeyFrom, PrefixLocation)

	c, _ := Connect(ConsulServer, Token)

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
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

	// If the data is long enough and the checksum matches, save to the new key location.
	if longEnough && checksumMatch {
		Log(fmt.Sprintf("copy='true' keyFrom='%s' keyTo='%s'", KeyFrom, KeyTo), "info")
		if Compress {
			KVData = CompressData(KVData)
		}
		// New destination key Locations
		KeyData = KeyDataPath(KeyTo, PrefixLocation)
		KeyChecksum = KeyChecksumPath(KeyTo, PrefixLocation)
		// Save it.
		saved := Set(c, KeyData, KVData, DogStatsd)
		if saved {
			KVDataBytes := len(KVData)
			Log(fmt.Sprintf("consul KeyData='%s' saved='true' size='%d'", KeyData, KVDataBytes), "info")
			Set(c, KeyChecksum, Checksum, DogStatsd)
			if DatadogAPIKey != "" && DatadogAPPKey != "" {
				DDCopyDataEvent(dog, KeyFrom, KeyTo)
			}
			if DogStatsd {
				StatsdIn(KeyTo, KVDataBytes, KVData)
			}
		} else {
			Log(fmt.Sprintf("consul KeyData='%s' saved='false'", KeyData), "info")
			RunTime(start, KeyTo, "consul_checksums_match", DogStatsd)
			os.Exit(0)
		}
	} else {
		Log(fmt.Sprintf("longEnough='%t' checksumMatch='%t'", longEnough, checksumMatch), "info")
		os.Exit(0)
	}

	// Run this command after the file is written.
	if PostExec != "" {
		Log(fmt.Sprintf("exec='%s'", PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, RawKeyOutLocation, "complete", DogStatsd)
}

func checkCopyFlags() {
	Log("Checking cli flags.", "debug")
	if KeyFrom == "" {
		fmt.Println("Need a key location in --keyfrom")
		os.Exit(1)
	}
	if KeyTo == "" {
		fmt.Println("Need a key destination in --keyto")
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
	// KeyFrom is the key in the Consul KV data store where we want to pull data.
	// This configuration variable will pull data from:
	//  /PrefixLocation/KeyFrom/data
	//  /PrefixLocation/KeyFrom/checksum
	KeyFrom string

	// KeyTo is the key in the Consul KV data store where we want to send data to.
	// This configuration variable will save data into:
	//  /PrefixLocation/KeyTo/data
	//  /PrefixLocation/KeyTo/checksum
	KeyTo string
)

func init() {
	RootCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVarP(&KeyFrom, "keyfrom", "", "", "key to pull data from")
	copyCmd.Flags().StringVarP(&KeyTo, "keyto", "", "", "key to write the data to")
}

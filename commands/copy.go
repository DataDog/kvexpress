// +build linux darwin freebsd

package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/zorkian/go-datadog-api.v1"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy a Consul key to another location.",
	Long:  `Copy is for copying already existing keys.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkCopyFlags()
		AutoEnable()
	},
	Run: copyRun,
}

func copyRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	var dog = new(datadog.Client)

	// Set the source key locations.
	KeyData := KeyPath(KeyFrom, "data")
	KeyChecksum := KeyPath(KeyFrom, "checksum")

	c, err := Connect(ConsulServer, Token)
	if err != nil {
		LogFatal("Could not connect to Consul.", KeyFrom, "consul_connect")
	}

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
	}

	// Get the KV data out of Consul.
	KVRaw, KVFlags := GetRaw(c, KeyData)

	// Decompress here if necessary.
	var KVData string
	if Compress {
		KVData = DecompressData(KVRaw, KVFlags)
	} else {
		KVData = string(KVRaw)
	}

	// Get the Checksum data out of Consul.
	Checksum := Get(c, KeyChecksum)

	// Is the data long enough?
	longEnough := LengthCheck(KVData, MinFileLength)
	Log(fmt.Sprintf("longEnough='%t'", longEnough), "debug")

	// Does the checksum match?
	checksumMatch := ChecksumCompare(KVData, Checksum)
	Log(fmt.Sprintf("checksumMatch='%t'", checksumMatch), "debug")

	// If the data is long enough and the checksum matches, save to the new key location.
	if longEnough && checksumMatch {
		Log(fmt.Sprintf("copy='true' keyFrom='%s' keyTo='%s'", KeyFrom, KeyTo), "info")
		var cFlags uint64
		if Compress {
			KVRaw, cFlags = CompressData(KVData)
		} else {
			KVRaw = []byte(KVData)
		}
		// New destination key Locations
		KeyData = KeyPath(KeyTo, "data")
		KeyChecksum = KeyPath(KeyTo, "checksum")
		// Save it.
		saved := SetRaw(c, KeyData, KVRaw, cFlags)
		if saved {
			KVDataBytes := len(KVRaw)
			Log(fmt.Sprintf("consul KeyData='%s' saved='true' size='%d'", KeyData, KVDataBytes), "info")
			Set(c, KeyChecksum, Checksum)
			if DatadogAPIKey != "" && DatadogAPPKey != "" {
				DDCopyDataEvent(dog, KeyFrom, KeyTo)
			}
			StatsdIn(KeyTo, KVDataBytes, KVData)
		} else {
			Log(fmt.Sprintf("consul KeyData='%s' saved='false'", KeyData), "info")
			RunTime(start, KeyTo, "consul_checksums_match")
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
	RunTime(start, KeyTo, "complete")
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

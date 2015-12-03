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
	var Direction = "copy"
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	checkCopyFlags(Direction)

	// Set the source key locations.
	KeyData := KeyDataPath(KeyFrom, PrefixLocation, Direction)
	KeyChecksum := KeyChecksumPath(KeyFrom, PrefixLocation, Direction)

	c, _ := Connect(ConsulServer, Token, Direction)

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
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
	Log(fmt.Sprintf("%s: longEnough='%t'", Direction, longEnough), "debug")

	// Does the checksum match?
	checksumMatch := ChecksumCompare(KVData, Checksum, Direction)
	Log(fmt.Sprintf("%s: checksumMatch='%t'", Direction, checksumMatch), "debug")

	// If the data is long enough and the checksum matches, save to the new key location.
	if longEnough && checksumMatch {
		Log(fmt.Sprintf("%s: copy='true' keyFrom='%s' keyTo='%s'", Direction, KeyFrom, KeyTo), "info")
		if Compress {
			KVData = CompressData(KVData, Direction)
		}
		// New destination key Locations
		KeyData = KeyDataPath(KeyTo, PrefixLocation, Direction)
		KeyChecksum = KeyChecksumPath(KeyTo, PrefixLocation, Direction)
		// Save it.
		saved := Set(c, KeyData, KVData, Direction, DogStatsd)
		if saved {
			KVDataBytes := len(KVData)
			Log(fmt.Sprintf("%s: consul KeyData='%s' saved='true' size='%d'", Direction, KeyData, KVDataBytes), "info")
			Set(c, KeyChecksum, Checksum, Direction, DogStatsd)
			if DatadogAPIKey != "" && DatadogAPPKey != "" {
				DDCopyDataEvent(dog, KeyFrom, KeyTo, Direction)
			}
			if DogStatsd {
				StatsdIn(KeyTo, KVDataBytes, KVData)
			}
		} else {
			Log(fmt.Sprintf("%s: consul KeyData='%s' saved='false'", Direction, KeyData), "info")
			RunTime(start, KeyTo, "consul_checksums_match", Direction, DogStatsd)
			os.Exit(0)
		}
	} else {
		Log(fmt.Sprintf("%s: longEnough='%t' checksumMatch='%t'", Direction, longEnough, checksumMatch), "info")
		os.Exit(0)
	}

	// Run this command after the file is written.
	if PostExec != "" {
		Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, RawKeyOutLocation, "complete", Direction, DogStatsd)
}

func checkCopyFlags(direction string) {
	Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if KeyFrom == "" {
		fmt.Println("Need a key location in --keyfrom")
		os.Exit(1)
	}
	if KeyTo == "" {
		fmt.Println("Need a key destination in --keyto")
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
	KeyFrom string
	KeyTo   string
)

func init() {
	RootCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVarP(&KeyFrom, "keyfrom", "", "", "key to pull data from")
	copyCmd.Flags().StringVarP(&KeyTo, "keyto", "", "", "key to write the data to")
}

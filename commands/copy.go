package commands

import (
	kvexpress "../kvexpress/"
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
	KeyData := kvexpress.KeyDataPath(KeyFrom, PrefixLocation, Direction)
	KeyChecksum := kvexpress.KeyChecksumPath(KeyFrom, PrefixLocation, Direction)

	c, _ := kvexpress.Connect(ConsulServer, Token, Direction)

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = kvexpress.DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
	}

	// Get the KV data out of Consul.
	KVData := kvexpress.Get(c, KeyData, Direction, DogStatsd)

	// Decompress here if necessary.
	if Compress {
		KVData = kvexpress.Decompress(KVData, Direction)
	}

	// Get the Checksum data out of Consul.
	Checksum := kvexpress.Get(c, KeyChecksum, Direction, DogStatsd)

	// Is the data long enough?
	longEnough := kvexpress.LengthCheck(KVData, MinFileLength, Direction)
	kvexpress.Log(fmt.Sprintf("%s: longEnough='%t'", Direction, longEnough), "debug")

	// Does the checksum match?
	checksumMatch := kvexpress.ChecksumCompare(KVData, Checksum, Direction)
	kvexpress.Log(fmt.Sprintf("%s: checksumMatch='%t'", Direction, checksumMatch), "debug")

	// If the data is long enough and the checksum matches, save to the new key location.
	if longEnough && checksumMatch {
		kvexpress.Log(fmt.Sprintf("%s: copy='true' keyFrom='%s' keyTo='%s'", Direction, KeyFrom, KeyTo), "info")
		if Compress {
			KVData = kvexpress.Compress(KVData, Direction)
		}
		// New destination key Locations
		KeyData = kvexpress.KeyDataPath(KeyTo, PrefixLocation, Direction)
		KeyChecksum = kvexpress.KeyChecksumPath(KeyTo, PrefixLocation, Direction)
		// Save it.
		saved := kvexpress.Set(c, KeyData, KVData, Direction, DogStatsd)
		if saved {
			KVDataBytes := len(KVData)
			kvexpress.Log(fmt.Sprintf("%s: consul KeyData='%s' saved='true' size='%d'", Direction, KeyData, KVDataBytes), "info")
			kvexpress.Set(c, KeyChecksum, Checksum, Direction, DogStatsd)
			if DatadogAPIKey != "" && DatadogAPPKey != "" {
				kvexpress.DDCopyDataEvent(dog, KeyFrom, KeyTo, Direction)
			}
			if DogStatsd {
				kvexpress.StatsdIn(KeyTo, KVDataBytes, KVData)
			}
		} else {
			kvexpress.Log(fmt.Sprintf("%s: consul KeyData='%s' saved='false'", Direction, KeyData), "info")
			kvexpress.RunTime(start, KeyTo, "consul_checksums_match", Direction, DogStatsd)
			os.Exit(0)
		}
	} else {
		kvexpress.Log(fmt.Sprintf("%s: longEnough='%t' checksumMatch='%t'", Direction, longEnough, checksumMatch), "info")
		os.Exit(0)
	}

	// Run this command after the file is written.
	if PostExec != "" {
		kvexpress.Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, RawKeyOutLocation, "complete", Direction, DogStatsd)
}

func checkCopyFlags(direction string) {
	kvexpress.Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if KeyFrom == "" {
		fmt.Println("Need a key location in --keyfrom")
		os.Exit(1)
	}
	if KeyTo == "" {
		fmt.Println("Need a key destination in --keyto")
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
	KeyFrom string
	KeyTo   string
)

func init() {
	RootCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVarP(&KeyFrom, "keyfrom", "", "", "key to pull data from")
	copyCmd.Flags().StringVarP(&KeyTo, "keyto", "", "", "key to write the data to")
}

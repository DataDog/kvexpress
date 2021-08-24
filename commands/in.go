// +build linux darwin freebsd

package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/zorkian/go-datadog-api.v1"
)

var inCmd = &cobra.Command{
	Use:   "in",
	Short: "Put configuration into Consul.",
	Long:  `In is for putting data into a Consul key so that you can write it on another networked node.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkInFlags()
		AutoEnable()
	},
	Run: inRun,
}

func inRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	var dog = new(datadog.Client)
	var CompareFile = ""
	var LastFile = ""
	var FileString = ""

	KeyStop := KeyPath(KeyInLocation, "stop")
	KeyData := KeyPath(KeyInLocation, "data")
	KeyChecksum := KeyPath(KeyInLocation, "checksum")

	if FiletoRead != "" {
		CompareFile = CompareFilename(FiletoRead)
		LastFile = LastFilename(FiletoRead)
	} else {
		CompareFile = RandomTmpFile()
		LastFile = LastFilename(CompareFile)
	}

	// Let's double check those files are safe to write.
	CheckFiletoWrite(CompareFile, "")
	CheckFiletoWrite(LastFile, "")

	c, err := Connect(ConsulServer, Token)
	if err != nil {
		LogFatal("Could not connect to Consul.", KeyInLocation, "consul_connect")
	}

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
	}

	StopKeyData := Get(c, KeyStop)

	if StopKeyData != "" {
		Log(fmt.Sprintf("Stop Key is present - stopping. Reason: %s", StopKeyData), "info")
		if DatadogAPIKey != "" && DatadogAPPKey != "" {
			DDStopEvent(dog, KeyStop, StopKeyData)
		}
		RunTime(start, KeyInLocation, "stop_key")
		os.Exit(1)
	} else {
		Log("Stop Key is NOT present - continuing.", "info")
	}

	// Read the file - if it's to be sorted - then make sure to sort.
	if FiletoRead != "" {
		FileString = ReadFile(FiletoRead)
	} else {
		FileString = ReadURL(UrltoRead)
	}

	// Sorting also removes any blank lines.
	if Sorted {
		FileString = SortFile(FileString)
	}

	// Is it long enough?
	longEnough := LengthCheck(FileString, MinFileLength)

	if !longEnough {
		Log("File is NOT long enough. Stopping.", "info")
		if DatadogAPIKey != "" && DatadogAPPKey != "" {
			DDLengthEvent(dog, KeyInLocation, FileString)
		}
		RunTime(start, KeyInLocation, "not_long_enough")
		os.Exit(1)
	}

	// Write the .compare file.
	WriteFile(FileString, CompareFile, FilePermissions, Owner)

	// Check for the .last file - touch if it doesn't exist.
	CheckLastFile(LastFile, FilePermissions, Owner)

	// Read compare and last files into string.
	CompareData := ReadFile(CompareFile)
	LastData := ReadFile(LastFile)

	if CompareData != "" && LastData != "" {
		Log("We have data - let's do the thing.", "info")
	} else {
		Log("We do NOT have data. This should never happen.", "info")
		RunTime(start, KeyInLocation, "error_no_data")
		os.Exit(1)
	}

	// Get SHA256 values for each string.
	CompareChecksum := ComputeChecksum(CompareData)
	LastChecksum := ComputeChecksum(LastData)

	// If they're different - let's update things.
	if CompareChecksum != LastChecksum {
		Log("file checksum='different' update='true'", "info")
	} else {
		Log("file checksum='match' update='false'", "info")
		RunTime(start, KeyInLocation, "file_checksums_match")
		os.Exit(0)
	}

	// Diff the files.
	diff := UnixDiff(LastFile, CompareFile)

	// If we get this far - copy the CompareData to the .last file.
	// This handles the case detailed in https://github.com/darron/kvexpress/issues/33
	WriteFile(CompareData, LastFile, FilePermissions, Owner)

	// Get the checksum from Consul.
	CurrentChecksum := Get(c, KeyChecksum)

	if CurrentChecksum != CompareChecksum {
		Log("consul checksum='different' update='true'", "info")
		// Compress data here.
		var CompareRaw []byte
		if Compress {
			CompareRaw = CompressData(CompareData)
		} else {
			CompareRaw = []byte(CompareData)
		}
		saved := SetRaw(c, KeyData, CompareRaw)
		if saved {
			CompareDataBytes := len(CompareRaw)
			Log(fmt.Sprintf("consul KeyData='%s' saved='true' size='%d'", KeyData, CompareDataBytes), "info")
			Set(c, KeyChecksum, CompareChecksum)
			if DatadogAPIKey != "" && DatadogAPPKey != "" {
				DDSaveDataEvent(dog, KeyData, diff)
			}
			StatsdIn(KeyInLocation, CompareDataBytes, CompareData)

			if UrltoRead != "" {
				urlOutput := fmt.Sprintf("\nURL: %s\n\nWhat was inserted into: '%s'\n===================\n%s\n===================\n", UrltoRead, KeyData, CompareData)
				fmt.Println(urlOutput)
			}

		} else {
			Log(fmt.Sprintf("consul KeyData='%s' saved='false'", KeyData), "info")
			RunTime(start, KeyInLocation, "consul_checksums_match")
			os.Exit(0)
		}
	} else {
		Log("consul checksum='match' update='false'", "info")
	}
	// Run this command after the data is input.
	if PostExec != "" {
		Log(fmt.Sprintf("exec='%s'", PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, KeyInLocation, "complete")
}

func checkInFlags() {
	Log("Checking cli flags.", "debug")
	if KeyInLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if FiletoRead == "" && UrltoRead == "" {
		fmt.Println("Need a file -f or url -u to read from.")
		os.Exit(1)
	}
	if FiletoRead != "" {
		if _, err := os.Stat(FiletoRead); err != nil {
			fmt.Println("File ", FiletoRead, " does not exist.")
			os.Exit(1)
		}
	}
	if FiletoRead != "" && UrltoRead != "" {
		fmt.Println("You cannot use both -f and -u.")
		os.Exit(1)
	}
	Log("Required cli flags present.", "debug")
}

var (
	// KeyInLocation is the key in the Consul KV data store where we want to store the data.
	// This configuration variable will save data into:
	//  /PrefixLocation/KeyInLocation/data
	//  /PrefixLocation/KeyInLocation/checksum
	KeyInLocation string

	// FiletoRead is the file to read to get the data.
	FiletoRead string

	// Sorted is an option to sort the file alphabetically. Doesn't work on many types
	// of files. But works great on files with many blank lines where ordering doesn't matter.
	Sorted bool

	// UrltoRead is an HTTP URL to read data from using ReadURL().
	UrltoRead string
)

func init() {
	RootCmd.AddCommand(inCmd)
	inCmd.Flags().StringVarP(&KeyInLocation, "key", "k", "", "key to push data to")
	inCmd.Flags().StringVarP(&FiletoRead, "file", "f", "", "filename to read data from")
	inCmd.Flags().StringVarP(&UrltoRead, "url", "u", "", "url to read data from")
	inCmd.Flags().BoolVarP(&Sorted, "sorted", "S", false, "sort the input file")
}

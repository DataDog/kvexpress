package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zorkian/go-datadog-api"
	"os"
	"time"
)

var inCmd = &cobra.Command{
	Use:   "in",
	Short: "Put configuration into Consul.",
	Long:  `in is for putting data into a Consul key so that you can write it on another networked node.`,
	Run:   inRun,
}

func inRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	var dog = new(datadog.Client)
	var CompareFile = ""
	var LastFile = ""
	var FileString = ""
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	checkInFlags(Direction)

	KeyStop := KeyStopPath(KeyInLocation, PrefixLocation, Direction)
	KeyData := KeyDataPath(KeyInLocation, PrefixLocation, Direction)
	KeyChecksum := KeyChecksumPath(KeyInLocation, PrefixLocation, Direction)

	if FiletoRead != "" {
		CompareFile = CompareFilename(FiletoRead, Direction)
		LastFile = LastFilename(FiletoRead, Direction)
	} else {
		CompareFile = RandomTmpFile(Direction)
		LastFile = LastFilename(CompareFile, Direction)
	}

	// Let's double check those files are safe to write.
	CheckFiletoWrite(CompareFile, "", Direction)
	CheckFiletoWrite(LastFile, "", Direction)

	c, _ := Connect(ConsulServer, Token, Direction)

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
	}

	StopKeyData := Get(c, KeyStop, Direction, DogStatsd)

	if StopKeyData != "" {
		Log(fmt.Sprintf("%s: Stop Key is present - stopping. Reason: %s", Direction, StopKeyData), "info")
		if DatadogAPIKey != "" && DatadogAPPKey != "" {
			DDStopEvent(dog, KeyStop, StopKeyData, Direction)
		}
		RunTime(start, KeyInLocation, "stop_key", Direction, DogStatsd)
		os.Exit(1)
	} else {
		Log(fmt.Sprintf("%s: Stop Key is NOT present - continuing.", Direction), "info")
	}

	// Read the file - if it's to be sorted - then make sure to sort.
	if FiletoRead != "" {
		FileString = ReadFile(FiletoRead)
	} else {
		FileString = ReadUrl(UrltoRead, DogStatsd)
	}

	// Sorting also removes any blank lines.
	if Sorted {
		FileString = SortFile(FileString)
	}

	// Is it long enough?
	longEnough := LengthCheck(FileString, MinFileLength, Direction)

	if !longEnough {
		Log(fmt.Sprintf("%s: File is NOT long enough. Stopping.", Direction), "info")
		// TODO: Add Datadog Event here.
		RunTime(start, KeyInLocation, "not_long_enough", Direction, DogStatsd)
		os.Exit(1)
	}

	// Write the .compare file.
	WriteFile(FileString, CompareFile, FilePermissions, Owner, Direction, DogStatsd)

	// Check for the .last file - touch if it doesn't exist.
	CheckLastFile(LastFile, FilePermissions, Owner, DogStatsd)

	// Read compare and last files into string.
	CompareData := ReadFile(CompareFile)
	LastData := ReadFile(LastFile)

	if CompareData != "" && LastData != "" {
		Log(fmt.Sprintf("%s: We have data - let's do the thing.", Direction), "info")
	} else {
		Log(fmt.Sprintf("%s: We do NOT have data. This should never happen.", Direction), "info")
		RunTime(start, KeyInLocation, "error_no_data", Direction, DogStatsd)
		os.Exit(1)
	}

	// Get SHA256 values for each string.
	CompareChecksum := ComputeChecksum(CompareData, Direction)
	LastChecksum := ComputeChecksum(LastData, Direction)

	// If they're different - let's update things.
	if CompareChecksum != LastChecksum {
		Log(fmt.Sprintf("%s: file checksum='different' update='true'", Direction), "info")
	} else {
		Log(fmt.Sprintf("%s: file checksum='match' update='false'", Direction), "info")
		RunTime(start, KeyInLocation, "file_checksums_match", Direction, DogStatsd)
		os.Exit(0)
	}

	// Diff the files.
	diff := UnixDiff(LastFile, CompareFile)

	// If we get this far - copy the CompareData to the .last file.
	// This handles the case detailed in https://github.com/darron/kvexpress/issues/33
	WriteFile(CompareData, LastFile, FilePermissions, Owner, Direction, DogStatsd)

	// Get the checksum from Consul.
	CurrentChecksum := Get(c, KeyChecksum, Direction, DogStatsd)

	if CurrentChecksum != CompareChecksum {
		Log(fmt.Sprintf("%s: consul checksum='different' update='true'", Direction), "info")
		// Compress data here.
		if Compress {
			CompareData = CompressData(CompareData, Direction)
		}
		saved := Set(c, KeyData, CompareData, Direction, DogStatsd)
		if saved {
			CompareDataBytes := len(CompareData)
			Log(fmt.Sprintf("%s: consul KeyData='%s' saved='true' size='%d'", Direction, KeyData, CompareDataBytes), "info")
			Set(c, KeyChecksum, CompareChecksum, Direction, DogStatsd)
			if DatadogAPIKey != "" && DatadogAPPKey != "" {
				DDSaveDataEvent(dog, KeyData, diff, Direction)
			}

			if DogStatsd {
				StatsdIn(KeyInLocation, CompareDataBytes, CompareData)
			}

			if UrltoRead != "" {
				urlOutput := fmt.Sprintf("\nURL: %s\n\nWhat was inserted into: '%s'\n===================\n%s\n===================\n", UrltoRead, KeyData, CompareData)
				fmt.Println(urlOutput)
			}

		} else {
			Log(fmt.Sprintf("%s: consul KeyData='%s' saved='false'", Direction, KeyData), "info")
			RunTime(start, KeyInLocation, "consul_checksums_match", Direction, DogStatsd)
			os.Exit(0)
		}
	} else {
		Log(fmt.Sprintf("%s: consul checksum='match' update='false'", Direction), "info")
	}
	// Run this command after the data is input.
	if PostExec != "" {
		Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, KeyInLocation, "complete", Direction, DogStatsd)
}

func checkInFlags(direction string) {
	Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if KeyInLocation == "" {
		fmt.Println(direction, ": Need a key location in -k")
		os.Exit(1)
	}
	if FiletoRead == "" && UrltoRead == "" {
		fmt.Println(direction, ": Need a file -f or url -u to read from.")
		os.Exit(1)
	}
	if FiletoRead != "" {
		if _, err := os.Stat(FiletoRead); err != nil {
			fmt.Println(direction, ": File ", FiletoRead, " does not exist.")
			os.Exit(1)
		}
	}
	if FiletoRead != "" && UrltoRead != "" {
		fmt.Println(direction, ": You cannot use both -f and -u.")
		os.Exit(1)
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		Log(fmt.Sprintf("%s: Enabling Datadog API.", direction), "debug")
	}
	if DogStatsd {
		Log(fmt.Sprintf("%s: Enabling Dogstatsd metrics.", direction), "debug")
	}
	if Owner == "" {
		Owner = GetCurrentUsername(direction)
	}
	Log(fmt.Sprintf("%s: Required cli flags present.", direction), "debug")
}

var (
	KeyInLocation string
	FiletoRead    string
	Sorted        bool
	UrltoRead     string
)

func init() {
	RootCmd.AddCommand(inCmd)
	inCmd.Flags().StringVarP(&KeyInLocation, "key", "k", "", "key to push data to")
	inCmd.Flags().StringVarP(&FiletoRead, "file", "f", "", "filename to read data from")
	inCmd.Flags().StringVarP(&UrltoRead, "url", "u", "", "url to read data from")
	inCmd.Flags().BoolVarP(&Sorted, "sorted", "S", false, "sort the input file")
}

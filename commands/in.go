package commands

import (
	kvexpress "../kvexpress/"
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
	var Direction = "in"
	var CompareFile = ""
	var LastFile = ""
	var FileString = ""
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	checkInFlags(Direction)

	KeyStop := kvexpress.KeyStopPath(KeyInLocation, PrefixLocation, Direction)
	KeyData := kvexpress.KeyDataPath(KeyInLocation, PrefixLocation, Direction)
	KeyChecksum := kvexpress.KeyChecksumPath(KeyInLocation, PrefixLocation, Direction)

	if FiletoRead != "" {
		CompareFile = kvexpress.CompareFilename(FiletoRead, Direction)
		LastFile = kvexpress.LastFilename(FiletoRead, Direction)
	} else {
		CompareFile = kvexpress.RandomTmpFile(Direction)
		LastFile = kvexpress.LastFilename(CompareFile, Direction)
	}

	// Let's double check those files are safe to write.
	kvexpress.CheckFiletoWrite(CompareFile, "", Direction)
	kvexpress.CheckFiletoWrite(LastFile, "", Direction)

	c, _ := kvexpress.Connect(ConsulServer, Token, Direction)

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = kvexpress.DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
	}

	StopKeyData := kvexpress.Get(c, KeyStop, Direction, DogStatsd)

	if StopKeyData != "" {
		kvexpress.Log(fmt.Sprintf("%s: Stop Key is present - stopping. Reason: %s", Direction, StopKeyData), "info")
		if DatadogAPIKey != "" && DatadogAPPKey != "" {
			kvexpress.DDStopEvent(dog, KeyStop, StopKeyData, Direction)
		}
		kvexpress.RunTime(start, KeyInLocation, "stop_key", Direction, DogStatsd)
		os.Exit(1)
	} else {
		kvexpress.Log(fmt.Sprintf("%s: Stop Key is NOT present - continuing.", Direction), "info")
	}

	// Read the file - if it's to be sorted - then make sure to sort.
	if FiletoRead != "" {
		FileString = kvexpress.ReadFile(FiletoRead)
	} else {
		FileString = kvexpress.ReadUrl(UrltoRead, DogStatsd)
	}

	// Sorting also removes any blank lines.
	if Sorted {
		FileString = kvexpress.SortFile(FileString)
	}

	// Is it long enough?
	longEnough := kvexpress.LengthCheck(FileString, MinFileLength, Direction)

	if !longEnough {
		kvexpress.Log(fmt.Sprintf("%s: File is NOT long enough. Stopping.", Direction), "info")
		// TODO: Add Datadog Event here.
		kvexpress.RunTime(start, KeyInLocation, "not_long_enough", Direction, DogStatsd)
		os.Exit(1)
	}

	// Write the .compare file.
	kvexpress.WriteFile(FileString, CompareFile, FilePermissions, Owner, Direction)

	// Check for the .last file - touch if it doesn't exist.
	kvexpress.CheckLastFile(LastFile, FilePermissions, Owner)

	// Read compare and last files into string.
	CompareData := kvexpress.ReadFile(CompareFile)
	LastData := kvexpress.ReadFile(LastFile)

	if CompareData != "" && LastData != "" {
		kvexpress.Log(fmt.Sprintf("%s: We have data - let's do the thing.", Direction), "info")
	} else {
		kvexpress.Log(fmt.Sprintf("%s: We do NOT have data. This should never happen.", Direction), "info")
		kvexpress.RunTime(start, KeyInLocation, "error_no_data", Direction, DogStatsd)
		os.Exit(1)
	}

	// Get SHA256 values for each string.
	CompareChecksum := kvexpress.ComputeChecksum(CompareData, Direction)
	LastChecksum := kvexpress.ComputeChecksum(LastData, Direction)

	// If they're different - let's update things.
	if CompareChecksum != LastChecksum {
		kvexpress.Log(fmt.Sprintf("%s: file checksum='different' update='true'", Direction), "info")
	} else {
		kvexpress.Log(fmt.Sprintf("%s: file checksum='match' update='false'", Direction), "info")
		kvexpress.RunTime(start, KeyInLocation, "file_checksums_match", Direction, DogStatsd)
		os.Exit(0)
	}

	// Diff the files.
	diff := kvexpress.UnixDiff(LastFile, CompareFile)

	// If we get this far - copy the CompareData to the .last file.
	// This handles the case detailed in https://github.com/darron/kvexpress/issues/33
	kvexpress.WriteFile(CompareData, LastFile, FilePermissions, Owner, Direction)

	// Get the checksum from Consul.
	CurrentChecksum := kvexpress.Get(c, KeyChecksum, Direction, DogStatsd)

	if CurrentChecksum != CompareChecksum {
		kvexpress.Log(fmt.Sprintf("%s: consul checksum='different' update='true'", Direction), "info")
		saved := kvexpress.Set(c, KeyData, CompareData, Direction, DogStatsd)
		if saved {
			CompareDataBytes := len(CompareData)
			kvexpress.Log(fmt.Sprintf("%s: consul KeyData='%s' saved='true' size='%d'", Direction, KeyData, CompareDataBytes), "info")
			kvexpress.Set(c, KeyChecksum, CompareChecksum, Direction, DogStatsd)
			if DatadogAPIKey != "" && DatadogAPPKey != "" {
				kvexpress.DDSaveDataEvent(dog, KeyData, diff, Direction)
			}

			if DogStatsd {
				kvexpress.StatsdIn(KeyInLocation, CompareDataBytes, CompareData)
			}

			if UrltoRead != "" {
				urlOutput := fmt.Sprintf("\nURL: %s\n\nWhat was inserted into: '%s'\n===================\n%s\n===================\n", UrltoRead, KeyData, CompareData)
				fmt.Println(urlOutput)
			}

		} else {
			kvexpress.Log(fmt.Sprintf("%s: consul KeyData='%s' saved='false'", Direction, KeyData), "info")
			kvexpress.RunTime(start, KeyInLocation, "consul_checksums_match", Direction, DogStatsd)
			os.Exit(0)
		}
	} else {
		kvexpress.Log(fmt.Sprintf("%s: consul checksum='match' update='false'", Direction), "info")
	}
	// Run this command after the data is input.
	if PostExec != "" {
		kvexpress.Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, KeyInLocation, "complete", Direction, DogStatsd)
}

func checkInFlags(direction string) {
	kvexpress.Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
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
		kvexpress.Log(fmt.Sprintf("%s: Enabling Datadog API.", direction), "debug")
	}
	if DogStatsd {
		kvexpress.Log(fmt.Sprintf("%s: Enabling Dogstatsd metrics.", direction), "debug")
	}
	if Owner == "" {
		Owner = kvexpress.GetCurrentUsername(direction)
	}
	kvexpress.Log(fmt.Sprintf("%s: Required cli flags present.", direction), "debug")
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

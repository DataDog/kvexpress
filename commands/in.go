package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/zorkian/go-datadog-api"
	"github.com/spf13/cobra"
	"log"
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
	checkInFlags(Direction)
	if EnvVars {
		ConfigEnvVars(Direction)
	}

	KeyStop := kvexpress.KeyStopPath(KeyInLocation, PrefixLocation, Direction)
	KeyData := kvexpress.KeyDataPath(KeyInLocation, PrefixLocation, Direction)
	KeyChecksum := kvexpress.KeyChecksumPath(KeyInLocation, PrefixLocation, Direction)

	CompareFile := kvexpress.CompareFilename(FiletoRead, Direction)
	LastFile := kvexpress.LastFilename(FiletoRead, Direction)

	// Let's double check those files are safe to write.
	kvexpress.CheckFiletoWrite(CompareFile, "", Direction)
	kvexpress.CheckFiletoWrite(LastFile, "", Direction)

	c, _ := kvexpress.Connect(ConsulServer, Token, Direction)

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = kvexpress.DDAPIConnect(DatadogAPIKey, DatadogAPPKey, DatadogHost)
	}

	StopKeyData := kvexpress.Get(c, KeyStop, Direction)

	if StopKeyData != "" {
		log.Print(Direction, ": Stop Key is present - stopping. Reason: ", StopKeyData)
		if DatadogAPIKey != "" && DatadogAPPKey != "" {
			kvexpress.DDStopEvent(dog, KeyStop, StopKeyData, Direction)
		}
		kvexpress.RunTime(start, "stop_key", Direction)
		os.Exit(1)
	} else {
		log.Print(Direction, ": Stop Key is NOT present - continuing.")
	}

	// Read the file - if it's to be sorted - then make sure to sort.
	FileString := kvexpress.ReadFile(FiletoRead)

	// Sorting also removes any blank lines.
	if Sorted {
		FileString = kvexpress.SortFile(FileString)
	}

	// Is it long enough?
	longEnough := kvexpress.LengthCheck(FileString, MinFileLength, Direction)

	if !longEnough {
		log.Print(Direction, ": File is NOT long enough. Stopping.")
		kvexpress.RunTime(start, "not_long_enough", Direction)
		os.Exit(1)
	}

	// Write the .compare file.
	kvexpress.WriteFile(FileString, CompareFile, FilePermissions, Direction)

	// Check for the .last file - touch if it doesn't exist.
	kvexpress.CheckLastFile(LastFile, FilePermissions)

	// Read compare and last files into string.
	CompareData := kvexpress.ReadFile(CompareFile)
	LastData := kvexpress.ReadFile(LastFile)

	if CompareData != "" && LastData != "" {
		log.Print(Direction, ": We have data - let's do the thing.")
	} else {
		log.Print(Direction, ": We do NOT have data. This should never happen.")
		kvexpress.RunTime(start, "error_no_data", Direction)
		os.Exit(1)
	}

	// Get SHA256 values for each string.
	CompareChecksum := kvexpress.ComputeChecksum(CompareData, Direction)
	LastChecksum := kvexpress.ComputeChecksum(LastData, Direction)

	// If they're different - let's update things.
	if CompareChecksum != LastChecksum {
		log.Print(Direction, ": file checksums are different - let's update some stuff!")
	} else {
		log.Print(Direction, ": checksums='match' saved='false'")
		kvexpress.RunTime(start, "file_checksums_match", Direction)
		os.Exit(0)
	}

	// If we get this far - copy the CompareData to the .last file.
	// This handles the case detailed in https://github.com/darron/kvexpress/issues/33
	kvexpress.WriteFile(CompareData, LastFile, FilePermissions, Direction)

	// Diff the file data.
	diff := kvexpress.Diff(LastData, CompareData)

	// Get the checksum from Consul.
	CurrentChecksum := kvexpress.Get(c, KeyChecksum, Direction)

	if CurrentChecksum != CompareChecksum {
		log.Print(Direction, ": current and previous Consul checksum are different - let's update the KV store.")
		saved := kvexpress.Set(c, KeyData, CompareData, Direction)
		if saved {
			CompareDataBytes := len(CompareData)
			log.Print(Direction, ": KeyData='", KeyData, "' saved='true' size='", CompareDataBytes, "'")
			kvexpress.Set(c, KeyChecksum, CompareChecksum, Direction)
			if DatadogAPIKey != "" && DatadogAPPKey != "" {
				kvexpress.DDSaveDataEvent(dog, KeyData, diff, Direction)
			}

			if DogStatsd {
				kvexpress.StatsdIn(KeyInLocation, CompareDataBytes, CompareData)
			}

		} else {
			log.Print(Direction, ": KeyData='", KeyData, "' saved='false'")
			kvexpress.RunTime(start, "consul_checksums_match", Direction)
			os.Exit(0)
		}

	}
	// Run this command after the data is input.
	if PostExec != "" {
		log.Print(Direction, ": exec='", PostExec, "'")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, "complete", Direction)
}

func checkInFlags(direction string) {
	log.Print(direction, ": Checking cli flags.")
	if KeyInLocation == "" {
		fmt.Println(direction, ": Need a key location in -k")
		os.Exit(1)
	}
	if FiletoRead == "" {
		fmt.Println(direction, ": Need a file to read in -f")
		os.Exit(1)
	}
	if _, err := os.Stat(FiletoRead); err != nil {
		fmt.Println(direction, ": File ", FiletoRead, " does not exist.")
		os.Exit(1)
	}
	if DogStatsd {
		log.Print(direction, ": Enabling Dogstatsd metrics.")
	}
	log.Print(direction, ": Required cli flags present.")
}

var (
	KeyInLocation string
	FiletoRead    string
	Sorted        bool
)

func init() {
	RootCmd.AddCommand(inCmd)
	inCmd.Flags().StringVarP(&KeyInLocation, "key", "k", "", "key to push data to")
	inCmd.Flags().StringVarP(&FiletoRead, "file", "f", "", "filename to read data from")
	inCmd.Flags().BoolVarP(&Sorted, "sorted", "S", false, "sort the input file")
}

package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var inCmd = &cobra.Command{
	Use:   "in",
	Short: "Put configuration into Consul.",
	Long:  `in is for putting data into a Consul key so that you can write it on another networked node.`,
	Run:   inRun,
}

func inRun(cmd *cobra.Command, args []string) {
	var Direction = "in"
	checkInFlags(Direction)

	key_stop := kvexpress.KeyStopPath(KeyInLocation, PrefixLocation, Direction)
	key_data := kvexpress.KeyDataPath(KeyInLocation, PrefixLocation, Direction)
	key_checksum := kvexpress.KeyChecksumPath(KeyInLocation, PrefixLocation, Direction)

	compare_file := kvexpress.CompareFilename(FiletoRead)
	last_file := kvexpress.LastFilename(FiletoRead)

	StopKeyData := kvexpress.Get(key_stop, ConsulServer, Token, Direction)

	if StopKeyData != "" {
		log.Print(Direction, ": Stop Key is present - stopping. Reason: ", StopKeyData)
		os.Exit(1)
	} else {
		log.Print(Direction, ": Stop Key is NOT present - continuing.")
	}

	// Read the file - if it's to be sorted - then make sure to sort.
	// TODO: Do we need to uniq as well?
	file_string := kvexpress.ReadFile(FiletoRead)
	log.Print(Direction, ": file_string='", file_string, "'")

	if Sorted {
		file_string = kvexpress.SortFile(file_string)
	}

	// Is it long enough?
	longEnough := kvexpress.LengthCheck(file_string, MinFileLength, Direction)

	if !longEnough {
		log.Print(Direction, ": File is NOT long enough. Stopping.")
		os.Exit(1)
	}

	// Write the .compare file.
	kvexpress.WriteFile(file_string, compare_file, FilePermissions, Direction)

	// Check for the .last file - touch if it doesn't exist.
	kvexpress.CheckLastFile(last_file, FilePermissions)

	// Read compare and last files into string.
	compare_data := kvexpress.ReadFile(compare_file)
	last_data := kvexpress.ReadFile(last_file)

	if compare_data != "" && last_data != "" {
		log.Print(Direction, ": We have data - let's do the thing.")
	} else {
		log.Print(Direction, ": We do NOT have data. This should never happen.")
		os.Exit(1)
	}

	// Get SHA256 values for each string.
	compare_checksum := kvexpress.ComputeChecksum(compare_data, Direction)
	last_checksum := kvexpress.ComputeChecksum(last_data, Direction)

	// If they're different - let's update things.
	if compare_checksum != last_checksum {
		log.Print(Direction, ": file checksums are different - let's update some stuff!")
	} else {
		log.Print(Direction, ": checksums='match' saved='false'")
		os.Exit(0)
	}

	// Diff the file data.
	// html_diff := kvexpress.HTMLDiff(last_data, compare_data)

	// TODO: To be removed.
	// fmt.Printf("%v", html_diff)

	// Get the checksum from Consul.
	current_checksum := kvexpress.Get(key_checksum, ConsulServer, Token, Direction)

	if current_checksum != compare_checksum {
		log.Print(Direction, ": current and previous Consul checksum are different - let's update the KV store.")
		saved := kvexpress.Set(key_data, compare_data, ConsulServer, Token, Direction)
		if saved {
			compare_data_bytes := len(compare_data)
			log.Print(Direction, ": key_data='", key_data, "' saved='true' size='", compare_data_bytes, "'")
			kvexpress.Set(key_checksum, compare_checksum, ConsulServer, Token, Direction)

			if DogStatsd {
				kvexpress.StatsdIn(KeyInLocation, compare_data_bytes, compare_data)
			}

			// Copy the compare_data to the .last file.
			kvexpress.WriteFile(compare_data, last_file, FilePermissions, Direction)
		} else {
			log.Print(Direction, ": key_data='", key_data, "' saved='false'")
			os.Exit(1)
		}

	}

	// Run this command after the data is input.
	if PostExec != "" {
		log.Print(Direction, ": exec='", PostExec, "'")
		kvexpress.RunCommand(PostExec)
	}
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

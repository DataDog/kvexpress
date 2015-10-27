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
	checkInFlags()

	key_stop := kvexpress.KeyStopPath(KeyInLocation, PrefixLocation, Direction)

	StopKeyData := kvexpress.Get(key_stop, ConsulServer, Token, Direction)

	if StopKeyData != "" {
		log.Print(Direction, ": Stop Key is present.")
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

	// Write the .compare file.
	compare_file := kvexpress.CompareFilename(FiletoRead)
	kvexpress.WriteFile(file_string, compare_file, FilePermissions, Direction)

	// key_data := kvexpress.KeyDataPath(KeyLocation, PrefixLocation)
	// key_checksum := kvexpress.KeyChecksumPath(KeyLocation, PrefixLocation)

	// Run this command after the data is input.
	if PostExec != "" {
		log.Print(Direction, ": exec='", PostExec, "'")
		kvexpress.RunCommand(PostExec)
	}
}

func checkInFlags() {
	log.Print("in: Checking cli flags.")
	if KeyInLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if FiletoRead == "" {
		fmt.Println("Need a file to read in -f")
		os.Exit(1)
	}
	if _, err := os.Stat(FiletoRead); err != nil {
		fmt.Println("File ", FiletoRead, " does not exist.")
		os.Exit(1)
	}
	log.Print("in: Required cli flags present.")
}

var KeyInLocation string
var FiletoRead string
var Sorted bool

func init() {
	RootCmd.AddCommand(inCmd)
	inCmd.Flags().StringVarP(&KeyInLocation, "key", "k", "", "key to push data to")
	inCmd.Flags().StringVarP(&FiletoRead, "file", "f", "", "filename to read data from")
	inCmd.Flags().BoolVarP(&Sorted, "sorted", "S", false, "sort the input file")
}

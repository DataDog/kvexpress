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
		log.Print("in: Stop Key is present.")
		os.Exit(1)
	}

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
	log.Print("in: Required cli flags present.")
}

var KeyInLocation string

func init() {
	RootCmd.AddCommand(inCmd)
	inCmd.Flags().StringVarP(&KeyInLocation, "key", "k", "", "key to pull data from")
}

package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Put stop value into Consul.",
	Long:  `stop is a convenient way to put stop values in Consul.`,
	Run:   stopRun,
}

func stopRun(cmd *cobra.Command, args []string) {
	var Direction = "stop"
	checkStopFlags(Direction)

  KeyStop := kvexpress.KeyStopPath(KeyStopLocation, PrefixLocation, Direction)

  saved := kvexpress.Set(KeyStop, KeyStopReason, ConsulServer, Token, Direction)

  if saved {
    log.Print(Direction, ": KeyStop='", KeyStop, "' saved='true' KeyStopReason='", KeyStopReason, "'")
  }

	// Run this command after the key is stopped.
	if PostExec != "" {
		log.Print(Direction, ": exec='", PostExec, "'")
		kvexpress.RunCommand(PostExec)
	}
}

func checkStopFlags(direction string) {
	log.Print(direction, ": Checking cli flags.")
	if KeyStopLocation == "" {
		fmt.Println("Need a key to stop in -k")
		os.Exit(1)
	}
  if KeyStopReason == "" {
    fmt.Println("Need a reason to stop in -r")
    os.Exit(1)
  }
	log.Print(direction, ": Required cli flags present.")
}

var (
	KeyStopLocation string
  KeyStopReason string
)

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVarP(&KeyStopLocation, "key", "k", "", "key to stop")
  stopCmd.Flags().StringVarP(&KeyStopReason, "reason", "r", "", "reason to stop")
}

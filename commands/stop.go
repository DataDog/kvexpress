package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/darron/go-datadog-api"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Put stop value into Consul.",
	Long:  `stop is a convenient way to put stop values in Consul.`,
	Run:   stopRun,
}

func stopRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	var dog = new(datadog.Client)
	var Direction = "stop"
	checkStopFlags(Direction)
	if EnvVars {
		ConfigEnvVars(Direction)
	}

	KeyStop := kvexpress.KeyStopPath(KeyStopLocation, PrefixLocation, Direction)

	c, _ := kvexpress.Connect(ConsulServer, Token, Direction)

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = kvexpress.DDAPIConnect(DatadogAPIKey, DatadogAPPKey, DatadogHost)
	}

	saved := kvexpress.Set(c, KeyStop, KeyStopReason, Direction)

	if saved {
		log.Print(Direction, ": KeyStop='", KeyStop, "' saved='true' KeyStopReason='", KeyStopReason, "'")
		if DatadogAPIKey != "" && DatadogAPPKey != "" {
			kvexpress.DDSaveStopEvent(dog, KeyStop, KeyStopReason, Direction)
		}
	}

	// Run this command after the key is stopped.
	if PostExec != "" {
		log.Print(Direction, ": exec='", PostExec, "'")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, "complete", Direction)
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
	KeyStopReason   string
)

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVarP(&KeyStopLocation, "key", "k", "", "key to stop")
	stopCmd.Flags().StringVarP(&KeyStopReason, "reason", "r", "", "reason to stop")
}

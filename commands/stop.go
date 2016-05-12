// +build linux darwin freebsd

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zorkian/go-datadog-api"
	"os"
	"time"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Put stop value into Consul.",
	Long:  `Stop is a convenient way to put stop values in Consul.  Stops ALL nodes from updating.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkStopFlags()
		AutoEnable()
	},
	Run: stopRun,
}

func stopRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	var dog = new(datadog.Client)

	KeyStop := KeyPath(KeyStopLocation, "stop")

	c, err := Connect(ConsulServer, Token)
	if err != nil {
		LogFatal("Could not connect to Consul.", KeyStopLocation, "consul_connect")
	}

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
	}

	saved := Set(c, KeyStop, KeyStopReason)

	if saved {
		Log(fmt.Sprintf("KeyStop='%s' saved='true' KeyStopReason='%s'", KeyStop, KeyStopReason), "info")
		if DatadogAPIKey != "" && DatadogAPPKey != "" {
			DDSaveStopEvent(dog, KeyStop, KeyStopReason)
		}
	}

	// Run this command after the key is stopped.
	if PostExec != "" {
		Log(fmt.Sprintf("exec='%s'", PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, KeyStopLocation, "complete")
}

func checkStopFlags() {
	Log("Checking cli flags.", "debug")
	if KeyStopLocation == "" {
		fmt.Println("Need a key to stop in -k")
		os.Exit(1)
	}
	if KeyStopReason == "" {
		fmt.Println("Need a reason to stop in -r")
		os.Exit(1)
	}
	Log("Required cli flags present.", "debug")
}

var (
	// KeyStopLocation This Consul key is the one we want to halt all updates and distribution of.
	KeyStopLocation string

	// KeyStopReason This is the reason that we are stopping distribution of the key.
	KeyStopReason string
)

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVarP(&KeyStopLocation, "key", "k", "", "key to stop")
	stopCmd.Flags().StringVarP(&KeyStopReason, "reason", "r", "", "reason to stop")
}

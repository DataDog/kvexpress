package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zorkian/go-datadog-api"
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
	if EnvVars {
		ConfigEnvVars(Direction)
	}
	checkStopFlags(Direction)

	KeyStop := kvexpress.KeyStopPath(KeyStopLocation, PrefixLocation, Direction)

	c, _ := kvexpress.Connect(ConsulServer, Token, Direction)

	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		dog = kvexpress.DDAPIConnect(DatadogAPIKey, DatadogAPPKey)
	}

	saved := kvexpress.Set(c, KeyStop, KeyStopReason, Direction)

	if saved {
		kvexpress.Log(fmt.Sprintf("%s: KeyStop='%s' saved='true' KeyStopReason='%s'", Direction, KeyStop, KeyStopReason), "info")
		if DatadogAPIKey != "" && DatadogAPPKey != "" {
			kvexpress.DDSaveStopEvent(dog, KeyStop, KeyStopReason, Direction)
		}
	}

	// Run this command after the key is stopped.
	if PostExec != "" {
		kvexpress.Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, "complete", Direction)
}

func checkStopFlags(direction string) {
	kvexpress.Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if KeyStopLocation == "" {
		fmt.Println("Need a key to stop in -k")
		os.Exit(1)
	}
	if KeyStopReason == "" {
		fmt.Println("Need a reason to stop in -r")
		os.Exit(1)
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		kvexpress.Log(fmt.Sprintf("%s: Enabling Datadog API.", direction), "debug")
	}
	kvexpress.Log(fmt.Sprintf("%s: Required cli flags present.", direction), "debug")
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

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "kvexpress",
	Short: "Configuration data -> Consul KV -> Filesytem",
	Long:  `Small Go program to put and pull configuration data out of Consul and write to filesystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("`kvexpress -h` for help information.")
		fmt.Println("`kvexpress -v` for version information.")
	},
}

var Token string
var PostExec string
var ConsulServer string
var PrefixLocation string
var MinFileLength int
var FilePermissions int
var DogStatsd bool
var DogStatsdAddress string
var DatadogAPIKey string
var DatadogAPPKey string
var DatadogHost string

func init() {
	RootCmd.PersistentFlags().StringVarP(&ConsulServer, "server", "s", "localhost:8500", "Consul server location")
	RootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", "Token for Consul access")
	RootCmd.PersistentFlags().StringVarP(&PrefixLocation, "prefix", "p", "kvexpress", "prefix for the key")
	RootCmd.PersistentFlags().StringVarP(&PostExec, "exec", "e", "", "Execute this command after")
	RootCmd.PersistentFlags().IntVarP(&MinFileLength, "length", "l", 10, "minimum amount of lines in the file")
	RootCmd.PersistentFlags().IntVarP(&FilePermissions, "chmod", "c", 0640, "permissions for the file")
	RootCmd.PersistentFlags().BoolVarP(&DogStatsd, "dogstatsd", "d", false, "send metrics to dogstatsd")
	RootCmd.PersistentFlags().StringVarP(&DogStatsdAddress, "dogstatsd_addr", "D", "localhost:8125", "address for dogstatsd server")
	RootCmd.PersistentFlags().StringVarP(&DatadogAPIKey, "datadog_api", "a", "", "Datadog API Key")
	RootCmd.PersistentFlags().StringVarP(&DatadogAPPKey, "datadog_app", "A", "", "Datadog App Key")
	RootCmd.PersistentFlags().StringVarP(&DatadogHost, "datadog_host", "", "https://app.datadoghq.com", "Datadog Host")
}

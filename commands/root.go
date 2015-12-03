package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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

var (
	Token            string
	PostExec         string
	ConsulServer     string
	PrefixLocation   string
	MinFileLength    int
	FilePermissions  int
	DogStatsd        bool
	Owner            string
	ConfigFile       string
	DogStatsdAddress string
	DatadogAPIKey    string
	DatadogAPPKey    string
	Compress         bool
	Direction        string
)

func init() {
	Direction = os.Args[1]
	RootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "C", "", "Config file location")
	RootCmd.PersistentFlags().StringVarP(&ConsulServer, "server", "s", "localhost:8500", "Consul server location")
	RootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", "Token for Consul access")
	RootCmd.PersistentFlags().StringVarP(&PrefixLocation, "prefix", "p", "kvexpress", "prefix for the key")
	RootCmd.PersistentFlags().StringVarP(&PostExec, "exec", "e", "", "Execute this command after")
	RootCmd.PersistentFlags().IntVarP(&MinFileLength, "length", "l", 10, "minimum amount of lines in the file")
	RootCmd.PersistentFlags().IntVarP(&FilePermissions, "chmod", "c", 0640, "permissions for the file")
	RootCmd.PersistentFlags().BoolVarP(&DogStatsd, "dogstatsd", "d", false, "send metrics to dogstatsd")
	RootCmd.PersistentFlags().BoolVarP(&Compress, "compress", "z", false, "gzip in and out of the KV store")
	RootCmd.PersistentFlags().StringVarP(&DogStatsdAddress, "dogstatsd_address", "D", "localhost:8125", "address for dogstatsd server")
	RootCmd.PersistentFlags().StringVarP(&DatadogAPIKey, "datadog_api_key", "a", "", "Datadog API Key")
	RootCmd.PersistentFlags().StringVarP(&DatadogAPPKey, "datadog_app_key", "A", "", "Datadog App Key")
	RootCmd.PersistentFlags().StringVarP(&Owner, "owner", "o", "", "who to write the file as")
}

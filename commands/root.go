package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// RootCmd is the default Cobra struct that starts it all off.
// https://github.com/spf13/cobra
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
	// Token is used by the Consul API and references an ACL that Consul uses
	// to control access the KV store: https://www.consul.io/docs/internals/acl.html
	Token string

	// PostExec refers to an optional command to run upon
	// successful completion of the command's task. An example:
	// kvexpress out -k hosts -f /etc/hosts -e "sudo pkill -HUP dnsmasq"
	PostExec string

	// ConsulServer if you are not talking to a Consul node on localhost - this is for you.
	ConsulServer string

	// PrefixLocation all Consul KV data related to kvexpress is stored underneath
	// this path. Defaults to `kvexpress` which
	PrefixLocation string

	// MinFileLength is the minimum number of lines a file is expected to have.
	// Keeps blank or truncated files out of the KV store.
	MinFileLength int

	// FilePermissions are the permissions for the files that are written to the filesystem.
	FilePermissions int

	// DogStatsd enables reporting of tagged statsd metrics to the local Datadog Agent.
	// http://docs.datadoghq.com/guides/dogstatsd/
	DogStatsd bool

	// Owner will be the owner of any file that's been written to the filesystem.
	Owner string

	// ConfigFile is the path to a yaml encoded configuration file.
	// Loaded with LoadConfig.
	ConfigFile string

	// DogStatsdAddress if you're not running a local Datadog agent.
	DogStatsdAddress string

	// DatadogAPIKey is for sending events to Datadog through the HTTP api.
	DatadogAPIKey string

	// DatadogAPPKey is for sending events to Datadog through the HTTP api.
	DatadogAPPKey string

	// Compress is for compressing data on the way in and out of Consul.
	Compress bool

	// Direction adds information about which command is running to the logs.
	Direction string

	// Verbose logs all output to stdout.
	Verbose bool
)

func init() {
	// Do some setup.
	Direction = SetDirection()
	AutoEnable()
	RootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "C", "", "Config file location")
	RootCmd.PersistentFlags().StringVarP(&ConsulServer, "server", "s", "localhost:8500", "Consul server location")
	RootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "anonymous", "Token for Consul access")
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
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "", false, "log output to stdout")
}

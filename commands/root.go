package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "kvexpress",
	Short: "Consul KV > Filesytem",
	Long:  `Small Go program to pull data out of Consul and write to filesystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("`kvexpress -h` for help information.")
		fmt.Println("`kvexpress -v` ver version information.")
	},
}

var Token string
var PostExec string
var ConsulServer string
var PrefixLocation string

func init() {
	RootCmd.PersistentFlags().StringVarP(&ConsulServer, "server", "s", "localhost:8500", "Consul server location")
	RootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", "Token for Consul access")
	RootCmd.PersistentFlags().StringVarP(&PrefixLocation, "prefix", "p", "kvexpress", "prefix for the key")
	RootCmd.PersistentFlags().StringVarP(&PostExec, "exec", "e", "", "Execute this command after")
}

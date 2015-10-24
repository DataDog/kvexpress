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

func init() {
	RootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", "Token for Consul access")
}

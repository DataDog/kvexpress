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
		fmt.Println("kvexpress: -h for help information.")
	},
}

func init() {
}

// +build linux darwin freebsd

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch a prefix for changed keys.",
	Long:  `Watch a prefix for changed keys. Convert incoming keys to kvexpress format and location.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkWatchFlags()
		AutoEnable()
	},
	Run: watchRun,
}

func watchRun(cmd *cobra.Command, args []string) {
	// Check and aquire lock.
	// If you've got a lock - set a watch on WatchPrefix.
	// If you don't get a lock - then wait until you do.

	// If the watch fires - get the changed key and write into the /kvexpress heirarchy.
}

func checkWatchFlags() {
	Log("Checking cli flags.", "debug")
	if WatchPrefix == "" {
		fmt.Println("Need a KV space to watch in -w")
		os.Exit(1)
	}
	Log("Required cli flags present.", "debug")
}

var (
	// WatchPrefix is the Consul KV space to watch.
	WatchPrefix string
)

func init() {
	RootCmd.AddCommand(watchCmd)
	watchCmd.Flags().StringVarP(&WatchPrefix, "watch", "w", "", "Consul KV space to watch.")
}

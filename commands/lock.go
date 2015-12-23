package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock a file on a single node so it stays the way it is.",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkLockFlags()
	},
	Long: `Lock is a convenient way to stop a file from being updated on a single node.`,
	Run:  lockRun,
}

func lockRun(cmd *cobra.Command, args []string) {
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}

	KeyLockLocation := FileLockPath(FiletoLock)

	result := LockFile(KeyLockLocation)
	if result {
		LockFileWrite(FiletoLock)
		Log(fmt.Sprintf("'%s' was locked.", FiletoLock), "info")
	} else {
		Log(fmt.Sprintf("'%s' was NOT locked - something went wrong.", FiletoLock), "info")
	}
}

func checkLockFlags() {
	Log("Checking cli flags.", "debug")
	if FiletoLock == "" {
		fmt.Println("Need a file to lock with -f")
		os.Exit(1)
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		Log("Enabling Datadog API.", "debug")
	}
	if Owner == "" {
		Owner = GetCurrentUsername()
	}
	CheckFullFilename(FiletoLock)
	Log("Required cli flags present.", "debug")
}

var (
	// FiletoLock is the location we want to write the data to.
	FiletoLock string

	// LockReason is the reason why you are locking the file.
	LockReason string
)

func init() {
	RootCmd.AddCommand(lockCmd)
	lockCmd.Flags().StringVarP(&FiletoLock, "file", "f", "", "file to lock")
	lockCmd.Flags().StringVarP(&LockReason, "reason", "r", "", "reason to lock")
}

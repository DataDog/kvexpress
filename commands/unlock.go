package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unock a file on a single node so it updates.",
	Long:  `Unlock is a convenient way to allow a previously locked file to be updated.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkUnlockFlags()
	},
	Run: unlockRun,
}

func unlockRun(cmd *cobra.Command, args []string) {
	KeyLockLocation := FileLockPath(FiletoUnlock)

	result := UnlockFile(KeyLockLocation)
	if result {
		LockFileRemove(FiletoUnlock)
		Log(fmt.Sprintf("'%s' was unlocked.", FiletoUnlock), "info")
	} else {
		Log(fmt.Sprintf("'%s' was NOT unlocked - something went wrong.", FiletoUnlock), "info")
	}
}

func checkUnlockFlags() {
	Log("Checking cli flags.", "debug")
	if FiletoUnlock == "" {
		fmt.Println("Need a file to lock with -f")
		os.Exit(1)
	}
	CheckFullFilename(FiletoUnlock)
	Log("Required cli flags present.", "debug")
}

var (
	// FiletoUnlock is the location we want to write the data to.
	FiletoUnlock string
)

func init() {
	RootCmd.AddCommand(unlockCmd)
	unlockCmd.Flags().StringVarP(&FiletoUnlock, "file", "f", "", "file to unlock")
}

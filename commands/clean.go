package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean local cache files.",
	Long:  `clean is for cleaning up local cache files.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkCleanFlags()
	},
	Run: cleanRun,
}

func cleanRun(cmd *cobra.Command, args []string) {
	start := time.Now()

	CompareFile := CompareFilename(FiletoClean)
	LastFile := LastFilename(FiletoClean)

	RemoveFile(FiletoClean)
	RemoveFile(CompareFile)
	RemoveFile(LastFile)

	// Run this command after the files are cleaned.
	if PostExec != "" {
		Log(fmt.Sprintf("exec='%s'", PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, "none", "complete")
}

func checkCleanFlags() {
	Log("Checking cli flags.", "debug")
	if FiletoClean == "" {
		fmt.Println("Need a file to clean in -f")
		os.Exit(1)
	}
	Log("Required cli flags present.", "debug")
}

var (
	// FiletoClean is the file we want to erase.
	FiletoClean string
)

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVarP(&FiletoClean, "file", "f", "", "file to clean")
}

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
	Run:   cleanRun,
}

func cleanRun(cmd *cobra.Command, args []string) {
	start := time.Now()
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	checkCleanFlags(Direction)

	CompareFile := CompareFilename(FiletoClean, Direction)
	LastFile := LastFilename(FiletoClean, Direction)

	RemoveFile(FiletoClean, Direction)
	RemoveFile(CompareFile, Direction)
	RemoveFile(LastFile, Direction)

	// Run this command after the files are cleaned.
	if PostExec != "" {
		Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		RunCommand(PostExec)
	}
	RunTime(start, "none", "complete", Direction, DogStatsd)
}

func checkCleanFlags(direction string) {
	Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if FiletoClean == "" {
		fmt.Println("Need a file to clean in -f")
		os.Exit(1)
	}
	Log(fmt.Sprintf("%s: Required cli flags present.", direction), "debug")
}

var (
	FiletoClean string
)

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVarP(&FiletoClean, "file", "f", "", "file to clean")
}

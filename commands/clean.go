package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/spf13/cobra"
	"log"
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
	var Direction = "clean"
	checkCleanFlags(Direction)
	if EnvVars {
		ConfigEnvVars(Direction)
	}

	CompareFile := kvexpress.CompareFilename(FiletoClean, Direction)
	LastFile := kvexpress.LastFilename(FiletoClean, Direction)

	kvexpress.RemoveFile(FiletoClean, Direction)
	kvexpress.RemoveFile(CompareFile, Direction)
	kvexpress.RemoveFile(LastFile, Direction)

	// Run this command after the files are cleaned.
	if PostExec != "" {
		log.Print(Direction, ": exec='", PostExec, "'")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, "complete", Direction)
}

func checkCleanFlags(direction string) {
	log.Print(direction, ": Checking cli flags.")
	if FiletoClean == "" {
		fmt.Println("Need a file to clean in -f")
		os.Exit(1)
	}
	log.Print(direction, ": Required cli flags present.")
}

var (
	FiletoClean string
)

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVarP(&FiletoClean, "file", "f", "", "file to clean")
}

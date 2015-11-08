package commands

import (
	kvexpress "../kvexpress/"
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
	var Direction = "clean"
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	if EnvVars {
		ConfigEnvVars(Direction)
	}
	checkCleanFlags(Direction)

	CompareFile := kvexpress.CompareFilename(FiletoClean, Direction)
	LastFile := kvexpress.LastFilename(FiletoClean, Direction)

	kvexpress.RemoveFile(FiletoClean, Direction)
	kvexpress.RemoveFile(CompareFile, Direction)
	kvexpress.RemoveFile(LastFile, Direction)

	// Run this command after the files are cleaned.
	if PostExec != "" {
		kvexpress.Log(fmt.Sprintf("%s: exec='%s'", Direction, PostExec), "debug")
		kvexpress.RunCommand(PostExec)
	}
	kvexpress.RunTime(start, "none", "complete", Direction, DogStatsd)
}

func checkCleanFlags(direction string) {
	kvexpress.Log(fmt.Sprintf("%s: Checking cli flags.", direction), "debug")
	if FiletoClean == "" {
		fmt.Println("Need a file to clean in -f")
		os.Exit(1)
	}
	kvexpress.Log(fmt.Sprintf("%s: Required cli flags present.", direction), "debug")
}

var (
	FiletoClean string
)

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVarP(&FiletoClean, "file", "f", "", "file to clean")
}

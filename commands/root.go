package commands

import (
        "github.com/spf13/cobra"
        "fmt"
)

var RootCmd = &cobra.Command{
    Use:   "kvexpress",
    Short: "Consul KV > Filesytem",
    Long: `Small Go program to pull data out of Consul and write to filesystem.`,
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("kvexpress: -h for help information.")
    },
}

var KeyLocation string
var FiletoWrite string
var MinFileLength int

func init() {
    RootCmd.PersistentFlags().StringVar(&KeyLocation, "key", "", "key to pull file from")
    RootCmd.PersistentFlags().StringVar(&FiletoWrite, "file", "", "where to write the file")
    RootCmd.PersistentFlags().IntVar(&MinFileLength, "length", 10, "minimum amount of lines in the file")
}

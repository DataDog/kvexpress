package commands

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var outCmd = &cobra.Command{
	Use:   "out",
	Short: "Write a file based on key data.",
	Long:  `out is for writing a file based on a Consul key.`,
	Run:   outRun,
}

func outRun(cmd *cobra.Command, args []string) {
	checkFlags()

	// Get the data out of Consul.
	config := consulapi.DefaultConfig()
	config.Address = ConsulServer
	consul, err := consulapi.NewClient(config)
	kv := consul.KV()
	pair, _, err := kv.Get(KeyLocation, nil)
	if err != nil {
		panic(err)
	} else {
		log.Print("out: value='", string(pair.Value[:]), "'")
	}

	// Is the data long enough?

	// If the data is long enough, write the file.
}

func checkFlags() {
	log.Print("out: Checking cli flags.")
	if KeyLocation == "" {
		fmt.Println("Need a key location in -k")
		os.Exit(1)
	}
	if FiletoWrite == "" {
		fmt.Println("Need a file to write in -f")
		os.Exit(1)
	}
	log.Print("out: Required cli flags present.")
}

var KeyLocation string
var FiletoWrite string
var ConsulServer string
var MinFileLength int

func init() {
	RootCmd.AddCommand(outCmd)
	outCmd.Flags().StringVarP(&KeyLocation, "key", "k", "", "key to pull data from")
	outCmd.Flags().StringVarP(&FiletoWrite, "file", "f", "", "where to write the data")
	outCmd.Flags().StringVarP(&ConsulServer, "server", "s", "localhost:8500", "Consul server location")
	outCmd.Flags().IntVarP(&MinFileLength, "length", "l", 10, "minimum amount of lines in the file")
}

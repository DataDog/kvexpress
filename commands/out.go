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

	key_data = KeyDataPath(KeyLocation)
	key_checksum = KeyChecksumPath(KeyLocation)

	// Get the KV data out of Consul.
	KVData = get(key_data)

	// Get the Checksum data out of Consul.
	Checksum = get(key_checksum)

	// Is the data long enough?

	// Does the checksum match?

	// If the data is long enough and the checksum matches, write the file.
}

func KeyDataPath(key string) string {
	full_path := fmt.Sprint(PrefixLocation, "/", key, "/data")
	log.Print("out: full_path='", full_path, "'")
	return full_path
}

func KeyChecksumPath(key string) string {
	full_path := fmt.Sprint(PrefixLocation, "/", key, "/checksum")
	log.Print("out: full_path='", full_path, "'")
	return full_path
}

func get(key string) string {
	var value string
	config := consulapi.DefaultConfig()
	config.Address = ConsulServer
	consul, err := consulapi.NewClient(config)
	kv := consul.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		panic(err)
	} else {
		value = string(pair.Value[:])
		log.Print("out: value='", value, "'")
	}
	return value
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

var KVData string
var KeyLocation string
var Checksum string
var key_data string
var key_checksum string
var FiletoWrite string
var PrefixLocation string
var ConsulServer string
var MinFileLength int

func init() {
	RootCmd.AddCommand(outCmd)
	outCmd.Flags().StringVarP(&PrefixLocation, "prefix", "p", "kvexpress", "prefix for the key")
	outCmd.Flags().StringVarP(&KeyLocation, "key", "k", "", "key to pull data from")
	outCmd.Flags().StringVarP(&FiletoWrite, "file", "f", "", "where to write the data")
	outCmd.Flags().StringVarP(&ConsulServer, "server", "s", "localhost:8500", "Consul server location")
	outCmd.Flags().IntVarP(&MinFileLength, "length", "l", 10, "minimum amount of lines in the file")
}

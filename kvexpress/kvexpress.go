package kvexpress

import (
	"bytes"
	"fmt"
  consulapi "github.com/hashicorp/consul/api"
	"strings"
	"crypto/sha256"
	"log"
	"os"
	"io/ioutil"
	"os/exec"
)

func init()  {
  // Nothing happens here.
}

func KeyDataPath(key string, prefix string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/data")
	log.Print("out: path='data' full_path='", full_path, "'")
	return full_path
}

func KeyChecksumPath(key string, prefix string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/checksum")
	log.Print("out: path='checksum' full_path='", full_path, "'")
	return full_path
}

func Get(key string, server string, token string) string {
	var value string
	config := consulapi.DefaultConfig()
	config.Address = server
	if token != "" {
		config.Token = token
	}
	consul, err := consulapi.NewClient(config)
	kv := consul.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		panic(err)
	} else {
		value = string(pair.Value[:])
		log.Print("out: key='", key, "' value='", value, "' address='", server, "' token='", token, "'")
	}
	return value
}

func LengthCheck(data string, min_length int) bool {
	var length int
	if strings.ContainsAny(data, "\n") {
		length = strings.Count(data, "\n")
	} else {
		length = 1
	}
	log.Print("out: length='", length, "' min_length='", min_length, "'")
	if length >= min_length {
		return true
	} else {
		return false
	}
}

func ComputeChecksum(data string) string {
	data_bytes := []byte(data)
	computed_checksum := sha256.Sum256(data_bytes)
	final_checksum := fmt.Sprintf("%x\n", computed_checksum)
	log.Print("out: computed_checksum='", final_checksum, "'")
	return final_checksum
}

func ChecksumCompare(data string, checksum string) bool {
	computed_checksum := ComputeChecksum(data)
	log.Print("out: checksum='", checksum, "' computed_checksum='", computed_checksum, "'")
	if strings.TrimSpace(computed_checksum) == strings.TrimSpace(checksum) {
		return true
	} else {
		return false
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteFile(data string, filepath string, perms int) {
	err := ioutil.WriteFile(filepath, []byte(data), os.FileMode(perms))
	check(err)
	log.Print("out: file_wrote='true' location='", filepath, "' permissions='", perms, "'")
}

func RunCommand(command string) bool {
	parts := strings.Fields(command)
	cli := parts[0]
	args := parts[1:len(parts)]
	cmd := exec.Command(cli, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Print(err)
		return false
	} else {
		return true
	}
}

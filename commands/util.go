package commands

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// ReturnCurrentUTC returns the current UTC time in RFC3339 format.
func ReturnCurrentUTC() string {
	t := time.Now().UTC()
	dateUpdated := (t.Format(time.RFC3339))
	return dateUpdated
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// RunTime sends time coded logs and dogstatsd metrics when called.
// Location is set when the RunTime function is called.
func RunTime(start time.Time, key string, location string, dogstatsd bool) {
	elapsed := time.Since(start)
	if dogstatsd {
		milliseconds := int64(elapsed / time.Millisecond)
		StatsdRunTime(key, location, milliseconds)
	}
	Log(fmt.Sprintf("location='%s', elapsed='%s'", location, elapsed), "info")
}

// Log adds the global Direction to a message and sends to syslog.
// Syslog is setup in main.go
func Log(message, priority string) {
	message = fmt.Sprintf("%s: %s", Direction, message)
	if Verbose {
		time := ReturnCurrentUTC()
		fmt.Printf("%s: %s\n", time, message)
	}
	switch {
	case priority == "debug":
		if os.Getenv("KVEXPRESS_DEBUG") != "" {
			log.Print(message)
		}
	default:
		log.Print(message)
	}
}

// GetCurrentUsername grabs the current user running the kvexpress binary.
func GetCurrentUsername() string {
	usr, _ := user.Current()
	username := usr.Username
	Log(fmt.Sprintf("username='%s'", username), "debug")
	return username
}

// GetOwnerID looks up the User Id for the owner passed.
func GetOwnerID(owner string) int {
	var uid = ""
	var status = ""
	usr, err := user.Lookup(owner)
	if err != nil {
		usr, _ = user.Current()
		uid = usr.Uid
		status = "not_found"
	} else {
		uid = usr.Uid
		status = "found"
	}
	Log(fmt.Sprintf("owner='%s' status='%s' uid='%s'", owner, status, uid), "debug")
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	return int(uidInt)
}

// GetGroupID looks up the Group Id for the owner passed.
func GetGroupID(owner string) int {
	var gid = ""
	var status = ""
	usr, err := user.Lookup(owner)
	if err != nil {
		usr, _ = user.Current()
		gid = usr.Gid
		status = "not_found"
	} else {
		gid = usr.Gid
		status = "found"
	}
	Log(fmt.Sprintf("owner='%s' status='%s' gid='%s'", owner, status, gid), "debug")
	gidInt, err := strconv.ParseInt(gid, 10, 64)
	return int(gidInt)
}

// CompressData compresses and base64 encodes a string to place into Consul's KV store.
func CompressData(data string) string {
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	gz.Write([]byte(data))
	gz.Flush()
	gz.Close()
	encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())
	Log(fmt.Sprintf("compressing='true' full_size='%d' compressed_size='%d'", len(data), len(encoded)), "info")
	return encoded
}

// DecompressData base64 decodes and decompresses a string taken from Consul's KV store.
func DecompressData(data string) string {
	// If it's been compressed, it's been base64 encoded.
	raw, _ := base64.StdEncoding.DecodeString(data)
	// gunzip the string.
	unzipped, _ := gzip.NewReader(strings.NewReader(string(raw)))
	uncompressed, _ := ioutil.ReadAll(unzipped)
	Log(fmt.Sprintf("decompressing='true' size='%d'", len(uncompressed)), "info")
	return string(uncompressed)
}

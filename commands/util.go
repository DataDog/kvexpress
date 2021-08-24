// +build linux darwin freebsd

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

// SetDirection returns the direction.
func SetDirection() string {
	args := fmt.Sprintf("%x", os.Args)
	direction := "main"
	if strings.ContainsAny(args, " ") {
		if strings.HasPrefix(os.Args[1], "-") {
			direction = "main"
		} else {
			direction = os.Args[1]
		}
	}
	return direction
}

// RunTime sends time coded logs and dogstatsd metrics when called.
// Location is set when the RunTime function is called.
func RunTime(start time.Time, key string, location string) {
	elapsed := time.Since(start)
	milliseconds := int64(elapsed / time.Millisecond)
	StatsdRunTime(key, location, milliseconds)
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

// LogFatal prints to screen, sends to syslog, creates a fatal error
// and stops
func LogFatal(message string, id string, location string) {
	fullMessage := fmt.Sprintf("%s id:%s location:%s\n", message, id, location)
	Log(fullMessage, "info")
	fmt.Printf(fullMessage)
	StatsdPanic(id, location)
	// StatsdPanic exists with os.Exit(0)
}

// GetCurrentUsername grabs the current user running the kvexpress binary.
func GetCurrentUsername() string {
	usr, err := user.Current()
	if err != nil {
		Log("GetCurrentUsername(): user.Current has failed.", "info")
		return ""
	}
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
		usr, err = user.Current()
		if err != nil {
			Log("GetOwnerID(): Both user.Lookup and user.Current have failed.", "info")
		}
		uid = usr.Uid
		status = "not_found"
	} else {
		uid = usr.Uid
		status = "found"
	}
	Log(fmt.Sprintf("owner='%s' status='%s' uid='%s'", owner, status, uid), "debug")
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		Log("GetOwnerID(): Could not convert string into UID.", "info")
	}
	return int(uidInt)
}

// GetGroupID looks up the Group Id for the owner passed.
func GetGroupID(owner string) int {
	var gid = ""
	var status = ""
	usr, err := user.Lookup(owner)
	if err != nil {
		usr, err = user.Current()
		if err != nil {
			Log("GetGroupID(): Both user.Lookup and user.Current have failed.", "info")
		}
		gid = usr.Gid
		status = "not_found"
	} else {
		gid = usr.Gid
		status = "found"
	}
	Log(fmt.Sprintf("owner='%s' status='%s' gid='%s'", owner, status, gid), "debug")
	gidInt, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		Log("GetGroupID(): Could not convert string into GID.", "info")
	}
	return int(gidInt)
}

// CompressData compresses and base64 encodes a string to place into Consul's KV store.
func CompressData(data string) string {
	var compressed bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	gz.Write([]byte(data))
	gz.Flush()
	gz.Close()
	encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())
	Log(fmt.Sprintf("compressing='true' full_size='%d' compressed_size='%d'", len(data), len(encoded)), "info")
	return encoded
}

// DecompressData base64 decodes and decompresses a string taken from Consul's KV store.
func DecompressData(data string) string {
	if data != "" {
		// If it's been compressed, it's been base64 encoded.
		raw, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			Log("function='DecompressData' panic='true' method='base64.StdEncoding.DecodeString'", "info")
			fmt.Println("Panic: Could not base64 decode string.")
			StatsdPanic("key", "DecompressData")
		}
		// gunzip the string.
		unzipped, err := gzip.NewReader(strings.NewReader(string(raw)))
		if err != nil {
			Log("function='DecompressData' panic='true' method='gzip.NewReader'", "info")
			fmt.Println("Panic: Could not gunzip string.")
			StatsdPanic("key", "DecompressData")
		}
		uncompressed, err := ioutil.ReadAll(unzipped)
		if err != nil {
			Log("function='DecompressData' panic='true' method='ioutil.ReadAll'", "info")
			fmt.Println("Panic: Could not ioutil.ReadAll string.")
			StatsdPanic("key", "DecompressData")
		}
		Log(fmt.Sprintf("decompressing='true' size='%d'", len(uncompressed)), "info")
		return string(uncompressed)
	}
	return ""
}

// GetHostname returns the hostname.
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// AutoEnable helps to automatically enable flags based on cues from the environment.
func AutoEnable() {
	// Load the config file if passed.
	if ConfigFile != "" {
		LoadConfig(ConfigFile)
	}
	// Check for dd-agent configuration file.
	if _, err := os.Stat("/etc/dd-agent/datadog.conf"); err == nil {
		DogStatsd = true
	}
	if Owner == "" {
		Owner = GetCurrentUsername()
	}
	if DogStatsd {
		Log("Enabling Dogstatsd metrics.", "debug")
	}
	if DatadogAPIKey != "" && DatadogAPPKey != "" {
		Log("Enabling Datadog API.", "debug")
	}
}

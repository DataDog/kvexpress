// +build linux darwin freebsd

package commands

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func init() {
	// Nothing happens here.
}

// LengthCheck makes sure a string has at least minLength lines.
func LengthCheck(data string, minLength int) bool {
	length := LineCount(data)
	Log(fmt.Sprintf("length='%d' minLength='%d'", length, minLength), "debug")
	if length >= minLength {
		return true
	}
	return false
}

// ReadURL grabs a URL and returns the string from the body.
func ReadURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		Log(fmt.Sprintf("function='ReadURL' panic='true' url='%s'", url), "info")
		fmt.Printf("Panic: Could not open URL: '%s'\n", url)
		StatsdPanic(url, "read_url")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log(fmt.Sprintf("ReadURL(): Error reading '%s'", url), "info")
		return fmt.Sprintf("There was an error reading the body of the url: %s", url)
	}
	return string(body)
}

// LineCount splits a string by linebreak and returns the number of lines.
func LineCount(data string) int {
	var length int
	if strings.ContainsAny(data, "\n") {
		length = strings.Count(data, "\n")
	} else {
		length = 1
	}
	return length
}

// ComputeChecksum takes a string and computes a SHA256 checksum.
func ComputeChecksum(data string) string {
	dataBytes := []byte(data)
	computedChecksum := sha256.Sum256(dataBytes)
	finalChecksum := fmt.Sprintf("%x", computedChecksum)
	Log(fmt.Sprintf("computedChecksum='%s'", finalChecksum), "debug")
	if finalChecksum == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		Log("WARNING: That checksum means the data/key is blank. WARNING", "info")
	}
	return finalChecksum
}

// ChecksumCompare takes a string, generates a SHA256 checksum and compares
// against the passed checksum to see if they match.
func ChecksumCompare(data string, checksum string) bool {
	computedChecksum := ComputeChecksum(data)
	Log(fmt.Sprintf("checksum='%s' computedChecksum='%s'", checksum, computedChecksum), "debug")
	if strings.TrimSpace(computedChecksum) == strings.TrimSpace(checksum) {
		return true
	}
	return false
}

// UnixDiff runs diff to generate text for the Datadog events.
func UnixDiff(old, new string) string {
	diff, err := exec.Command("diff", "-u", old, new).Output()
	if err != nil {
		return "There was an error generating the diff."
	}
	text := string(diff)
	finalText := removeLines(text, 3)
	return finalText
}

// removeLines trims the top n number of lines from a string.
func removeLines(text string, number int) string {
	lines := strings.Split(text, "\n")
	var cleaned []string
	cleaned = append(cleaned, lines[number:]...)
	finalText := strings.Join(cleaned, "\n")
	return finalText
}

// RunCommand runs a cli command with arguments.
func RunCommand(command string) bool {
	parts := strings.Fields(command)
	cli := parts[0]
	args := parts[1:len(parts)]
	cmd := exec.Command(cli, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		Log(fmt.Sprintf("exec='error' message='%v'", err), "info")
		return false
	}
	return true
}

// GenerateLockReason creates a reason with filename, username and date.
func GenerateLockReason() string {
	reason := fmt.Sprintf("No reason given for '%s' by '%s' at '%s'.", FiletoLock, GetCurrentUsername(), ReturnCurrentUTC())
	return reason
}

// LockFile sets a key in Consul so that a particular file won't be updated. See commands/lock.go
func LockFile(key string) bool {
	c, err := Connect(ConsulServer, Token)
	if err != nil {
		LogFatal("Could not connect to Consul.", key, "consul_connect")
	}
	saved := Set(c, key, LockReason)
	if saved {
		StatsdLock(key)
		return true
	}
	return false
}

// UnlockFile removes a key in Consul so that a particular file can be updated. See commands/unlock.go
func UnlockFile(key string) bool {
	c, err := Connect(ConsulServer, Token)
	if err != nil {
		LogFatal("copy: Could not connect to Consul.", key, "consul_connect")
	}
	value := Del(c, key)
	StatsdUnlock(key)
	return value
}

package commands

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/aryann/difflib"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func init() {
	// Nothing happens here.
}

func LengthCheck(data string, min_length int) bool {
	length := LineCount(data)
	Log(fmt.Sprintf("length='%d' min_length='%d'", length, min_length), "debug")
	if length >= min_length {
		return true
	} else {
		return false
	}
}

func ReadUrl(url string, dogstatsd bool) string {
	resp, err := http.Get(url)
	if err != nil {
		Log(fmt.Sprintf("function='ReadUrl' panic='true' url='%s'", url), "info")
		if dogstatsd {
			StatsdPanic(url, "read_url")
		}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func LineCount(data string) int {
	var length int
	if strings.ContainsAny(data, "\n") {
		length = strings.Count(data, "\n")
	} else {
		length = 1
	}
	return length
}

func ComputeChecksum(data string) string {
	data_bytes := []byte(data)
	computed_checksum := sha256.Sum256(data_bytes)
	final_checksum := fmt.Sprintf("%x\n", computed_checksum)
	Log(fmt.Sprintf("computed_checksum='%s'", final_checksum), "debug")
	return final_checksum
}

func ChecksumCompare(data string, checksum string) bool {
	computed_checksum := ComputeChecksum(data)
	Log(fmt.Sprintf("checksum='%s' computed_checksum='%s'", checksum, computed_checksum), "debug")
	if strings.TrimSpace(computed_checksum) == strings.TrimSpace(checksum) {
		return true
	} else {
		return false
	}
}

func UnixDiff(old, new string) string {
	diff, _ := exec.Command("diff", "-u", old, new).Output()
	text := string(diff)
	finalText := removeLines(text, 3)
	return finalText
}

func removeLines(text string, number int) string {
	lines := strings.Split(text, "\n")
	cleaned := make([]string, 0)
	cleaned = append(cleaned, lines[number:]...)
	finalText := strings.Join(cleaned, "\n")
	return finalText
}

func Diff(last string, current string) string {
	var buffer bytes.Buffer

	// Split lines.
	last_strings := strings.Split(string(last), "\n")
	current_strings := strings.Split(string(current), "\n")

	diff := difflib.Diff(last_strings, current_strings)
	diffString := fmt.Sprintf("%v", diff)

	Log("doing the diff", "debug")
	buffer.WriteString(diffString)
	return buffer.String()
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
		Log(fmt.Sprintf("exec='error' message='%v'", err), "info")
		return false
	} else {
		return true
	}
}

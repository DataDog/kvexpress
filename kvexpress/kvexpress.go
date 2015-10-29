package kvexpress

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/aryann/difflib"
	"html"
	"log"
	"os/exec"
	"strings"
)

func init() {
	// Nothing happens here.
}

func LengthCheck(data string, min_length int, direction string) bool {
	length := LineCount(data)
	log.Print(direction, ": length='", length, "' min_length='", min_length, "'")
	if length >= min_length {
		return true
	} else {
		return false
	}
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

func ComputeChecksum(data string, direction string) string {
	data_bytes := []byte(data)
	computed_checksum := sha256.Sum256(data_bytes)
	final_checksum := fmt.Sprintf("%x\n", computed_checksum)
	log.Print(direction, ": computed_checksum='", final_checksum, "'")
	return final_checksum
}

func ChecksumCompare(data string, checksum string, direction string) bool {
	computed_checksum := ComputeChecksum(data, direction)
	log.Print(direction, ": checksum='", checksum, "' computed_checksum='", computed_checksum, "'")
	if strings.TrimSpace(computed_checksum) == strings.TrimSpace(checksum) {
		return true
	} else {
		return false
	}
}

func HTMLDiff(last string, current string) string {
	var buffer bytes.Buffer

	// Split lines.
	last_strings := strings.Split(html.EscapeString(string(last)), "\n")
	current_strings := strings.Split(html.EscapeString(string(current)), "\n")

	log.Print("in: Doing the diff.")
	buffer.WriteString("<table>")
	buffer.WriteString(difflib.HTMLDiff(last_strings, current_strings))
	buffer.WriteString("</table>")
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
		log.Print(err)
		return false
	} else {
		return true
	}
}

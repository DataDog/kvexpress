package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

// ReadFile reads a file in the filesystem and returns a string.
func ReadFile(filepath string) string {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		dat = []byte("")
	}
	return string(dat)
}

// SortFile takes a string, splits it into lines, removes all blank lines using
// BlankLineStrip() and then sorts the remaining lines.
func SortFile(file string) string {
	Log("sorting='true'", "debug")
	lines := strings.Split(file, "\n")
	lines = BlankLineStrip(lines)
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// BlankLineStrip takes a slice of strings, ranges over them and only returns
// a slice of strings where the lines weren't blank.
func BlankLineStrip(data []string) []string {
	Log(fmt.Sprintf("in: stripping_blank_lines='true'"), "debug")
	var stripped []string
	for _, str := range data {
		if str != "" {
			stripped = append(stripped, str)
		}
	}
	return stripped
}

// WriteFile writes a string to a filepath. It also chowns the file to the owner and group
// of the user running the program if it's not set as a different user.
func WriteFile(data string, filepath string, perms int, owner string, dogstatsd bool) {
	var fileChown = false
	err := ioutil.WriteFile(filepath, []byte(data), os.FileMode(perms))
	if err != nil {
		Log(fmt.Sprintf("function='WriteFile' panic='true' file='%s'", filepath), "info")
		if dogstatsd {
			StatsdPanic(filepath, "write_file")
		}
	}
	oid := GetOwnerID(owner)
	gid := GetGroupID(owner)
	err = os.Chown(filepath, oid, gid)
	if err != nil {
		fileChown = false
		if dogstatsd {
			StatsdPanic(filepath, "chown_file")
		}
	} else {
		fileChown = true
	}
	Log(fmt.Sprintf("file_wrote='true' location='%s' permissions='%s'", filepath, strconv.FormatInt(int64(perms), 8)), "debug")
	Log(fmt.Sprintf("file_chown='%t' location='%s' owner='%d' group='%d'", fileChown, filepath, oid, gid), "debug")
}

// CheckFiletoWrite takes a filename and checksum and stops execution if
// there is a directory OR the file has the same checksum.
func CheckFiletoWrite(filename, checksum string) {
	// Try to open the file.
	file, err := os.Open(filename)
	f, err := file.Stat()
	switch {
	case err != nil:
		Log(fmt.Sprintf("there is NO file at %s", filename), "debug")
		break
	case f.IsDir():
		Log(fmt.Sprintf("Can NOT write a directory %s", filename), "info")
		os.Exit(1)
	default:
		data, _ := ioutil.ReadFile(filename)
		computedChecksum := ComputeChecksum(string(data))
		if computedChecksum == checksum {
			Log(fmt.Sprintf("'%s' has the same checksum. Stopping.", filename), "info")
			os.Exit(0)
		}
	}
	// If there's no file - then great - there's nothing to check
}

// RemoveFile takes a filename and stops if it's a directory. It will log success
// or failure of removal.
func RemoveFile(filename string) {
	file, err := os.Open(filename)
	f, err := file.Stat()
	switch {
	case err != nil:
		Log(fmt.Sprintf("Could NOT stat %s", filename), "debug")
	case f.IsDir():
		Log(fmt.Sprintf("Would NOT remove a directory %s", filename), "info")
		os.Exit(1)
	default:
		err = os.Remove(filename)
		if err != nil {
			Log(fmt.Sprintf("Could NOT remove %s", filename), "info")
		} else {
			Log(fmt.Sprintf("Removed %s", filename), "info")
		}
	}
}

// RandomTmpFile is used to create a .compare or .last file for UrltoRead()
func RandomTmpFile() string {
	file, err := ioutil.TempFile(os.TempDir(), "kvexpress")
	if err != nil {
		Log("function='RandomTmpFile' panic='true'", "info")
	}
	fileName := file.Name()
	Log(fmt.Sprintf("tempfile='%s'", fileName), "debug")
	return fileName
}

// CompareFilename returns a .compare filename based on the passed file.
func CompareFilename(file string) string {
	compare := fmt.Sprintf("%s.compare", path.Base(file))
	fullPath := path.Join(path.Dir(file), compare)
	Log(fmt.Sprintf("file='compare' fullPath='%s'", fullPath), "debug")
	return fullPath
}

// LastFilename returns a .last filename based on the passed file.
func LastFilename(file string) string {
	last := fmt.Sprintf("%s.last", path.Base(file))
	fullPath := path.Join(path.Dir(file), last)
	Log(fmt.Sprintf("file='last' fullPath='%s'", fullPath), "debug")
	return fullPath
}

// CheckLastFile creates a .last file if it doesn't exist.
func CheckLastFile(file string, perms int, owner string, dogstatsd bool) {
	if _, err := os.Stat(file); err != nil {
		Log(fmt.Sprintf("file='last' file='%s' does_not_exist='true'", file), "debug")
		WriteFile("This is a blank file.\n", file, perms, owner, dogstatsd)
	}
}

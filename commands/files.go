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

func ReadFile(filepath string) string {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		dat = []byte("")
	}
	return string(dat)
}

func SortFile(file string) string {
	Log("in: sorting='true'", "debug")
	lines := strings.Split(file, "\n")
	lines = BlankLineStrip(lines)
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

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

func WriteFile(data string, filepath string, perms int, owner string, direction string, dogstatsd bool) {
	var fileChown = false
	err := ioutil.WriteFile(filepath, []byte(data), os.FileMode(perms))
	if err != nil {
		Log(fmt.Sprintf("%s: function='WriteFile' panic='true' file='%s'", direction, filepath), "info")
		if dogstatsd {
			StatsdPanic(filepath, direction, "write_file")
		}
	}
	oid := GetOwnerId(owner, direction)
	gid := GetGroupId(owner, direction)
	err = os.Chown(filepath, oid, gid)
	if err != nil {
		fileChown = false
		if dogstatsd {
			StatsdPanic(filepath, direction, "chown_file")
		}
	} else {
		fileChown = true
	}
	Log(fmt.Sprintf("%s: file_wrote='true' location='%s' permissions='%s'", direction, filepath, strconv.FormatInt(int64(perms), 8)), "debug")
	Log(fmt.Sprintf("%s: file_chown='%t' location='%s' owner='%d' group='%d'", direction, fileChown, filepath, oid, gid), "debug")
}

func CheckFiletoWrite(filename, checksum, direction string) {
	// Try to open the file.
	file, err := os.Open(filename)
	f, err := file.Stat()
	switch {
	case err != nil:
		Log(fmt.Sprintf("%s: there is NO file at %s", direction, filename), "debug")
		break
	case f.IsDir():
		Log(fmt.Sprintf("%s: Can NOT write a directory %s", direction, filename), "info")
		os.Exit(1)
	default:
		data, _ := ioutil.ReadFile(filename)
		computedChecksum := ComputeChecksum(string(data), direction)
		if computedChecksum == checksum {
			Log(fmt.Sprintf("%s: '%s' has the same checksum. Stopping.", direction, filename), "info")
			os.Exit(0)
		}
	}

	// If there's no file - then great - there's nothing to check
}

func RemoveFile(filename string, direction string) {
	file, err := os.Open(filename)
	f, err := file.Stat()
	switch {
	case err != nil:
		Log(fmt.Sprintf("%s: Could NOT stat %s", direction, filename), "debug")
	case f.IsDir():
		Log(fmt.Sprintf("%s: Would NOT remove a directory %s", direction, filename), "info")
		os.Exit(1)
	default:
		err = os.Remove(filename)
		if err != nil {
			Log(fmt.Sprintf("%s: Could NOT remove %s", direction, filename), "info")
		} else {
			Log(fmt.Sprintf("%s: Removed %s", direction, filename), "info")
		}
	}
}

func RandomTmpFile(direction string) string {
	file, err := ioutil.TempFile(os.TempDir(), "kvexpress")
	if err != nil {
		Log(fmt.Sprintf("%s: function='RandomTmpFile' panic='true'", direction), "info")
	}
	fileName := file.Name()
	Log(fmt.Sprintf("%s: tempfile='%s'", direction, fileName), "debug")
	return fileName
}

func CompareFilename(file string, direction string) string {
	compare := fmt.Sprintf("%s.compare", path.Base(file))
	full_path := path.Join(path.Dir(file), compare)
	Log(fmt.Sprintf("%s: file='compare' full_path='%s'", direction, full_path), "debug")
	return full_path
}

func LastFilename(file string, direction string) string {
	last := fmt.Sprintf("%s.last", path.Base(file))
	full_path := path.Join(path.Dir(file), last)
	Log(fmt.Sprintf("%s: file='last' full_path='%s'", direction, full_path), "debug")
	return full_path
}

func CheckLastFile(file string, perms int, owner string, dogstatsd bool) {
	if _, err := os.Stat(file); err != nil {
		Log(fmt.Sprintf("in: file='last' file='%s' does_not_exist='true'", file), "debug")
		WriteFile("This is a blank file.\n", file, perms, owner, "in", dogstatsd)
	}
}

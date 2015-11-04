package kvexpress

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
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
	lines := strings.Split(file, "\n")
	lines = BlankLineStrip(lines)
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func BlankLineStrip(data []string) []string {
	var stripped []string
	for _, str := range data {
		if str != "" {
			stripped = append(stripped, str)
		}
	}
	return stripped
}

func WriteFile(data string, filepath string, perms int, direction string) {
	err := ioutil.WriteFile(filepath, []byte(data), os.FileMode(perms))
	check(err)
	log.Print(direction, ": file_wrote='true' location='", filepath, "' permissions='", perms, "'")
}

func CheckFiletoWrite(filename, checksum, direction string) {
	// Try to open the file.
	file, err := os.Open(filename)
	f, err := file.Stat()
	switch {
	case err != nil:
		log.Print(direction, ": there is NO file at ", filename)
		break
	case f.IsDir():
		log.Print(direction, ": Can NOT write a directory ", filename)
		os.Exit(1)
	default:
		data, _ := ioutil.ReadFile(filename)
		computedChecksum := ComputeChecksum(string(data), direction)
		if computedChecksum == checksum {
			log.Print(direction, ": already a file with the same checksum. Stopping.")
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
		log.Print(direction, ": Could NOT stat ", filename)
	case f.IsDir():
		log.Print(direction, ": Would NOT remove a directory ", filename)
		os.Exit(1)
	default:
		err = os.Remove(filename)
		if err != nil {
			log.Print(direction, ": Could NOT remove ", filename)
		} else {
			log.Print(direction, ": Removed ", filename)
		}
	}
}

func CompareFilename(file string, direction string) string {
	compare := fmt.Sprintf("%s.compare", path.Base(file))
	full_path := path.Join(path.Dir(file), compare)
	log.Print(direction, ": file='compare' full_path='", full_path, "'")
	return full_path
}

func LastFilename(file string, direction string) string {
	last := fmt.Sprintf("%s.last", path.Base(file))
	full_path := path.Join(path.Dir(file), last)
	log.Print(direction, ": file='last' full_path='", full_path, "'")
	return full_path
}

func CheckLastFile(file string, perms int) {
	if _, err := os.Stat(file); err != nil {
		log.Print("in: Last File: ", file, " does not exist.")
		WriteFile("This is a blank file.\n", file, perms, "in")
	}
}

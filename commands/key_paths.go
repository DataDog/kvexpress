package commands

import (
	"fmt"
)

// KeyDataPath returns the key to fetch kvexpress data.
func KeyDataPath(key string, prefix string) string {
	fullPath := fmt.Sprint(prefix, "/", key, "/data")
	Log(fmt.Sprintf("path='data' fullPath='%s'", fullPath), "debug")
	return fullPath
}

// KeyChecksumPath returns the key to fetch kvexpress checksum.
func KeyChecksumPath(key string, prefix string) string {
	fullPath := fmt.Sprint(prefix, "/", key, "/checksum")
	Log(fmt.Sprintf("path='checksum' fullPath='%s'", fullPath), "debug")
	return fullPath
}

// KeyStopPath returns the key to fetch kvexpress stop information.
func KeyStopPath(key string, prefix string) string {
	fullPath := fmt.Sprint(prefix, "/", key, "/stop")
	Log(fmt.Sprintf("path='stop' fullPath='%s'", fullPath), "debug")
	return fullPath
}

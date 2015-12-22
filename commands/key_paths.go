package commands

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// KeyDataPath returns the key to fetch kvexpress data.
func KeyDataPath(key string) string {
	fullPath := fmt.Sprint(strings.TrimPrefix(PrefixLocation, "/"), "/", key, "/data")
	Log(fmt.Sprintf("path='data' fullPath='%s'", fullPath), "debug")
	return fullPath
}

// KeyChecksumPath returns the key to fetch kvexpress checksum.
func KeyChecksumPath(key string) string {
	fullPath := fmt.Sprint(strings.TrimPrefix(PrefixLocation, "/"), "/", key, "/checksum")
	Log(fmt.Sprintf("path='checksum' fullPath='%s'", fullPath), "debug")
	return fullPath
}

// KeyStopPath returns the key to fetch kvexpress stop information.
func KeyStopPath(key string) string {
	fullPath := fmt.Sprint(strings.TrimPrefix(PrefixLocation, "/"), "/", key, "/stop")
	Log(fmt.Sprintf("path='stop' fullPath='%s'", fullPath), "debug")
	return fullPath
}

// GenerateFileLockPath generates the path for the KV store for a particular file.
func FileLockPath(file string) string {
	hostname := GetHostname()
	fileSHA := sha256.Sum256([]byte(file))
	fileSHAs := fmt.Sprintf("%x", fileSHA)
	path := fmt.Sprintf("%s/locks/%s/%s", strings.TrimPrefix(PrefixLocation, "/"), fileSHAs, hostname)
	Log(fmt.Sprintf("path='%s'", path), "debug")
	return path
}

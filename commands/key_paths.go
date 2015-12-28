// +build linux darwin freebsd

package commands

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// KeyPath returns the standard kvexpress paths for data, checksum and stop.
func KeyPath(key string, suffix string) string {
	fullPath := fmt.Sprintf("%s/%s/%s", strings.TrimPrefix(PrefixLocation, "/"), key, suffix)
	Log(fmt.Sprintf("path='%s' fullPath='%s'", suffix, fullPath), "debug")
	return fullPath
}

// FileLockPath generates the path for the KV store for a particular file.
func FileLockPath(file string) string {
	hostname := GetHostname()
	fileSHA := sha256.Sum256([]byte(file))
	fileSHAs := fmt.Sprintf("%x", fileSHA)
	path := fmt.Sprintf("%s/locks/%s/%s", strings.TrimPrefix(PrefixLocation, "/"), fileSHAs, hostname)
	Log(fmt.Sprintf("path='%s'", path), "debug")
	return path
}

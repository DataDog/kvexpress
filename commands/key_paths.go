package commands

import (
	"fmt"
)

func KeyDataPath(key string, prefix string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/data")
	Log(fmt.Sprintf("path='data' full_path='%s'", full_path), "debug")
	return full_path
}

func KeyChecksumPath(key string, prefix string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/checksum")
	Log(fmt.Sprintf("path='checksum' full_path='%s'", full_path), "debug")
	return full_path
}

func KeyStopPath(key string, prefix string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/stop")
	Log(fmt.Sprintf("path='stop' full_path='%s'", full_path), "debug")
	return full_path
}

func KeyUpdatedPath(key string, prefix string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/updated")
	Log(fmt.Sprintf("path='updated' full_path='%s'", full_path), "debug")
	return full_path
}

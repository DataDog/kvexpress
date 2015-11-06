package kvexpress

import (
	"fmt"
)

func KeyDataPath(key string, prefix string, direction string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/data")
	Log(fmt.Sprintf("%s: path='data' full_path='%s'", direction, full_path), "debug")
	return full_path
}

func KeyChecksumPath(key string, prefix string, direction string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/checksum")
	Log(fmt.Sprintf("%s: path='checksum' full_path='%s'", direction, full_path), "debug")
	return full_path
}

func KeyStopPath(key string, prefix string, direction string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/stop")
	Log(fmt.Sprintf("%s: path='stop' full_path='%s'", direction, full_path), "debug")
	return full_path
}

func KeyUpdatedPath(key string, prefix string, direction string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/updated")
	Log(fmt.Sprintf("%s: path='updated' full_path='%s'", direction, full_path), "debug")
	return full_path
}

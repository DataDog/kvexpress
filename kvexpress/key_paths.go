package kvexpress

import (
	"fmt"
	"log"
)

func KeyDataPath(key string, prefix string, direction string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/data")
	log.Print(direction, ": path='data' full_path='", full_path, "'")
	return full_path
}

func KeyChecksumPath(key string, prefix string, direction string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/checksum")
	log.Print(direction, ": path='checksum' full_path='", full_path, "'")
	return full_path
}

func KeyStopPath(key string, prefix string, direction string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/stop")
	log.Print(direction, ": path='stop' full_path='", full_path, "'")
	return full_path
}

func KeyUpdatedPath(key string, prefix string, direction string) string {
	full_path := fmt.Sprint(prefix, "/", key, "/updated")
	log.Print(direction, ": path='updated' full_path='", full_path, "'")
	return full_path
}

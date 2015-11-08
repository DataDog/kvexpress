package kvexpress

import (
	"fmt"
	"log"
	"os"
	"time"
)

func ReturnCurrentUTC() string {
	t := time.Now().UTC()
	date_updated := (t.Format(time.RFC3339))
	return date_updated
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func RunTime(start time.Time, key string, location string, direction string, dogstatsd bool) {
	elapsed := time.Since(start)
	if dogstatsd {
		milliseconds := int64(elapsed / time.Millisecond)
		StatsdRunTime(direction, key, location, milliseconds)
	}
	Log(fmt.Sprintf("%s: location='%s', elapsed='%s'", direction, location, elapsed), "info")
}

func Log(message, priority string) {
	switch {
	case priority == "debug":
		if os.Getenv("KVEXPRESS_DEBUG") != "" {
			log.Print(message)
		}
	default:
		log.Print(message)
	}

}

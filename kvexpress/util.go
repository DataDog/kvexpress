package kvexpress

import (
	"fmt"
	"log"
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

func RunTime(start time.Time, location string, direction string) {
	elapsed := time.Since(start)
	log.Print(fmt.Sprintf("%s: location='%s', elapsed='%s'", direction, location, elapsed))
}

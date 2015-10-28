package kvexpress

import (
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

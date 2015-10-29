package kvexpress

import (
  "fmt"
	"github.com/PagerDuty/godspeed"
)

func StatsdIn(key string, data_length int, data string) {
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	statsdTags := []string{fmt.Sprintf("kvkey:%s", key)}
	statsd.Incr("kvexpress.in", statsdTags)
	statsd.Gauge("kvexpress.bytes", float64(data_length), statsdTags)
	statsd.Gauge("kvexpress.lines", float64(LineCount(data)), statsdTags)
}

package kvexpress

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"github.com/zorkian/go-datadog-api"
	"os"
)

func StatsdIn(key string, data_length int, data string) {
	Log(fmt.Sprintf("in: dogstatsd='true' key='%s' stats='in'", key), "debug")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(key, "in")
	statsd.Incr("kvexpress.in", tags)
	statsd.Gauge("kvexpress.bytes", float64(data_length), tags)
	statsd.Gauge("kvexpress.lines", float64(LineCount(data)), tags)
}

func StatsdOut(key string) {
	Log(fmt.Sprintf("out: dogstatsd='true' key='%s' stats='out'", key), "debug")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(key, "out")
	statsd.Incr("kvexpress.out", tags)
}

func StatsdRunTime(direction string, key string, location string, msec int64) {
	Log(fmt.Sprintf("%s: dogstatsd='true' key='%s' location='%s' msec='%d'", direction, key, location, msec), "debug")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(key, direction)
	locationTag := fmt.Sprintf("location:%s", location)
	tags = append(tags, locationTag)
	statsd.Gauge("kvexpress.time", float64(msec), tags)
}

func DDAPIConnect(api, app string) *datadog.Client {
	client := datadog.NewClient(api, app)
	return client
}

func makeTags(key, direction string) []string {
	tags := make([]string, 3)
	keyTag := fmt.Sprintf("key:%s", key)
	hostname, _ := os.Hostname()
	hostTag := fmt.Sprintf("host:%s", hostname)
	directionTag := fmt.Sprintf("direction:%s", direction)
	tags = append(tags, keyTag)
	tags = append(tags, hostTag)
	tags = append(tags, directionTag)
	return tags
}

// TODO: These three functions are ripe for refactoring to be more Golang like.
func DDStopEvent(dd *datadog.Client, key, value, direction string) {
	Log(fmt.Sprintf("%s: datadog='true' DDStopEvent='true' key='%s'", direction, key), "debug")
	tags := makeTags(key, direction)
	tags = append(tags, "kvexpress:stop")
	title := fmt.Sprintf("Stop key is present: %s. Stopping.", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "error", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

func DDSaveDataEvent(dd *datadog.Client, key, value, direction string) {
	Log(fmt.Sprintf("%s: datadog='true' DDSaveDataEvent='true' key='%s'", direction, key), "debug")
	tags := makeTags(key, direction)
	tags = append(tags, "kvexpress:success")
	title := fmt.Sprintf("Updated: %s", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "info", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

func DDSaveStopEvent(dd *datadog.Client, key, value, direction string) {
	Log(fmt.Sprintf("%s: datadog='true' DDSaveStopEvent='true' key='%s'", direction, key), "debug")
	tags := makeTags(key, direction)
	tags = append(tags, "kvexpress:stop_set")
	title := fmt.Sprintf("Set Stop Key: %s", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "warning", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

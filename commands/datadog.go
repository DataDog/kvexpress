package commands

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"github.com/zorkian/go-datadog-api"
	"os"
)

func StatsdIn(key string, dataLength int, data string) {
	Log(fmt.Sprintf("dogstatsd='true' key='%s' stats='in'", key), "debug")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(key, "complete")
	statsd.Incr("kvexpress.in", tags)
	statsd.Gauge("kvexpress.bytes", float64(dataLength), tags)
	statsd.Gauge("kvexpress.lines", float64(LineCount(data)), tags)
}

func StatsdOut(key string) {
	Log(fmt.Sprintf("dogstatsd='true' key='%s' stats='out'", key), "debug")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(key, "complete")
	statsd.Incr("kvexpress.out", tags)
}

func StatsdRaw(key string) {
	Log(fmt.Sprintf("dogstatsd='true' key='%s' stats='raw'", key), "debug")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(key, "complete")
	statsd.Incr("kvexpress.raw", tags)
}

func StatsdRunTime(key string, location string, msec int64) {
	Log(fmt.Sprintf("dogstatsd='true' key='%s' location='%s' msec='%d'", key, location, msec), "debug")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(key, location)
	locationTag := fmt.Sprintf("location:%s", location)
	tags = append(tags, locationTag)
	statsd.Gauge("kvexpress.time", float64(msec), tags)
}

func StatsdPanic(key, location string) {
	Log(fmt.Sprintf("dogstatsd='true' key='%s' location='%s' stats='panic'", key, location), "debug")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(key, location)
	statsd.Incr("kvexpress.panic", tags)
	// If we're going to panic, we might as well stop right here.
	// Means we can't connect to Consul, download a URL or
	// write and/or chown files.
	os.Exit(0)
}

func DDAPIConnect(api, app string) *datadog.Client {
	client := datadog.NewClient(api, app)
	return client
}

func makeTags(key, location string) []string {
	tags := make([]string, 4)
	keyTag := fmt.Sprintf("key:%s", key)
	hostname, _ := os.Hostname()
	hostTag := fmt.Sprintf("host:%s", hostname)
	directionTag := fmt.Sprintf("direction:%s", Direction)
	locationTag := fmt.Sprintf("location:%s", location)
	tags = append(tags, keyTag)
	tags = append(tags, hostTag)
	tags = append(tags, directionTag)
	tags = append(tags, locationTag)
	return tags
}

// TODO: These three functions are ripe for refactoring to be more Golang like.
func DDStopEvent(dd *datadog.Client, key, value string) {
	Log(fmt.Sprintf("datadog='true' DDStopEvent='true' key='%s'", key), "debug")
	tags := makeTags(key, "stop_key_present")
	tags = append(tags, "kvexpress:stop")
	title := fmt.Sprintf("Stop key is present: %s. Stopping.", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "error", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

func DDSaveDataEvent(dd *datadog.Client, key, value string) {
	Log(fmt.Sprintf("datadog='true' DDSaveDataEvent='true' key='%s'", key), "debug")
	tags := makeTags(key, "complete")
	tags = append(tags, "kvexpress:success")
	title := fmt.Sprintf("Updated: %s", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "info", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

func DDCopyDataEvent(dd *datadog.Client, keyFrom, keyTo string) {
	Log(fmt.Sprintf("datadog='true' DDCopyDataEvent='true' keyFrom='%s' keyTo='%s'", keyFrom, keyTo), "debug")
	tags := makeTags(keyTo, "complete")
	tags = append(tags, "kvexpress:success")
	tags = append(tags, fmt.Sprintf("keyFrom:%s", keyFrom))
	title := fmt.Sprintf("Copy: %s to %s", keyFrom, keyTo)
	event := datadog.Event{Title: title, Text: title, AlertType: "info", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

func DDSaveStopEvent(dd *datadog.Client, key, value string) {
	Log(fmt.Sprintf("datadog='true' DDSaveStopEvent='true' key='%s'", key), "debug")
	tags := makeTags(key, "stop_key_save")
	tags = append(tags, "kvexpress:stop_set")
	title := fmt.Sprintf("Set Stop Key: %s", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "warning", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

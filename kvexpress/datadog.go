package kvexpress

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"github.com/darron/go-datadog-api"
	"os"
)

func StatsdIn(key string, data_length int, data string) {
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	statsdTags := []string{fmt.Sprintf("kvkey:%s", key)}
	statsd.Incr("kvexpress.in", statsdTags)
	statsd.Gauge("kvexpress.bytes", float64(data_length), statsdTags)
	statsd.Gauge("kvexpress.lines", float64(LineCount(data)), statsdTags)
}

func StatsdOut(key string) {
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	statsdTags := []string{fmt.Sprintf("kvkey:%s", key)}
	statsd.Incr("kvexpress.out", statsdTags)
}

func DDAPIConnect(api, app, host string) *datadog.Client {
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

func DDStopEvent(dd *datadog.Client, key, value, direction string) {
	tags := makeTags(key, direction)
	tags = append(tags, "kvexpress:stop")
	title := fmt.Sprintf("Stop key is present: %s. Stopping.", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "error", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

func DDSaveDataEvent(dd *datadog.Client, key, value, direction string) {
	tags := makeTags(key, direction)
	tags = append(tags, "kvexpress:success")
	title := fmt.Sprintf("Updated: %s", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "info", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

func DDSaveStopEvent(dd *datadog.Client, key, value, direction string) {
	tags := makeTags(key, direction)
	tags = append(tags, "kvexpress:stop_set")
	title := fmt.Sprintf("Set Stop Key: %s", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "warning", Tags: tags}
	post, _ := dd.PostEvent(&event)
	if post != nil {

	}
}

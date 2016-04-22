// +build linux darwin freebsd

package commands

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"github.com/zorkian/go-datadog-api"
	"os"
)

// StatsdSetup sets up the connection to dogstatsd.
func StatsdSetup() *godspeed.Godspeed {
	statsd, err := godspeed.NewDefault()
	if err != nil {
		Log("StatsdSetup(): Problem setting up connection.", "info")
		return nil
	}
	return statsd
}

// StatsdIn sends metrics to Dogstatsd on a `kvexpress in` operation.
func StatsdIn(key string, dataLength int, data string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' stats='in'", DogStatsd, key), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, "complete")
			statsd.Incr("kvexpress.in", tags)
			statsd.Gauge("kvexpress.bytes", float64(dataLength), tags)
			// If the data is compressed - then LineCount will always return 1.
			// That's not useful or accurate, so let's decompress and count that.
			if Compress {
				data = DecompressData(data)
			}
			statsd.Gauge("kvexpress.lines", float64(LineCount(data)), tags)
		}
	}
}

// StatsdOut sends metrics to Dogstatsd on a `kvexpress out` operation.
func StatsdOut(key string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' stats='out'", DogStatsd, key), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, "complete")
			statsd.Incr("kvexpress.out", tags)
		}
	}
}

// StatsdLocked sends metrics to Dogstatsd on a `kvexpress out` operation
// that is blocked by a locked file.
func StatsdLocked(file string) {
	Log(fmt.Sprintf("dogstatsd='%t' file='%s' stats='locked'", DogStatsd, file), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(file, "complete")
			statsd.Incr("kvexpress.locked", tags)
		}
	}
}

// StatsdLength sends metrics to Dogstatsd on a `kvexpress out` operation
// where the file isn't long enough.
func StatsdLength(key string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' stats='not_long_enough'", DogStatsd, key), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, "not_long_enough")
			statsd.Incr("kvexpress.not_long_enough", tags)
		}
	}
}

// StatsdChecksum sends metrics to Dogstatsd on a `kvexpress out` operation
// where the checksum doesn't match.
func StatsdChecksum(key string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' stats='checksum_mismatch'", DogStatsd, key), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, "checksum_mismatch")
			statsd.Incr("kvexpress.checksum_mismatch", tags)
		}
	}
}

// StatsdLock sends metrics to Dogstatsd on a `kvexpress lock` operation.
func StatsdLock(key string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' stats='lock'", DogStatsd, key), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, "complete")
			statsd.Incr("kvexpress.lock", tags)
		}
	}
}

// StatsdUnlock sends metrics to Dogstatsd on a `kvexpress unlock` operation.
func StatsdUnlock(key string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' stats='unlock'", DogStatsd, key), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, "complete")
			statsd.Incr("kvexpress.unlock", tags)
		}
	}
}

// StatsdRaw sends metrics to Dogstatsd on a `kvexpress raw` operation.
func StatsdRaw(key string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' stats='raw'", DogStatsd, key), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, "complete")
			statsd.Incr("kvexpress.raw", tags)
		}
	}
}

// StatsdReconnect sends metrics when we have Consul connection retries.
func StatsdReconnect(times int) {
	Log(fmt.Sprintf("dogstatsd='%t' reconnect='%d'", DogStatsd, times), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := make([]string, 2)
			hostname := GetHostname()
			hostTag := fmt.Sprintf("host:%s", hostname)
			directionTag := fmt.Sprintf("direction:%s", Direction)
			tags = append(tags, hostTag)
			tags = append(tags, directionTag)
			statsd.Incr("kvexpress.consul_reconnect", tags)
		}
	}
}

// StatsdRunTime sends metrics to Dogstatsd on various operations.
func StatsdRunTime(key string, location string, msec int64) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' location='%s' msec='%d'", DogStatsd, key, location, msec), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, location)
			locationTag := fmt.Sprintf("location:%s", location)
			tags = append(tags, locationTag)
			statsd.Gauge("kvexpress.time", float64(msec), tags)
		}
	}
}

// StatsdPanic sends metrics to Dogstatsd when something really bad happens.
// It also stops the execution of kvexpress.
func StatsdPanic(key, location string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' location='%s' stats='panic'", DogStatsd, key, location), "debug")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, location)
			statsd.Incr("kvexpress.panic", tags)
		}
	}
	// If we're going to panic, we might as well stop right here.
	// Means we can't connect to Consul, download a URL or
	// write and/or chown files.
	os.Exit(0)
}

// StatsdConsul sends metrics to DogStatsd when Consul has a KV write or delete error.
func StatsdConsul(key, location string) {
	Log(fmt.Sprintf("dogstatsd='%t' key='%s' location='%s' stats='consul_error'", DogStatsd, key, location), "info")
	if DogStatsd {
		statsd := StatsdSetup()
		if statsd != nil {
			defer statsd.Conn.Close()
			tags := makeTags(key, location)
			statsd.Incr("kvexpress.consul_error", tags)
		}
	}
}

// DDAPIConnect connects to the Datadog API and returns a client object.
func DDAPIConnect(api, app string) *datadog.Client {
	client := datadog.NewClient(api, app)
	return client
}

// makeTags creates some standard tags for use with Dogstatsd and the Datadog API.
func makeTags(key, location string) []string {
	tags := make([]string, 4)
	keyTag := fmt.Sprintf("key:%s", key)
	hostname := GetHostname()
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

// DDStopEvent sends a Datadog event to the API when there's a stop key present.
func DDStopEvent(dd *datadog.Client, key, value string) {
	Log(fmt.Sprintf("datadog='true' DDStopEvent='true' key='%s'", key), "debug")
	tags := makeTags(key, "stop_key_present")
	tags = append(tags, "kvexpress:stop")
	title := fmt.Sprintf("Stop key is present: %s. Stopping.", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "error", Tags: tags}
	post, err := dd.PostEvent(&event)
	if (post == nil) || (err != nil) {
		Log("DDStopEvent(): Error posting to Datadog.", "info")
	}
}

// DDLengthEvent sends a Datadog event to the API when the file/url is too short.
func DDLengthEvent(dd *datadog.Client, key, value string) {
	Log(fmt.Sprintf("datadog='true' DDLengthEvent='true' key='%s'", key), "debug")
	tags := makeTags(key, "not_long_enough")
	tags = append(tags, "kvexpress:length")
	title := fmt.Sprintf("Not long enough: %s. Stopping.", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "error", Tags: tags}
	post, err := dd.PostEvent(&event)
	if (post == nil) || (err != nil) {
		Log("DDLengthEvent(): Error posting to Datadog.", "info")
	}
}

// DDSaveDataEvent sends a Datadog event to the API when we have updated a Consul key.
func DDSaveDataEvent(dd *datadog.Client, key, value string) {
	Log(fmt.Sprintf("datadog='true' DDSaveDataEvent='true' key='%s'", key), "debug")
	tags := makeTags(key, "complete")
	tags = append(tags, "kvexpress:success")
	title := fmt.Sprintf("Updated: %s", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "info", Tags: tags}
	post, err := dd.PostEvent(&event)
	if (post == nil) || (err != nil) {
		Log("DDSaveDataEvent(): Error posting to Datadog.", "info")
	}
}

// DDCopyDataEvent sends a Datadog event to the API when we have used `kvexpress copy`
// to copy a Consul key.
func DDCopyDataEvent(dd *datadog.Client, keyFrom, keyTo string) {
	Log(fmt.Sprintf("datadog='true' DDCopyDataEvent='true' keyFrom='%s' keyTo='%s'", keyFrom, keyTo), "debug")
	tags := makeTags(keyTo, "complete")
	tags = append(tags, "kvexpress:success")
	tags = append(tags, fmt.Sprintf("keyFrom:%s", keyFrom))
	title := fmt.Sprintf("Copy: %s to %s", keyFrom, keyTo)
	event := datadog.Event{Title: title, Text: title, AlertType: "info", Tags: tags}
	post, err := dd.PostEvent(&event)
	if (post == nil) || (err != nil) {
		Log("DDCopyDataEvent(): Error posting to Datadog.", "info")
	}
}

// DDSaveStopEvent sends a Datadog event when we have added a stop key to Consul.
func DDSaveStopEvent(dd *datadog.Client, key, value string) {
	Log(fmt.Sprintf("datadog='true' DDSaveStopEvent='true' key='%s'", key), "debug")
	tags := makeTags(key, "stop_key_save")
	tags = append(tags, "kvexpress:stop_set")
	title := fmt.Sprintf("Set Stop Key: %s", key)
	event := datadog.Event{Title: title, Text: value, AlertType: "warning", Tags: tags}
	post, err := dd.PostEvent(&event)
	if (post == nil) || (err != nil) {
		Log("DDSaveStopEvent(): Error posting to Datadog.", "info")
	}
}

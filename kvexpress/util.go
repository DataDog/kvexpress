package kvexpress

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
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

func GetOwnerId(owner string) int {
	var uid = ""
	usr, err := user.Lookup(owner)
	if err != nil {
		usr, _ = user.Current()
		uid = usr.Uid
		Log(fmt.Sprintf("out: owner='%s' status='not_found' uid='%s'", owner, uid), "debug")
	} else {
		uid = usr.Uid
		Log(fmt.Sprintf("out: owner='%s' status='found' uid='%s'", owner, uid), "debug")
	}
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	return int(uidInt)
}

func GetGroupId(group string) int {
	var gid = ""
	usr, err := user.Lookup(group)
	if err != nil {
		usr, _ = user.Current()
		gid = usr.Gid
		Log(fmt.Sprintf("out: group='%s' status='not_found' gid='%s'", group, gid), "debug")
	} else {
		gid = usr.Gid
		Log(fmt.Sprintf("out: group='%s' status='found' gid='%s'", group, gid), "debug")
	}
	gidInt, err := strconv.ParseInt(gid, 10, 64)
	return int(gidInt)
}

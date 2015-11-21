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

func GetCurrentUsername(direction string) string {
	usr, _ := user.Current()
	username := usr.Username
	Log(fmt.Sprintf("%s: username='%s'", direction, username), "debug")
	return username
}

func GetOwnerId(owner, direction string) int {
	var uid = ""
	var status = ""
	usr, err := user.Lookup(owner)
	if err != nil {
		usr, _ = user.Current()
		uid = usr.Uid
		status = "not_found"
	} else {
		uid = usr.Uid
		status = "found"
	}
	Log(fmt.Sprintf("%s: owner='%s' status='%s' uid='%s'", direction, owner, status, uid), "debug")
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	return int(uidInt)
}

func GetGroupId(owner, direction string) int {
	var gid = ""
	var status = ""
	usr, err := user.Lookup(owner)
	if err != nil {
		usr, _ = user.Current()
		gid = usr.Gid
		status = "not_found"
	} else {
		gid = usr.Gid
		status = "found"
	}
	Log(fmt.Sprintf("%s: owner='%s' status='%s' gid='%s'", direction, owner, status, gid), "debug")
	gidInt, err := strconv.ParseInt(gid, 10, 64)
	return int(gidInt)
}

// +build linux darwin freebsd

package commands

import (
	"fmt"
	"testing"
)

func TestKeyPath(t *testing.T) {
	PrefixLocation = "testing"
	fullPath := KeyPath("keyname", "data")
	if fullPath != "testing/keyname/data" {
		t.Error("Got the wrong path.")
	}
}

func TestFileLockPath(t *testing.T) {
	PrefixLocation = "testing"
	file := "/etc/datadog/hosts.consul"
	sha := "fba4f4f80fd22d7f7c26b00ad0a6c92d38c5f860870446eed4105ab170db2a9e"
	hostname := GetHostname()
	comparisonLockPath := fmt.Sprintf("%s/locks/%s/%s", PrefixLocation, sha, hostname)
	fileLockPath := FileLockPath(file)
	if fileLockPath != comparisonLockPath {
		t.Error("Got the wrong lock path.")
	}
}

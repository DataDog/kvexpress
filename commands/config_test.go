// +build linux darwin freebsd

package commands

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/smallfish/simpleyaml"
)

var yaml = `---
  datadog_api_key: api-key-goes-here
  datadog_app_key: app-key-goes-here
  consul_server: 127.0.0.1:8501
  token: abcd-efgh-hijk-lmno-pqrs-tuvw
  dogstatsd: true
  dogstatsd_address: 127.0.0.1:8125
  datadog_host: https://app.datadoghq.com`

var stringVars = []string{"datadog_api_key", "datadog_app_key", "consul_server", "token", "dogstatsd_address", "datadog_host"}

func loadTestConfigValues(data string) *simpleyaml.Yaml {
	yamlBytes := []byte(data)
	config := ParseConfig(yamlBytes)
	return config
}

func TestConfigValues(t *testing.T) {
	config := loadTestConfigValues(yaml)
	for _, configName := range stringVars {
		t.Logf("Getting configuration for '%s'", configName)
		configValue := GetStringConfig(config, configName)
		if configValue == "" {
			t.Errorf("Could not get the config for '%s'", configName)
		}
		t.Logf("Value: %s", configValue)
	}
}

func TestParseStatsdAddress(t *testing.T) {
	var goodAddrs = []string{"localhost:8100", "10.0.0.4:8100", "[2001:db8::1]:8080"}
	for _, addr := range goodAddrs {
		host, port := ParseStatsdAddress(addr)
		if host == "" {
			t.Errorf("Host should not be empty")
		}
		if reflect.TypeOf(port).Kind() != reflect.Int {
			t.Errorf("Port should be of type 'int'")
		}
		if &port == nil {
			t.Errorf("Port should not be empty")
		}

	}

	// test os.Exit() functionality
	if os.Getenv("CRASH_TEST") == "1" {
		ParseStatsdAddress("localhost")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestParseStatsdAddress")
	cmd.Env = append(os.Environ(), "CRASH_TEST=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		fmt.Printf("Error is %v\n", e)
		return
	}
	t.Fatalf("ParseStatsdAddress ran with err %v, want exit status 1", err)
}

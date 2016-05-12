// +build linux darwin freebsd

package commands

import (
	"github.com/smallfish/simpleyaml"
	"testing"
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

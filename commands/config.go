package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type kvexpressConfig struct {
	DatadogAPIKey    string `yaml:"datadog_api_key"`
	DatadogAPPKey    string `yaml:"datadog_app_key"`
	ConsulServer     string `yaml:"consul_server"`
	Token            string `yaml:"token"`
	DogStatsdAddress string `yaml:"dogstatsd_address"`
	DatadogHost      string `yaml:"datadog_host"`
}

func ConfigEnvVars(direction string) {
	if os.Getenv("CONSUL_SERVER") != "" {
		ConsulServer = os.Getenv("CONSUL_SERVER")
		kvexpress.Log(fmt.Sprintf("%s: Using CONSUL_SERVER ENV variable.", direction), "debug")
	}
	if os.Getenv("CONSUL_TOKEN") != "" {
		Token = os.Getenv("CONSUL_TOKEN")
		kvexpress.Log(fmt.Sprintf("%s: Using CONSUL_TOKEN ENV variable.", direction), "debug")
	}
	if os.Getenv("DATADOG_API_KEY") != "" {
		DatadogAPIKey = os.Getenv("DATADOG_API_KEY")
		kvexpress.Log(fmt.Sprintf("%s: Using DATADOG_API_KEY ENV variable.", direction), "debug")
	}
	if os.Getenv("DATADOG_APP_KEY") != "" {
		DatadogAPPKey = os.Getenv("DATADOG_APP_KEY")
		kvexpress.Log(fmt.Sprintf("%s: Using DATADOG_APP_KEY ENV variable.", direction), "debug")
	}
}

func setConfig(value, name string) {
	if value != "" {
		kvexpress.Log(fmt.Sprintf("config: Setting '%s' to '%s'", name, value), "info")
		name = value
	}
}

// Thanks https://mlafeldt.github.io/blog/decoding-yaml-in-go/ for a clear explanataion
// of how you did it.

func LoadConfig(filename string) {
	kvexpress.Log(fmt.Sprintf("config: filename='%s'", filename), "info")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte("")
	}
	var config kvexpressConfig
	if err := config.Parse(data); err != nil {
		kvexpress.Log(fmt.Sprintf("%v", err), "info")
	}
	if config.DatadogHost != "" {
		os.Setenv("DATADOG_HOST", config.DatadogHost)
	}
	setConfig(fmt.Sprintf("%s", config.DatadogAPIKey), "DatadogAPIKey")
	setConfig(fmt.Sprintf("%s", config.DatadogAPPKey), "DatadogAPPKey")
	setConfig(fmt.Sprintf("%s", config.ConsulServer), "ConsulServer")
	setConfig(fmt.Sprintf("%s", config.Token), "Token")
	setConfig(fmt.Sprintf("%s", config.DogStatsdAddress), "DogStatsdAddress")
}

func (c *kvexpressConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

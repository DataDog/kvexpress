// +build linux darwin freebsd

package commands

import (
	"fmt"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"os"
)

// GetStringConfig grabs the string from the config object.
func GetStringConfig(c *simpleyaml.Yaml, configValue string) string {
	result, err := c.Get(configValue).String()
	if err != nil {
		Log(fmt.Sprintf("Could not get '%s' from config.", configValue), "info")
	}
	if result != "" {
		return result
	}
	return ""
}

// ParseConfig takes the data from a file and parses the config.
func ParseConfig(data []byte) *simpleyaml.Yaml {
	config, err := simpleyaml.NewYaml(data)
	if err != nil {
		Log("Could not parse the configuration.", "info")
	}
	return config
}

// LoadConfig opens a file and reads the yaml formatted configuration data.
// It will set configuration globals and/or ENV variables as required.
func LoadConfig(filename string) {
	Log(fmt.Sprintf("config: filename='%s'", filename), "info")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte("")
	}
	config := ParseConfig(data)

	datadogHost := GetStringConfig(config, "datadog_host")
	if datadogHost != "" {
		os.Setenv("DATADOG_HOST", datadogHost)
	}

	datadogAPIKey := GetStringConfig(config, "datadog_api_key")
	if datadogAPIKey != "" {
		DatadogAPIKey = datadogAPIKey
	}

	datadogAPPKey := GetStringConfig(config, "datadog_app_key")
	if datadogAPPKey != "" {
		DatadogAPPKey = datadogAPPKey
	}

	token := GetStringConfig(config, "token")
	if token != "" {
		Token = token
	}

	consulServer := GetStringConfig(config, "consul_server")
	if consulServer != "" {
		ConsulServer = consulServer
	}

	dogstatsd, err := config.Get("dogstatsd").Bool()
	if err != nil {
		Log("Could not get 'dogstatsd' from config.", "info")
	}
	if dogstatsd {
		DogStatsd = true
	}

	dogstatsdAddress := GetStringConfig(config, "dogstatsd_address")
	if dogstatsdAddress != "" {
		DogStatsdAddress = dogstatsdAddress
	}

}

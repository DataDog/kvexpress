package commands

import (
	"fmt"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"os"
)

// RIPE for refactoring.

// LoadConfig opens a file and reads the yaml formatted configuration data.
// It will set configuration globals and/or ENV variables as required.
func LoadConfig(filename string) {
	Log(fmt.Sprintf("config: filename='%s'", filename), "info")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte("")
	}
	config, err := simpleyaml.NewYaml(data)

	datadogHost, _ := config.Get("datadog_host").String()

	if datadogHost != "" {
		os.Setenv("DATADOG_HOST", datadogHost)
	}

	datadogAPIKey, _ := config.Get("datadog_api_key").String()
	if datadogAPIKey != "" {
		DatadogAPIKey = datadogAPIKey
	}

	datadogAPPKey, _ := config.Get("datadog_app_key").String()
	if datadogAPPKey != "" {
		DatadogAPPKey = datadogAPPKey
	}

	token, _ := config.Get("token").String()
	if token != "" {
		Token = token
	}

	consulServer, _ := config.Get("consul_server").String()
	if consulServer != "" {
		ConsulServer = consulServer
	}

	dogstatsd, _ := config.Get("dogstatsd").Bool()
	if dogstatsd {
		DogStatsd = true
	}

	dogstatsdAddress, _ := config.Get("dogstatsd_address").String()
	if dogstatsdAddress != "" {
		DogStatsdAddress = dogstatsdAddress
	}

}

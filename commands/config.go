package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"os"
)

// RIPE for refactoring.

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

func LoadConfig(filename string) {
	kvexpress.Log(fmt.Sprintf("config: filename='%s'", filename), "info")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte("")
	}
	config, err := simpleyaml.NewYaml(data)

	datadog_host, _ := config.Get("datadog_host").String()

	if datadog_host != "" {
		os.Setenv("DATADOG_HOST", datadog_host)
	}

	datadog_api_key, _ := config.Get("datadog_api_key").String()
	if datadog_api_key != "" {
		DatadogAPIKey = datadog_api_key
	}

	datadog_app_key, _ := config.Get("datadog_app_key").String()
	if datadog_app_key != "" {
		DatadogAPPKey = datadog_app_key
	}

	token, _ := config.Get("token").String()
	if token != "" {
		Token = token
	}

	consul_server, _ := config.Get("consul_server").String()
	if consul_server != "" {
		ConsulServer = consul_server
	}

	dogstatsd, _ := config.Get("dogstatsd").Bool()
	if dogstatsd {
		DogStatsd = true
	}

	dogstatsd_address, _ := config.Get("dogstatsd_address").String()
	if dogstatsd_address != "" {
		DogStatsdAddress = dogstatsd_address
	}

}

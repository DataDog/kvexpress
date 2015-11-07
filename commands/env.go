package commands

import (
	kvexpress "../kvexpress/"
	"fmt"
	"os"
)

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

package commands

import (
	"log"
	"os"
)

func ConfigEnvVars(direction string) {
	if os.Getenv("CONSUL_SERVER") != "" {
		ConsulServer = os.Getenv("CONSUL_SERVER")
		log.Print(direction, ": Using CONSUL_SERVER ENV variable.")
	}
	if os.Getenv("CONSUL_TOKEN") != "" {
		Token = os.Getenv("CONSUL_TOKEN")
		log.Print(direction, ": Using CONSUL_TOKEN ENV variable.")
	}
	if os.Getenv("DATADOG_HOST") != "" {
		DatadogHost = os.Getenv("DATADOG_HOST")
		log.Print(direction, ": Using DATADOG_HOST ENV variable.")
	}
	if os.Getenv("DATADOG_API_KEY") != "" {
		DatadogAPIKey = os.Getenv("DATADOG_API_KEY")
		log.Print(direction, ": Using DATADOG_API_KEY ENV variable.")
	}
	if os.Getenv("DATADOG_APP_KEY") != "" {
		DatadogAPPKey = os.Getenv("DATADOG_APP_KEY")
		log.Print(direction, ": Using DATADOG_APP_KEY ENV variable.")
	}
}

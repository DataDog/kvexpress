package commands

import (
	"os"
)

// AutoEnable helps to automatically enable flags based on cues from the environment.
func AutoEnable() {
	// Check for dd-agent configuration file.
	if _, err := os.Stat("/etc/dd-agent/datadog.conf"); err == nil {
		DogStatsd = true
	}
}

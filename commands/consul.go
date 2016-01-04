// +build linux darwin freebsd

package commands

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"strings"
	"time"
)

const (
	consulTries = 5
)

// Connect sets up a connection to Consul.
func Connect(server string, token string) (*consul.Client, error) {
	consul, err := consulConnect(server, token)
	if err != nil {
		Log("Consul connection is bad.", "info")
	}
	return consul, err
}

// consulConnect to the Consul server and hand back a client object.
func consulConnect(server, token string) (*consul.Client, error) {
	var cleanedToken = ""
	config := consul.DefaultConfig()
	config.Address = server
	// Let's clean up the token so it doesn't appear in the logs.
	if token != "" && strings.ContainsAny(token, "-") {
		config.Token = token
		cleanedToken = cleanupToken(token)
	}
	// If it's anonymous - then that's a special case.
	if token == "anonymous" {
		config.Token = token
		cleanedToken = token
	}
	consul, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	Log(fmt.Sprintf("server='%s' token='%s'", server, cleanedToken), "debug")
	return consul, nil
}

// Standard Consul tokens have lots of dashes in them - let's split on the dash
// so that we can see part of the token in the logs - helps with debugging.
// We don't want the full token in any logs - that's bad.
func cleanupToken(token string) string {
	first := strings.Split(token, "-")
	firstString := fmt.Sprintf("%s", first[0])
	Log(firstString, "info")
	return firstString
}

// Get the value from a key in the Consul KV store.
func Get(c *consul.Client, key string) string {
	maxTries := consulTries
	for tries := 1; tries <= maxTries; tries++ {
		value, err := consulGet(c, key)
		if err == nil {
			return value
		}
		waitTime := time.Duration(tries) * time.Second
		Log(fmt.Sprintf("Consul Get failure (%d) - trying again. Max: %d", tries, maxTries), "info")
		if tries < maxTries {
			time.Sleep(waitTime)
		}
	}
	Log("Panic: Giving up on Consul Get.", "info")
	fmt.Printf("Panic: Could not get from Consul: '%s'\n", key)
	StatsdPanic(key, "consul_get")
	return ""
}

// consulGet the value from a key in the Consul KV store.
func consulGet(c *consul.Client, key string) (string, error) {
	var value string
	kv := c.KV()
	key = strings.TrimPrefix(key, "/")
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return "", err
	}
	if pair != nil {
		value = string(pair.Value[:])
	} else {
		value = ""
	}
	Log(fmt.Sprintf("action='consulGet' key='%s'", key), "debug")
	return value, err
}

// Set the value for a key in the Consul KV store.
func Set(c *consul.Client, key string, value string) bool {
	maxTries := consulTries
	for tries := 1; tries <= maxTries; tries++ {
		if consulSet(c, key, value) {
			return true
		}
		waitTime := time.Duration(tries) * time.Second
		Log(fmt.Sprintf("Consul Set failure (%d) - trying again. Max: %d", tries, maxTries), "info")
		if tries < maxTries {
			time.Sleep(waitTime)
		}
	}
	Log("Panic: Giving up on Consul Set.", "info")
	fmt.Printf("Panic: Could not set in Consul: '%s'\n", key)
	StatsdPanic(key, "consul_set")
	return false
}

// consulSet a value for a key in the Consul KV store.
func consulSet(c *consul.Client, key string, value string) bool {
	key = strings.TrimPrefix(key, "/")
	p := &consul.KVPair{Key: key, Value: []byte(value)}
	kv := c.KV()
	_, err := kv.Put(p, nil)
	if err != nil {
		return false
	}
	Log(fmt.Sprintf("action='consulSet' key='%s'", key), "debug")
	return true
}

// Del removes a key from the Consul KV store.
func Del(c *consul.Client, key string) bool {
	maxTries := consulTries
	for tries := 1; tries <= maxTries; tries++ {
		if consulDel(c, key) {
			return true
		}
		waitTime := time.Duration(tries) * time.Second
		Log(fmt.Sprintf("Consul Del failure (%d) - trying again. Max: %d", tries, maxTries), "info")
		if tries < maxTries {
			time.Sleep(waitTime)
		}
	}
	Log("Panic: Giving up on Consul Del.", "info")
	fmt.Printf("Panic: Could not delete from Consul: '%s'\n", key)
	StatsdPanic(key, "consul_del")
	return false
}

// consulDel removes a key from the Consul KV store.
func consulDel(c *consul.Client, key string) bool {
	kv := c.KV()
	key = strings.TrimPrefix(key, "/")
	_, err := kv.Delete(key, nil)
	if err != nil {
		return false
	}
	Log(fmt.Sprintf("action='consulDel' key='%s'", key), "info")
	return true
}

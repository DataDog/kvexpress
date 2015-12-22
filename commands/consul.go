package commands

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"strings"
)

// Connect to the Consul server and hand back a client object.
func Connect(server, token string) (*consul.Client, error) {
	var cleanedToken = ""
	config := consul.DefaultConfig()
	config.Address = server
	if token != "" {
		config.Token = token
		cleanedToken = cleanupToken(token)
	}
	consul, _ := consul.NewClient(config)
	Log(fmt.Sprintf("server='%s' token='%s'", server, cleanedToken), "debug")
	return consul, nil
}

// TODO: Can likely refactor the following two functions into one.

// Get the value from a kvexpress formatted key in the Consul KV store.
func Get(c *consul.Client, key string) string {
	var value string
	kv := c.KV()
	key = strings.TrimPrefix(key, "/")
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		Log(fmt.Sprintf("action='get' panic='true' key='%s'", key), "info")
		StatsdPanic(key, "consul_get")
	} else {
		if pair != nil {
			value = string(pair.Value[:])
		} else {
			value = ""
		}
		Log(fmt.Sprintf("action='get' key='%s'", key), "debug")
	}
	return value
}

// GetRaw the value from any key in the Consul KV store.
func GetRaw(c *consul.Client, key string) string {
	var value string
	kv := c.KV()
	fullKey := fmt.Sprintf("%s/%s", PrefixLocation, key)
	fullKey = strings.TrimPrefix(fullKey, "/")
	pair, _, err := kv.Get(fullKey, nil)
	if err != nil {
		Log(fmt.Sprintf("action='get_raw' panic='true' key='%s'", fullKey), "info")
		StatsdPanic(fullKey, "consul_get_raw")
	} else {
		if pair != nil {
			value = string(pair.Value[:])
		} else {
			value = ""
		}
		Log(fmt.Sprintf("action='get_raw' key='%s'", fullKey), "debug")
	}
	return value
}

func cleanupToken(token string) string {
	first := strings.Split(token, "-")
	firstString := fmt.Sprintf("%s", first[0])
	return firstString
}

// Set a value in a kvexpress formatted key in the Consul KV store.
func Set(c *consul.Client, key string, value string) bool {
	key = strings.TrimPrefix(key, "/")
	p := &consul.KVPair{Key: key, Value: []byte(value)}
	kv := c.KV()
	_, err := kv.Put(p, nil)
	if err != nil {
		Log(fmt.Sprintf("action='set' panic='true' key='%s'", key), "info")
		StatsdPanic(key, "consul_set")
	} else {
		Log(fmt.Sprintf("action='set' key='%s'", key), "debug")
		return true
	}
	return true
}

// Del removes a key from the Consul KV store.
func Del(c *consul.Client, key string) bool {
	kv := c.KV()
	key = strings.TrimPrefix(key, "/")
	_, err := kv.Delete(key, nil)
	if err != nil {
		Log(fmt.Sprintf("action='Del' panic='true' key='%s'", key), "info")
		StatsdPanic(key, "consul_del")
		return false
	}
	Log(fmt.Sprintf("action='Del' panic='false' key='%s'", key), "info")
	return true
}

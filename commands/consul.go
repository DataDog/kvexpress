package commands

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"strings"
)

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

func Get(c *consul.Client, key string, dogstatsd bool) string {
	var value string
	kv := c.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		Log(fmt.Sprintf("action='get' panic='true' key='%s'", key), "info")
		if dogstatsd {
			StatsdPanic(key, "consul_get")
		}
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

func GetRaw(c *consul.Client, prefix string, key string, dogstatsd bool) string {
	var value string
	kv := c.KV()
	full_key := fmt.Sprintf("%s/%s", prefix, key)
	pair, _, err := kv.Get(full_key, nil)
	if err != nil {
		Log(fmt.Sprintf("action='get_raw' panic='true' key='%s'", full_key), "info")
		if dogstatsd {
			StatsdPanic(full_key, "consul_get_raw")
		}
	} else {
		if pair != nil {
			value = string(pair.Value[:])
		} else {
			value = ""
		}
		Log(fmt.Sprintf("action='get_raw' key='%s'", full_key), "debug")
	}
	return value
}

func cleanupToken(token string) string {
	first := strings.Split(token, "-")
	firstString := fmt.Sprintf("%s", first[0])
	return firstString
}

func Set(c *consul.Client, key string, value string, dogstatsd bool) bool {
	p := &consul.KVPair{Key: key, Value: []byte(value)}
	kv := c.KV()
	_, err := kv.Put(p, nil)
	if err != nil {
		Log(fmt.Sprintf("action='set' panic='true' key='%s'", key), "info")
		if dogstatsd {
			StatsdPanic(key, "consul_set")
		}
	} else {
		Log(fmt.Sprintf("action='set' key='%s'", key), "debug")
		return true
	}
	return true
}

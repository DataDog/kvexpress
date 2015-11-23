package kvexpress

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"strings"
)

func Connect(server, token, direction string) (*consul.Client, error) {
	var cleanedToken = ""
	config := consul.DefaultConfig()
	config.Address = server
	if token != "" {
		config.Token = token
		cleanedToken = cleanupToken(token)
	}
	consul, _ := consul.NewClient(config)
	Log(fmt.Sprintf("%s: server='%s' token='%s'", direction, server, cleanedToken), "debug")
	return consul, nil
}

func Get(c *consul.Client, key, direction string) string {
	var value string
	kv := c.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		Log(fmt.Sprintf("%s: action='get' panic='true' key='%s'", direction, key), "info")
	} else {
		if pair != nil {
			value = string(pair.Value[:])
		} else {
			value = ""
		}
		Log(fmt.Sprintf("%s: action='get' key='%s'", direction, key), "debug")
	}
	return value
}

func cleanupToken(token string) string {
	first := strings.Split(token, "-")
	firstString := fmt.Sprintf("%s", first[0])
	return firstString
}

func Set(c *consul.Client, key, value, direction string) bool {
	p := &consul.KVPair{Key: key, Value: []byte(value)}
	kv := c.KV()
	_, err := kv.Put(p, nil)
	if err != nil {
		Log(fmt.Sprintf("%s: action='set' panic='true' key='%s'", direction, key), "info")
	} else {
		Log(fmt.Sprintf("%s: action='set' key='%s'", direction, key), "debug")
		return true
	}
	return true
}

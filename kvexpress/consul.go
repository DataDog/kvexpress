package kvexpress

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
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
	log.Print(direction, ": address='", server, "' token='", cleanedToken, "'")
	return consul, nil
}

func Get(c *consul.Client, key, direction string) string {
	var value string
	kv := c.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		panic(err)
	} else {
		if pair != nil {
			value = string(pair.Value[:])
		} else {
			value = ""
		}
		log.Print(direction, ": key='", key)
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
		panic(err)
	} else {
		return true
	}
	return true
}

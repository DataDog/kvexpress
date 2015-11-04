package kvexpress

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"strings"
)

func Get(key string, server string, token string, direction string) string {
	var value string
	config := consulapi.DefaultConfig()
	config.Address = server
	if token != "" {
		config.Token = token
	}
	consul, err := consulapi.NewClient(config)
	kv := consul.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		panic(err)
	} else {
		if pair != nil {
			value = string(pair.Value[:])
		} else {
			value = ""
		}
		cleanedToken := cleanupToken(token)
		log.Print(direction, ": key='", key, "' address='", server, "' token='", cleanedToken, "'")
	}
	return value
}

func cleanupToken(token string) string {
	first := strings.Split(token, "-")
	firstString := fmt.Sprintf("%s", first[0:1])
	return firstString
}

func Set(key string, value string, server string, token string, direction string) bool {
	config := consulapi.DefaultConfig()
	config.Address = server
	if token != "" {
		config.Token = token
	}
	consul, err := consulapi.NewClient(config)
	kv := consul.KV()
	p := &consulapi.KVPair{Key: key, Value: []byte(value)}
	_, err = kv.Put(p, nil)
	if err != nil {
		panic(err)
	} else {
		return true
	}
}

package kvexpress

import (
	consulapi "github.com/hashicorp/consul/api"
	"log"
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
		log.Print(direction, ": key='", key, "' value='", value, "' address='", server, "' token='", token, "'")
	}
	return value
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

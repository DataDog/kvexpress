// +build linux darwin freebsd

package commands

import (
	"fmt"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
)

const (
	consulTries = 5
)

// Connect sets up a connection to Consul.
func Connect(server string, token string) (*consul.Client, error) {
	consul, err := consulConnect(server, token)
	if err != nil {
		Log("Consul connection is bad.", "info")
		return nil, err
	}
	return consul, err
}

// consulConnect to the Consul server and hand back a client object.
func consulConnect(server, token string) (*consul.Client, error) {
	config := consul.DefaultConfig()
	config.Address = server
	// Let's clean up the token so it doesn't appear in the logs.
	if token != "" {
		config.Token = token
	}
	consul, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	Log(fmt.Sprintf("server='%s' token='%s'", server, cleanupToken(token)), "debug")
	return consul, nil
}

// Standard Consul tokens have lots of dashes in them - let's split on the dash
// so that we can see part of the token in the logs - helps with debugging.
// We don't want the full token in any logs - that's bad.
func cleanupToken(token string) string {
	first := strings.SplitN(token, "-", 2)[0]
	return first
}

// Get the value from a key in the Consul KV store.
func Get(c *consul.Client, key string) string {
	data, _ := GetRaw(c, key)
	return string(data)
}

// Get the raw value + flags from a key in the Consul KV store.
func GetRaw(c *consul.Client, key string) ([]byte, uint64) {
	var data []byte
	var flags uint64
	Retry(func() error {
		pair, err := consulGet(c, key)
		if pair != nil {
			data = pair.Value[:]
			flags = pair.Flags
		}
		return err
	}, consulTries)
	return data, flags
}

// Retry loops through the callback func and tries several times to do the thing.
func Retry(callback func() error, tries int) {
	var err error
	for i := 1; i <= tries; i++ {
		err = callback()
		if err == nil {
			return
		}
		waitTime := time.Duration(tries) * time.Second
		Log(fmt.Sprintf("Consul Failure (%d) - trying again. Max: %d. Error: %s", i, tries, err), "info")
		StatsdReconnect(tries)
		if i < tries {
			time.Sleep(waitTime)
		}
	}
	LogFatal("Panic: Giving up on Consul.", "nokey", "no_more_retries")
}

// consulGet the value from a key in the Consul KV store.
func consulGet(c *consul.Client, key string) (*consul.KVPair, error) {
	kv := c.KV()
	key = strings.TrimPrefix(key, "/")
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	Log(fmt.Sprintf("action='consulGet' key='%s'", key), "debug")
	return pair, err
}

// Set the value for a key in the Consul KV store.
func Set(c *consul.Client, key string, value string) bool {
	return SetRaw(c, key, []byte(value), 0)
}

// Set the value for a key in the Consul KV store.
func SetRaw(c *consul.Client, key string, value []byte, flags uint64) bool {
	var success bool
	Retry(func() error {
		var err error
		success, err = consulSet(c, key, value, flags)
		if success != true {
			StatsdConsul(key, "set")
		}
		return err
	}, consulTries)
	return true
}

// consulSet a value for a key in the Consul KV store.
func consulSet(c *consul.Client, key string, value []byte, flags uint64) (bool, error) {
	key = strings.TrimPrefix(key, "/")
	p := &consul.KVPair{Key: key, Value: value, Flags: flags}
	kv := c.KV()
	_, err := kv.Put(p, nil)
	if err != nil {
		return false, err
	}
	Log(fmt.Sprintf("action='consulSet' key='%s'", key), "debug")
	return true, err
}

// Del removes a key from the Consul KV store.
func Del(c *consul.Client, key string) bool {
	var success bool
	Retry(func() error {
		var err error
		success, err = consulDel(c, key)
		if success != true {
			StatsdConsul(key, "delete")
		}
		return err
	}, consulTries)
	return true
}

// consulDel removes a key from the Consul KV store.
func consulDel(c *consul.Client, key string) (bool, error) {
	kv := c.KV()
	key = strings.TrimPrefix(key, "/")
	_, err := kv.Delete(key, nil)
	if err != nil {
		return false, err
	}
	Log(fmt.Sprintf("action='consulDel' key='%s'", key), "info")
	return true, err
}

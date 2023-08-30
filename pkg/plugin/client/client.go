package client

import (
	"fmt"
	"strings"
	"sync"
)

type PluginClientRegistry struct {
	client sync.Map
}

func (c *PluginClientRegistry) RegisterClient(id string, client interface{}) error {
	if client == nil {
		return fmt.Errorf("register client fail: %s", "client is nil")
	}
	_, loaded := c.client.LoadOrStore(id, client)
	if loaded {
		return fmt.Errorf("client %s already exists", id)
	}
	return nil
}

func (c *PluginClientRegistry) GetPluginClient(id string) any {
	value, _ := c.client.Load(id)
	if value == nil {
		return nil
	}
	return value
}

func (c *PluginClientRegistry) DeletePluginClient(id string) {
	c.client.Delete(id)
}

func (c *PluginClientRegistry) RangePluginClient(name string) interface{} {
	var client interface{}
	c.client.Range(func(key, value interface{}) bool {
		parts := strings.SplitN(key.(string), "-", 2)
		prefix := parts[0]
		if prefix == name {
			client = value
			return false
		}
		return true
	})
	return client
}

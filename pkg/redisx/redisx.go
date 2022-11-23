package redisx

import (
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// NewClient 从配置信息中创建client
func NewClient(conf map[string]interface{}) redis.UniversalClient {
	if mode, ok := conf["mode"].(string); ok {
		if mode == "cluster" {
			return redis.NewClusterClient(newClusterOptions(conf))
		}
	}

	return redis.NewClient(newOptions(conf))
}

// newClusterOptions 集群的Options
func newClusterOptions(conf map[string]interface{}) *redis.ClusterOptions {
	options := &redis.ClusterOptions{}
	for k, v := range conf {
		switch strings.ToLower(k) {
		case "addrs":
			_addrs := v.([]interface{})
			for _, s := range _addrs {
				options.Addrs = append(options.Addrs, s.(string))
			}
		}
	}

	return options
}

// newOptions 普通Options
func newOptions(conf map[string]interface{}) *redis.Options {
	options := &redis.Options{}

	for k, v := range conf {
		switch strings.ToLower(k) {
		case "addr":
			options.Addr = v.(string)
		case "minidleconns":
			options.MinIdleConns = v.(int)
		case "network":
			options.Network = v.(string)
		case "pooltimeout":
			options.PoolTimeout = time.Duration(v.(int)) * time.Second
		}
	}

	return options
}

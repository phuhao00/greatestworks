package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/phuhao00/broker"
)

type ClusterClient struct {
	*broker.BaseComponent
	real *redis.ClusterClient
}

func NewClusterClient(conf *Config) *ClusterClient {
	c := &ClusterClient{
		BaseComponent: broker.NewBaseComponent(),
		real: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: conf.Addrs,

			//Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
			// To route commands by latency or randomly, enable one of the following.
			//RouteByLatency: true,
			//RouteRandomly: true,
		}),
	}
	return c
}

package redis

import (
	"github.com/phuhao00/broker/redis"
	"sync"
)

var (
	mockInstance *redis.ClusterClient
	onceMockInit sync.Once
)

func GetMockInstance() *redis.ClusterClient {
	onceMockInit.Do(func() {
		conf := &redis.Config{Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"}}
		mockInstance = redis.NewClusterClient(conf)
	})
	return mockInstance
}

package redis

import (
	"time"
)

type Config struct {
	Addrs              []string      `json:"addrs"`
	MaxRedirects       int           `json:"maxRedirects"`
	ReadOnly           bool          `json:"readOnly"`
	RouteByLatency     bool          `json:"routeByLatency"`
	RouteRandomly      bool          `json:"routeRandomly"`
	Username           string        `json:"username"`
	Password           string        `json:"password"`
	MaxRetries         int           `json:"maxRetries"`
	MinRetryBackoff    time.Duration `json:"minRetryBackoff"`
	MaxRetryBackoff    time.Duration `json:"maxRetryBackoff"`
	DialTimeout        time.Duration `json:"dialTimeout"`
	ReadTimeout        time.Duration `json:"readTimeout"`
	WriteTimeout       time.Duration `json:"writeTimeout"`
	PoolFIFO           bool          `json:"poolFIFO"`
	PoolSize           int           `json:"poolSize"`
	MinIdleConns       int           `json:"minIdleConns"`
	MaxConnAge         time.Duration `json:"maxConnAge"`
	PoolTimeout        time.Duration `json:"poolTimeout"`
	IdleTimeout        time.Duration `json:"idleTimeout"`
	IdleCheckFrequency time.Duration `json:"idleCheckFrequency"`
}

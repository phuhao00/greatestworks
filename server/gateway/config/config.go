package config

import "greatestworks/aop/redis"

// GlobalConfig ...
type GlobalConfig struct {
	ServiceUpdateTime int
	IsOpenNow         bool
	ZoneId            int
	RedisInfo         *redis.Config
}

// ServerConfig ...
type ServerConfig struct {
	PrivateIP       string
	PublicIP        string
	Port            int
	InnerPort       int
	MaxConnNum      int
	PriMsgBuffSize  string
	PriConnBuffSize string
	PubMsgBuffSize  string
	PubConnBuffSize string
}

// LogConfig ...
type LogConfig struct {
	LogPath  string
	LogFile  string
	LogLevel string
}

// HTTPConfig ...
type HTTPConfig struct {
	HTTPAddr    string
	HTTPPort    int
	TLSCertFile *string
	TLSKeyFile  *string
}

// Config ...
type Config struct {
	Global       *GlobalConfig
	Server       *ServerConfig
	Log          *LogConfig
	HTTP         *HTTPConfig
	UpdateInfoCd int64
	DeploymentId string
	NodeName     string
	Version      string
}

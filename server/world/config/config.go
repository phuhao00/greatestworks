package config

import "greatestworks/aop/redis"

type Config struct {
	MaxPlayerNum int32
	RpcAddr      string
	HttpAddress  string
	Global       *GlobalConfig
	Server       *ServerConfig
	Log          *LogConfig
	HTTP         *HTTPConfig
	RpcServer    *RpcConfig
	Stat         *StatConfig
	Settings     *SettingsConfig
}

type Global struct {
	ServerType int
}

type GlobalConfig struct {
	Mongo                string
	MongoName            string
	NsqLookup            string
	NsqPoolSize          int
	MongoWinSize         int
	ServiceDiscoveryTime int
	ServiceUpdateTime    int
	LoadBalanceRatio     int32
	ZoneId               int
	RedisInfo            *redis.Config
	CtrlPlayersServerCnt int32 // 开服的时候有多少个排头的服做人数上浮 - 避免排头的服瞬间爆掉
	CtrlPlayersDeltaCnt  int32 // 人数上浮时，按照服爆满值(阻拦玩家进入)减少的数
	CtrlPlayerNumHour    int32 // 上浮的时限制， 从该服启动时开始算。
	EnablePayWhitelist   bool  // 重置白名单检测开启
	ReservationReward    bool  // 官网预约奖励
	MMIDBindRewardClose  bool  // 米米号绑定检查
	InvitePlayerClose    bool  // 人拉人奖励
	ServerType           int   // 服务器类型 0:官服 1:渠道服
}

// ServerConfig ...
type ServerConfig struct {
	MsgBuffSize  string
	ConnBuffSize string
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

// RpcConfig ...
type RpcConfig struct {
	RpcIp   string
	RpcPort int
}

// StatConfig ...
type StatConfig struct {
	GameID       int
	ReqUrl       string
	ReqCustomUrl string
}

// SettingsConfig  ...
type SettingsConfig struct {
	GMCommand bool
}

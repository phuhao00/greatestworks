package config

import "time"

type Me struct {
	Name                 string
	HTTPAddr             string  `json:"http_addr"`
	HTTPPort             int     `json:"http_port"`
	TLSCertFile          string  `json:"tls_cert_file"`
	TLSKeyFile           string  `json:"tls_key_file"`
	LimitGoroutinesNum   int     `json:"limit_goroutines_num"`
	ServiceDiscoveryTime int     `json:"service_discovery_time"`
	LoadBalanceRatio     int32   `json:"load_balance_ratio"`
	WindowSize           int     `json:"window_size"`
	RetryRatio           float64 `json:"retry_ratio"`
	TimeStamp            int     `json:"time_stamp"`
	EmptyRatioValue      float32 `json:"empty_ratio_value"`
	BusyRatioValue       float32 `json:"busy_ratio_value"`
	MaxLoginQueueLength  int64   `json:"max_login_queue_length"`
	EnableLoginQueue     bool    `json:"enable_login_queue"`
	PlayersServerCnt     int32   `json:"players_server_cnt"`
	PlayersDeltaCnt      int32   `json:"players_delta_cnt"`
	PlayerNumHour        int32   `json:"player_num_hour"`
	QueryGateWayRatio    int     `json:"query_gate_way_ratio"`
	MaxWorldPlayerNum    uint32  `json:"max_world_player_num"`
}

type EndPoint struct {
	ZoneId  int
	ID      string
	IP      string
	Port    int
	Name    string
	Weights int
	InnerIP string
}

const (
	MAXHold = 1000
	LEVEL0  = int(0.01 * float32(MAXHold))
	LEVEL1  = int(0.3 * float32(MAXHold))
	LEVEL2  = int(0.5 * float32(MAXHold))
	LEVEL3  = int(0.9 * float32(MAXHold))
)

var QueryToGateWayRatio = int(3)

const (
	GateWayServiceName = "GateWay-Tcp"
	WorldServiceName   = "World-http"
)

const (
	HoursSeconds = 60 * 60

	WorldMaxCoefficient = 0.7
)

const (
	CloseStatus = int32(0)
	EmptyStatus = int32(1)
	OKStatus    = int32(2)
	FullStatus  = int32(3)

	RecommendWorldMaxCnt = 5
)

const (
	LoginWindowSize      = int64(100)
	EmptyRatio           = float32(0.2)
	BusyRatio            = float32(0.7)
	StartPreRegisterTime = int64(0)
	EndPreRegisterTime   = int64(0)
	ServerOpenTime       = int64(0)

	LoginTimeStampMcu = 60
)

const (
	WaitTimeLimit    int64 = 10
	WaitTimeBitLen   int   = 20
	MaxRetryWaitTime int   = 30
)

var (
	QueueLength         int64 = 0       // 排队队列长度
	MaxLoginQueueLength int64 = 1000000 // 排队的最大排队长度 100w

	DailyLimitKey = "LoginDailyRegisterInfo"
)

const (
	Succ           = 0  //
	DecodeErr      = 1  // josn
	WhiteListErr   = 2  // 白名单
	VerifyTokenErr = 3  // Token验证
	BanUserErr     = 4  // 封号
	UnknownErr     = 5  //
	PassWdErr      = 6  //
	LoginBusy      = 7  // 登录服繁忙，稍后再来
	ZoneError      = 8  // 指定的Zone不存在
	DailyIncrOver  = 9  // 日新增达到上限，不能注册
	OnlineError    = 10 // 指定的online不存在
	BanMacLoginErr = 11 // 封Mac登录
	BanMacChatErr  = 12 // 封Mac聊天
)
const (
	TokenExpireDuration = 24 * time.Hour
)

type AccountData struct {
	ZoneId   int
	Account  string
	Password string
	Sign     string
	Sid      string
	Token    string
}

type LimitInfo struct {
	Cond      uint32 `bson:"con"` // 原因
	InReason  string `bson:"ir"`  // 内部原因
	OutReason string `bson:"or"`  // 内部原因
	ExpTime   int64  `bson:"etm"` // 过期时间
}

func GetIdxByLimitInfo(source []LimitInfo, cond uint32) (bool, int) {
	for idx, val := range source {
		if val.Cond == cond {
			return true, idx
		}
	}
	return false, 0
}

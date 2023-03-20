package config

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
	TimeSlice            int     `json:"time_slice"`
	EmptyRatioValue      float32 `json:"empty_ratio_value"`
	BusyRatioValue       float32 `json:"busy_ratio_value"`
	MaxLoginQueueLength  int64   `json:"max_login_queue_length"`
	EnableLoginQueue     bool    `json:"enable_login_queue"`
	PlayersServerCnt     int32   `json:"players_server_cnt"`
	PlayersDeltaCnt      int32   `json:"players_delta_cnt"`
	PlayerNumHour        int32   `json:"player_num_hour"`
	QueryGateWayRatio    int32   `json:"query_gate_way_ratio"`
}

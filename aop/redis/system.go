package redis

const (
	SystemCloseList       = "system_close_list"        // 系统关闭
	ServerMaxPlayer       = "server_max_player"        // 最大人数设置
	LoginWindowCtr        = "login_window_ctr"         // 登录限流窗口
	ServerPlayerNumRate   = "server_player_num_Rate"   // 服务器人数系数
	ServerEmptyRate       = "server_empty_Rate"        // 服务器空闲系数
	ServerOkRate          = "server_ok_Rate"           // 服务器繁忙系数
	ServerFullRate        = "server_full_Rate"         // 服务器爆满系数
	ServerPreRegisterTime = "server_pre_register_time" // 服务器开启预创角的时间
	ServerOpenTime        = "server_open_time"         // 服务器开启的时间
)

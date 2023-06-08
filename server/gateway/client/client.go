package client

import (
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/network"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	Conn   *network.TcpConnX
	Router *fuse.Router // 消息路由器

	CharacterId          uint64        // 角色Id
	WorldServerId        atomic.Value  // string
	IsBindWorldServer    atomic.Value  // bool
	JoinWorldServerId    atomic.Value  // string
	IsDisconnected       atomic.Value  // bool
	DisconnectedTime     atomic.Value  // time.Time
	DisconnectedMessages []interface{} // 断线缓存消息
	IsReconnection       atomic.Value  // bool
	RemoteIp             string        // 对端IP
	LastPingTime         time.Time     // 上次ping的时间
	LastCheckTime        int64         // 上次收到数据的时间点
	RequestSpeed         int           // 平均上行流速
	ReqFrequency         int           // 平均上行评率
	MsgRegisterOnce      sync.Once     // 代理消息注册一次
	ProcIndex            uint32        // 服id编号
}

func (c *Client) Init() {
	now := time.Now()
	c.Router = fuse.NewRouter()
	c.CharacterId = 0
	c.LastCheckTime = now.Unix()
	c.RequestSpeed = 0
	c.ReqFrequency = 0
	c.WorldServerId.Store("")
	c.JoinWorldServerId.Store("")
	c.IsBindWorldServer.Store(false)
	c.IsDisconnected.Store(false)
	c.IsReconnection.Store(false)
	c.DisconnectedTime.Store(now)
	c.LastPingTime = now
}

func (c *Client) HandlerRegister() {

}

func (s *Client) Loop() {

	for {
		select {
		//impl Message

		}
	}
}

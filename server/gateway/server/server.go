package server

import (
	"context"
	"github.com/phuhao00/broker/timerassistant"
	_ "github.com/phuhao00/broker/timerassistant"
	"github.com/phuhao00/network"
	"greatestworks/aop/consul"
	"greatestworks/aop/fn"
	"greatestworks/aop/idgenerator"
	"greatestworks/aop/logger"
	"greatestworks/aop/redis"
	"greatestworks/server"
	"greatestworks/server/gateway/client"
	"greatestworks/server/gateway/config"
	"greatestworks/server/gateway/gm"
	"math/rand"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

//todo 字节对其

type Server struct {
	*server.BaseService
	Config          *config.Config
	id              string
	runPath         string
	processIdx      int
	tcpAddr         string
	tcpPort         int
	rpcPort         int
	httpAddr        string
	httpPort        int
	tcpSvcID        string
	httpSvcID       string
	innerPort       int
	innerSvcAddr    string
	innerSvcID      string
	startTM         int64
	httpHandler     *HTTPHandler
	tcpServer       *network.TcpServer
	innerServer     *network.TcpServer
	timer           *timerassistant.TimerNormalAssistant
	lastUpdateCount uint32
	PriMsgBuffSize  int64
	PriConnBuffSize int64
	PubMsgBuffSize  int64
	PubConnBuffSize int64
}

var (
	once           sync.Once
	InstanceServer *Server
)

const (
	httpService     = "gateway-http"
	tcpService      = "gateway-tcp"
	innerTCPService = "gateway-inner-tcp"
)

func GetServer() *Server {
	once.Do(func() {
		InstanceServer = &Server{
			tcpServer:   &network.TcpServer{},
			innerServer: &network.TcpServer{},
			timer:       timerassistant.NewTimerNormalAssistant(time.Second),
		}
		baseService, err := server.NewBaseService(InstanceServer.GetName(), InstanceServer.GetDeploymentId())
		if err != nil {
			panic(err)
		}
		InstanceServer.BaseService = baseService
	})
	return InstanceServer
}

func (s *Server) GetName() string {
	return s.Name
}

func (s *Server) GetDeploymentId() string {
	return s.DeploymentId
}

func (s *Server) GetTcpAddress() string {
	return s.tcpAddr
}

func (s *Server) clearConnRecords() {
	key := redis.MakeGatewayKey(s.tcpAddr)
	records := redis.CacheRedis().SMembers(context.TODO(), key)
	for _, idStr := range records.Val() {
		uid, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			continue
		}

		key := redis.MakeAccountKey(int64(uid))
		ret := redis.CacheRedis().HGet(context.TODO(), key, "Gateway")
		if ret.Val() == s.tcpAddr {
			redis.CacheRedis().HDel(context.TODO(), key, "Gateway", "ZoneId")
			logger.Debug("[clearConnRecords] redis HDel ZoneId of %v", key)
		}
	}
	redis.CacheRedis().Del(context.TODO(), key)
}

func (s *Server) registerTimer() {
	if s.Config.Global.ServiceUpdateTime <= 0 {
		s.Config.Global.ServiceUpdateTime = 2
	}

	suDuration := s.Config.Global.ServiceUpdateTime

	suRandSec := 3

	if suDuration >= 5 {
		suRandSec = suDuration/2 + 1
	}

	suDuration = rand.Intn(suRandSec) + suDuration
	serviceUpdate := &timerassistant.CallInfo{
		Category: &timerassistant.Interval{
			Duration: time.Duration(suDuration) * time.Second,
		},
		Fn:           s.serviceUpdateTimer,
		ResumeCallCh: nil,
	}
	s.timer.AddCallBack(serviceUpdate)
	disconnectedCheck := &timerassistant.CallInfo{
		Category: &timerassistant.Interval{
			Duration: 1 * time.Second,
		},
		Fn:           client.GetMe().DisconnectedCheck,
		ResumeCallCh: nil,
	}
	s.timer.AddCallBack(disconnectedCheck)
	gmOnTimer := &timerassistant.CallInfo{
		Category: &timerassistant.Interval{
			Duration: 1 * time.Second,
		},
		Fn:           gm.GMInstance.OnTimer,
		ResumeCallCh: nil,
	}
	s.timer.AddCallBack(gmOnTimer)
}

func (s *Server) Run() {
	go s.tcpServer.Run()
	go s.innerServer.Run()
	startHTTPServer(s.httpPort, s.httpHandler, s.Config.HTTP.TLSCertFile, s.Config.HTTP.TLSKeyFile)
	s.serviceRegister()

	go func() {
		tick := time.NewTicker(time.Second * 1)
		defer func() {
			tick.Stop()
			if err := recover(); err != nil {
				logger.Error("[Run] 全局定时器出错", err, "\n", string(debug.Stack()))
			}
		}()

		for {
			select {
			case <-tick.C:
				s.timer.Loop()
			}
		}
	}()
}

func (s *Server) serviceRegister() {

	s.httpSvcID = idgenerator.FetchServiceID(httpService, s.httpAddr)
	err := consul.RegisterService(s.httpSvcID, httpService, s.httpAddr, s.httpAddr, s.runPath, s.processIdx, s.Config.Global.ZoneId)
	if err != nil {
		logger.Error("[serviceRegister] %v service error : %v", httpService, err)
	}

	s.tcpSvcID = idgenerator.FetchServiceID(tcpService, s.tcpAddr)
	err = consul.RegisterServiceEx(s.tcpSvcID, tcpService, s.tcpAddr, s.httpAddr, s.runPath, s.processIdx, s.Config.Global.ZoneId,
		s.Config.Server.PrivateIP)
	if err != nil {
		logger.Error("serviceRegister] %v service error : %v", tcpService, err)
	}

	s.innerSvcID = idgenerator.FetchServiceID(innerTCPService, s.innerSvcAddr)
	err = consul.RegisterService(s.innerSvcID, innerTCPService, s.innerSvcAddr, s.httpAddr, s.runPath, s.processIdx, s.Config.Global.ZoneId)
	if err != nil {
		logger.Error("serviceRegister] %v service error : %v", innerTCPService, err)
	}

}

func (s *Server) serviceUpdateTimer() {
	clientNum := client.GetMe().GetClientNum()
	if clientNum != s.lastUpdateCount {
		err := consul.UpdateService(s.tcpSvcID,
			tcpService, s.tcpAddr, s.httpAddr,
			fn.GetPerformanceExt(s.tcpSvcID, s.httpAddr, s.tcpPort, s.runPath, s.startTM, int32(clientNum), s.processIdx,
				s.Config.Global.ZoneId, s.Config.Server.PrivateIP),
			s.runPath)
		if err != nil {
			logger.Error("[serviceUpdateTimer] register gateway-tcp service error : ", err)
		}
		s.lastUpdateCount = clientNum
	}
}

func (s *Server) serviceDeRegister() {

	err := consul.DeregisterService(s.httpSvcID)
	if err != nil {
		logger.Error("[serviceDeRegister] http service error %v", err)
	}

	err = consul.DeregisterService(s.tcpSvcID)
	if err != nil {
		logger.Error("[serviceDeRegister] tcp service error %v", err)
	}

	err = consul.DeregisterService(s.innerSvcID)
	if err != nil {
		logger.Error("[serviceDeRegister] inner tcp service error %v", err)
	}
}

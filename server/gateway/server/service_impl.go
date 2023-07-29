package server

import (
	"greatestworks/aop/fn"
	"greatestworks/aop/logger"
	"greatestworks/server/gateway/client"
	"greatestworks/server/gateway/config"
	"greatestworks/server/gateway/gm"
	"greatestworks/server/gateway/world"
	"strconv"
	"time"
)

func (s *Server) Start() {
	s.Run()
}

func (s *Server) Init(cfg interface{}, processIdx int) {

	configInstance, ok := cfg.(*config.Config)
	if !ok {
		logger.Error("[Init] load config error init !!!")
	}
	s.Config = configInstance

	s.runPath = fn.GetCurrentDirectory()

	s.processIdx = processIdx

	sizeNums := fn.SplitStringToUint32Slice(configInstance.Server.PriMsgBuffSize, "*")
	if len(sizeNums) == 0 {
		s.PriMsgBuffSize = 16 * 1024
	} else {
		var buffSize int64 = 1
		for _, num := range sizeNums {
			buffSize = buffSize * int64(num)
		}
		s.PriMsgBuffSize = buffSize
	}

	sizeNums = fn.SplitStringToUint32Slice(configInstance.Server.PriConnBuffSize, "*")
	if len(sizeNums) == 0 {
		s.PriConnBuffSize = 128 * 1024
	} else {
		var buffSize int64 = 1
		for _, num := range sizeNums {
			buffSize = buffSize * int64(num)
		}
		s.PriConnBuffSize = buffSize
	}

	sizeNums = fn.SplitStringToUint32Slice(configInstance.Server.PubMsgBuffSize, "*")
	if len(sizeNums) == 0 {
		s.PubMsgBuffSize = 16 * 1024
	} else {
		var buffSize int64 = 1
		for _, num := range sizeNums {
			buffSize = buffSize * int64(num)
		}
		s.PubMsgBuffSize = buffSize
	}

	sizeNums = fn.SplitStringToUint32Slice(configInstance.Server.PubConnBuffSize, "*")
	if len(sizeNums) == 0 {
		s.PubConnBuffSize = 128 * 1024
	} else {
		var buffSize int64 = 1
		for _, num := range sizeNums {
			buffSize = buffSize * int64(num)
		}
		s.PubConnBuffSize = buffSize
	}

	s.tcpPort = s.Config.Server.Port + processIdx
	addr := s.Config.Server.PrivateIP + ":" + strconv.Itoa(s.tcpPort)
	s.tcpAddr = s.Config.Server.PublicIP + ":" + strconv.Itoa(s.tcpPort)
	s.tcpServer.Addr = addr
	s.tcpServer.MaxConnNum = s.Config.Server.MaxConnNum
	s.tcpServer.Init()
	s.tcpServer.NewSession = client.NewSession

	s.innerPort = s.Config.Server.InnerPort + processIdx
	innerAddr := s.Config.Server.PrivateIP + ":" + strconv.Itoa(s.innerPort)
	s.innerSvcAddr = innerAddr
	s.innerServer.Addr = innerAddr
	s.innerServer.MaxConnNum = s.Config.Server.MaxConnNum
	s.innerServer.Init()
	s.innerServer.NewSession = world.NewSession
	s.startTM = time.Now().Unix()

	s.httpPort = s.Config.HTTP.HTTPPort + processIdx
	s.httpAddr = s.Config.HTTP.HTTPAddr + ":" + strconv.Itoa(s.httpPort)
	s.id = strconv.FormatUint(fn.IpAddressStringToUint64(s.httpAddr), 10)
	s.Name = configInstance.NodeName
	s.DeploymentId = configInstance.DeploymentId
	gm.Init(s.startTM, s.id, "gateway", 10, s.Config.Global.IsOpenNow)
	s.registerTimer()
	s.clearConnRecords()
}

func (s *Server) Stop() {
	s.tcpServer.Close()
	s.innerServer.Close()
	s.serviceDeRegister()
}

func (s *Server) Reload() {
	logger.Info("[Reload] 服务器收到热更新信号...")
}

package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"greatestworks/aop/logger"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	consulClient *Client
)

func GetConsul() *Client {
	return consulClient
}

func InitConsul(conf *Config) error {
	config := api.DefaultConfig()
	addrs := conf.Nodes
	if len(addrs) <= 0 {
		return errors.New("[InitConsul] addrs length is zero")
	}

	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	config.Address = addrs[randSeed.Intn(len(conf.Nodes))]
	config.Token = conf.Token
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	consulClient = &Client{
		client:   client,
		services: &sync.Map{},
		addrs:    addrs,
	}
	return nil
}

func RegisterService(svcId, svcName, svcAddr, checkAddr, runPath string, processIdx int, zoneId int) error {
	var userName string
	if runtime.GOOS == "windows" {
		userName = "win" + GetOsUserName()
	} else {
		userName = GetOsUserName()
	}
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	interval := randSeed.Intn(5) + 10
	healthCheck := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s/health", checkAddr),
		Timeout:                        "3s",
		Interval:                       fmt.Sprintf("%vs", interval),
		Status:                         "passing",
		DeregisterCriticalServiceAfter: "1m",
	}

	host, port, err := net.SplitHostPort(svcAddr)
	if err != nil {
		logger.Logger.ErrorF("[RegisterService] error: %v", err.Error())
		return err
	}
	r := new(api.AgentServiceRegistration)
	r.ID = svcId
	r.Name = fmt.Sprintf("%s-%s", userName, svcName)
	r.Address = host
	r.Port, _ = strconv.Atoi(port)
	r.EnableTagOverride = true
	performance := &Performance{
		Zid:       zoneId,
		PIdx:      processIdx,
		SvrPort:   r.Port,
		StartTM:   time.Now().Unix(),
		SvrAddr:   svcAddr,
		SvrId:     svcId,
		SvrPath:   runPath,
		InnerAddr: "",
	}
	r.Tags = []string{GetPerformance(performance)}
	r.Check = healthCheck

	return consulClient.register(r)
}

func RegisterServiceEx(svcID, svcName, svcAddr, checkAddr, runPath string, processIdx int, zoneId int, innerIp string) error {
	var userName string
	if runtime.GOOS == "windows" {
		userName = "win" + GetOsUserName()
	} else {
		userName = GetOsUserName()
	}
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	interval := randSeed.Intn(5) + 10
	healthCheck := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s/health", checkAddr),
		Timeout:                        "3s",
		Interval:                       fmt.Sprintf("%vs", interval),
		Status:                         "passing",
		DeregisterCriticalServiceAfter: "1m",
	}

	host, port, err := net.SplitHostPort(svcAddr)
	if err != nil {
		logger.Logger.ErrorF("[RegisterServiceEx] error:%v ", err.Error())
		return err
	}
	r := new(api.AgentServiceRegistration)
	r.ID = svcID
	r.Name = fmt.Sprintf("%s-%s", userName, svcName)
	r.Address = host
	r.Port, _ = strconv.Atoi(port)
	r.EnableTagOverride = true
	performance := &Performance{
		Zid:       zoneId,
		PIdx:      processIdx,
		SvrPort:   r.Port,
		StartTM:   time.Now().Unix(),
		SvrAddr:   svcAddr,
		SvrId:     svcID,
		SvrPath:   runPath,
		InnerAddr: innerIp,
	}
	r.Tags = []string{GetPerformance(performance)}
	r.Check = healthCheck

	return consulClient.register(r)
}

func UpdateService(svcID, svcName, svcAddr, checkAddr, weight, runPath string) error {
	var userName string
	if runtime.GOOS == "windows" {
		userName = "win" + GetOsUserName()
	} else {
		userName = GetOsUserName()
	}
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	interval := randSeed.Intn(5) + 10
	healthCheck := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s/health", checkAddr),
		Timeout:                        "3s",
		Interval:                       fmt.Sprintf("%vs", interval),
		Status:                         "passing",
		DeregisterCriticalServiceAfter: "1m",
	}

	host, port, err := net.SplitHostPort(svcAddr)
	if err != nil {
		logger.Logger.ErrorF("[UpdateService]  error:%v", err.Error())
		return err
	}
	r := new(api.AgentServiceRegistration)
	r.ID = svcID
	r.Name = fmt.Sprintf("%s-%s", userName, svcName)
	r.Address = host
	r.Port, _ = strconv.Atoi(port)
	r.EnableTagOverride = true
	r.Tags = []string{weight}
	r.Check = healthCheck

	return consulClient.register(r)
}

func DeregisterService(svcID string) error {
	return consulClient.deregister(svcID)
}

func QueryServices(service string) ([]*api.ServiceEntry, *api.QueryMeta, error) {
	var userName string
	if runtime.GOOS == "windows" {
		userName = "win" + GetOsUserName()
	} else {
		userName = GetOsUserName()
	}
	return consulClient.service(userName+"-"+service, "", false, nil)
}

package main

import (
	"context"
	"sync"
	"time"
)

var gm GM

type GM struct {
	fetchIntervalSec int64
	lastFetchTime    int64
	ipWhiteList      *sync.Map
}

func (g *GM) Init(intervalSec int64) {
	g.fetchIntervalSec = intervalSec
	g.ipWhiteList = &sync.Map{}
}

func (g *GM) OnTimer() {
	now := time.Now().Unix()
	if now-g.lastFetchTime >= g.fetchIntervalSec {
		g.lastFetchTime = now
		g.fetchInfoFromRedis(now)
	}
}

func (g *GM) fetchInfoFromRedis(now int64) {
	g.fetchIpWhiteList(now)
}

func (g *GM) fetchIpWhiteList(now int64) {
	ipList := redisCluster.HGetAll(context.TODO(), "GatewayIpWhiteList").Val()
	for k, v := range ipList {
		if v == "1" {
			g.ipWhiteList.Store(k, true)
		} else {
			g.ipWhiteList.Store(k, false)
		}
	}
}

func (g *GM) IsIpInWhiteList(ip string) bool {
	ret, ok := g.ipWhiteList.Load(ip)
	return ok && ret.(bool)
}

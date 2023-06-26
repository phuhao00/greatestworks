package server

import (
	"greatestworks/server/gateway/client"
	"greatestworks/server/gateway/config"
	"time"
)

type Gateway struct {
	clientServer client.Session
	serverServer Server
	Metrics      *MetricInfo
	Config       config.Config
}

type MetricInfo struct {
	ClientCount int32 `json:"client_count"`
	ServerCount int32 `json:"server_count"`
}

func (g *Gateway) Update() {
	//g.Metrics.ServerCount = g.serverServer.GetServerCount()
	//g.Metrics.ClientCount = g.clientServer.GetServerCount()
}

func (g *Gateway) Loop() {

	tick := time.NewTicker(time.Duration(g.Config.UpdateInfoCd) * time.Second)

	for {
		select {
		case <-tick.C:
			g.Update()
		}
	}
}

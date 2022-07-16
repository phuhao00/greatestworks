package world

import (
	"greatestworks/network"
	"greatestworks/player"
)

func (mm *MgrMgr) UserLogin(message *network.SessionPacket) {
	newPlayer := player.NewPlayer()
	newPlayer.UId = 111
	newPlayer.HandlerParamCh = message.Sess.WriteCh
	message.Sess.IsPlayerOnline = true
	newPlayer.Run()
}

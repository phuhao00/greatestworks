package player

import (
	"github.com/phuhao00/greatestworks-proto/gm"
	"greatestworks/aop/fn"
	"greatestworks/aop/logger"
	"greatestworks/server/world/server"
)

func (p *Player) playerGMHandler(msgId uint16, data []byte) {
	msgReceive := &gm.CSPlayerGMCmd{}
	err := p.LogicRouter.Unmarshal(data, msgReceive)
	if err != nil {
		logger.Error("receive data:%v msg:%v", data, msgReceive)
		return
	}
	params := fn.SplitStringToUint32Slice(msgReceive.ParamStr, ",")

	if server.Oasis.Config.Settings == nil || !server.Oasis.Config.Settings.GMCommand {
		logger.Error("[playerGMHandler] server close PlayerID:%v op:%v %v %v %v", p.PlayerID, msgReceive.Op,
			msgReceive.Params[0], msgReceive.Params[1], msgReceive.Params[2])
		return
	}
	logger.Debug("[playerGMHandler] execute PlayerID:%v data:%v", p.PlayerID, msgReceive)
	switch msgReceive.GetOp() {
	case 0:
		_ = params[0]
	}
}

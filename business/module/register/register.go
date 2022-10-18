package register

import (
	"github.com/phuhao00/network"
	"greatestworks/business/module/player"
)

type Fn func(player *player.Player, packet *network.Packet)

func Register(cmd uint32, fn Fn) {
	//装饰器
	
	//
}

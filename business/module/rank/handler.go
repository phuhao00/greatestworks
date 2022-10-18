package rank

import (
	"github.com/phuhao00/network"
	"greatestworks/business/module/player"
	"greatestworks/business/module/register"
)

func init() {
	register.Register(222, GetRankList)

}

func GetRankList(player *player.Player, packet *network.Packet) {

}
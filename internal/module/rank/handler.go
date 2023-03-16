package rank

import (
	"github.com/phuhao00/network"
	"greatestworks/business/module/register"
	"greatestworks/internal/module/player"
)

func init() {
	register.Register(222, GetRankList)

}

func GetRankList(player *player.Player, packet *network.Packet) {

}
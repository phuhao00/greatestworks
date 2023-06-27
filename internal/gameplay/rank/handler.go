package rank

import (
	"github.com/phuhao00/network"
	"greatestworks/aop/module_router"
	"greatestworks/internal/communicate/player"
)

func init() {

}

func GetRankList(player *player.Player, packet *network.Packet) {

}

func RegisterHandler() {
	module_router.RegisterModuleMessageHandler(0, 0, nil)
}

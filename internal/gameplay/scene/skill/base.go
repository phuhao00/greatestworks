package skill

import (
	"greatestworks/internal/gameplay/scene/buff"
)

type Base struct {
	Id     uint32
	Desc   string
	Cd     int64
	Damage int64
	Buffs  []buff.Abstract
}

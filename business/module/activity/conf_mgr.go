package activity

import (
	"greatestworks/business/module/base"
	"sync"
)

type ConfigManager struct {
	base.ConfigManagerBase
	Configs sync.Map //策划配置

}

func (m *ConfigManager) Get(id uint32) interface{} {
	var ret any
	m.Configs.Range(func(key, value any) bool {
		idAssert := key.(uint32)
		if idAssert == id {
			ret = value
			return false
		}
		return true
	})
	return ret
}

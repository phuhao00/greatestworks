package excel_object

import "sync"

type Manager struct {
	SkillConfigs sync.Map //技能表
	DropsConfigs sync.Map //掉落
}

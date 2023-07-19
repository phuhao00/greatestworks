package task_category_group

import (
	"greatestworks/internal/gameplay/task"
	"sync"
)

type Achievement struct {
	Records   sync.Map //map[string]int //历史计数，用于完成成就判断
	Completed []uint32 //完成的
	Submitted []uint32 //提交的
	Tasks     []task.BaseTask
}

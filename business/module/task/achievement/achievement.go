package achievement

import (
	"greatestworks/business/module/task"
	"sync"
)

// 成就
type Data struct {
	Records   sync.Map //map[string]int //历史计数，用于完成成就判断
	Completed []uint32 //完成的
	Submitted []uint32 //提交的
}

func (d *Data) AddRecords(group task.TargetCategory, subGroup int, delta int) {

}

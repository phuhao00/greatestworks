package idgenerator

import (
	"github.com/sony/sonyflake"
	"greatestworks/aop/idgenerator/local"
)

var (
	sf          *sonyflake.Sonyflake
	node2NodeId = map[string]uint16{
		"login":   1,
		"gateway": 2,
		"world":   3,
		"game":    4,
		"gm":      5,
		"upload":  6,
		"robot":   7,
		"unknown": 8,
	}
)

func init() {
	var st sonyflake.Settings
	st.MachineID = local.MachineID //设置为nil 默认拿本机ip
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func GenerateId() map[string]uint64 {
	id, err := sf.NextID()
	if err != nil {
		return nil
	}
	return sonyflake.Decompose(id)
}

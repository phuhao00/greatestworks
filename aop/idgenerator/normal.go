package idgenerator

import (
	"github.com/sony/sonyflake"
	"greatestworks/aop/idgenerator/local"
)

var sf *sonyflake.Sonyflake

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

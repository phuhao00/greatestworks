package condition

import "sync"

type Register struct {
	conditions map[uint32]func(str string) Condition //Id->condition
}

var (
	ConRegs *Register
	onceDo  sync.Once
)

func GetMe() *Register {
	return ConRegs
}

func (r *Register) OnInit() {
	onceDo.Do(func() {
		ConRegs = &Register{conditions: make(map[uint32]func(str string) Condition)}
	})
}

func (r *Register) Reg(id uint32, condition func(str string) Condition) {
	r.conditions[id] = condition
}

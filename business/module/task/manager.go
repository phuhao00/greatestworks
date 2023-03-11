package task

import (
	"greatestworks/aop/event"
	"sync"
)

var (
	manager        *Manager
	newManagerOnce sync.Once
)

type Manager struct {
	configs      sync.Map
	ChIn         chan *PlayerActionParam
	ChOut        chan interface{}
	ChEvent      chan *EventWrap
	LoopNum      int
	MonitorNum   int
	events       sync.Map
	eventHandles map[event.IEvent]EventHandle
}

func NewManager(conf *ManagerConfig) {

	var (
		loopNum    int
		monitorNum int
		chOutSize  int
		chInSize   int
	)
	if conf.LoopNum == 0 {
		loopNum = defaultLoopNum
	}
	if conf.MonitorNum == 0 {
		monitorNum = defaultMonitor
	}
	if conf.ChInSize == 0 {
		chInSize = defaultChanInSize
	}
	if conf.ChOutSize == 0 {
		chOutSize = defaultChanOutSize
	}

	manager = &Manager{
		configs:    sync.Map{},
		ChIn:       make(chan *PlayerActionParam, chInSize),
		ChOut:      make(chan interface{}, chOutSize),
		LoopNum:    loopNum,
		MonitorNum: monitorNum,
	}
}

func GetMe() *Manager {

	return manager
}

func (m *Manager) Run() {

	for i := 0; i < m.LoopNum; i++ {
		go m.Loop()
	}

	for i := 0; i < m.MonitorNum; i++ {
		go m.Monitor()
	}
}

func (m *Manager) Loop() {
	for {
		select {
		case <-m.ChOut:

		}
	}
}

func (m *Manager) Monitor() {
	for {
		select {
		case p := <-m.ChIn:
			m.Handle(p)
		case e := <-m.ChEvent:
			m.OnEvent(e)
		}
	}
}

func (m *Manager) Handle(param *PlayerActionParam) {
	handler, err := GetHandler(param.MessageId)
	if err != nil {
		//todo log
	}
	handler.Fn(param.Player, param.Packet)
}

// getTaskConfig get task config
func (m *Manager) getTaskConfig(confId uint32) (ret *Config) {
	m.configs.Range(func(key, value any) bool {
		if val, ok := value.(*Config); ok {
			if val.Id == confId {
				ret = val
				return false
			}
		}
		return true
	})
	return ret
}

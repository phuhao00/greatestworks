package task

import "sync"

var (
	manager        *Manager
	newManagerOnce sync.Once
)

type Manager struct {
	configs    sync.Map
	ChIn       chan *PlayerActionParam
	ChOut      chan interface{}
	ChEvent    chan *EventParam
	LoopNum    int
	MonitorNum int
	events     sync.Map
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
			m.HandleEvent(e)
		}
	}
}

func (m *Manager) HandleEvent(p *EventParam) {
	e := m.getEvent(p.EventCategory)
	e.Notify(p.Param, p.Player)
}

func (m *Manager) Handle(param *PlayerActionParam) {
	handler, err := GetHandler(param.MessageId)
	if err != nil {
		//todo log
	}
	handler.Fn(param.Player, param.Packet)
}

// getTaskConfig get task config
func (m *Manager) getTaskConfig(taskConfigId uint32) (ret *Config) {
	m.configs.Range(func(key, value any) bool {
		if val, ok := value.(*Config); ok {
			if val.Id == taskConfigId {
				ret = val
				return false
			}
		}
		return true
	})
	return ret
}

func (m *Manager) getEvent(e EventCategory) Event {
	value, ok := m.events.Load(e)
	if ok {
		return value.(Event)
	}
	return nil
}

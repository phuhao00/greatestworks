package family

type Manager struct {
	families map[uint64]*Family
	Owner
	ChIn  chan ManagerHandlerParam
	ChOut chan interface{}
}

func (m *Manager) Loop() {
	for {
		select {
		case msg := <-m.ChOut:
			m.ForwardMsg(msg)
		}
	}
}

func (m *Manager) Monitor() {
	for {
		select {
		case param := <-m.ChIn:
			m.Handler(param)
		}
	}
}

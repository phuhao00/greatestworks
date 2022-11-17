package family

type Manager struct {
	families map[uint32]*Family
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

func (m *Manager) Handler(param interface{}) {

}

func (m *Manager) ForwardMsg(msg interface{}) {
	m.Owner.BroadcastMsg(nil, 0, nil)

}

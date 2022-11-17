package family

type Manager struct {
	families map[uint32]*Family
	Owner
	HandlerParam chan interface{}
	broadcastMsg chan interface{}
}

func (m *Manager) Loop() {
	for {
		select {
		case msg := <-m.broadcastMsg:
			m.ForwardMsg(msg)
		}
	}
}

func (m *Manager) Monitor() {
	for {
		select {
		case param := <-m.HandlerParam:
			m.Handler(param)
		}
	}
}

func (m *Manager) Handler(param interface{}) {

}

func (m *Manager) ForwardMsg(msg interface{}) {

}

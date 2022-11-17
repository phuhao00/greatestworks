package family

func (m *Manager) GetFamily(id uint64) *Family {
	return m.families[id]
}

func (m *Manager) Handler(param ManagerHandlerParam) {
	if family := m.GetFamily(param.FamilyId); family != nil {
		family.ChIn <- param.MemberActionParam
	}
}

func (m *Manager) ForwardMsg(msg interface{}) {
	m.Owner.BroadcastMsg(nil, 0, nil)

}

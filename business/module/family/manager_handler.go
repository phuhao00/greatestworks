package family

func (m *Module) GetFamily(id uint64) *Family {
	return m.families[id]
}

func (m *Module) Handler(param ManagerHandlerParam) {
	if family := m.GetFamily(param.FamilyId); family != nil {
		family.ChIn <- param.MemberActionParam
	}
}

func (m *Module) ForwardMsg(msg interface{}) {
	m.IWorld.BroadcastMsg(nil, 0, nil)

}

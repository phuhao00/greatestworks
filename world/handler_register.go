package world

func (mm *MgrMgr) HandlerRegister() {
	mm.Handlers[1] = mm.UserLogin
}

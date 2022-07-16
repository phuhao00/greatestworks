package player

func (p *Player) HandlerRegister() {
	p.handlers[111] = p.AddFriend
	p.handlers[222] = p.DelFriend
	p.handlers[333] = p.ResolveChatMsg
}

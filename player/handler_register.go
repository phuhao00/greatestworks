package player

func (p *Player) HandlerRegister() {
	p.handlers["add_friend"] = p.AddFriend
	p.handlers["del_friend"] = p.DelFriend
	p.handlers["del_friend"] = p.ResolveChatMsg
}

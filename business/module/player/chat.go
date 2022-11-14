package player

import "github.com/nsqio/go-nsq"

func (p *Player) InitNsqHandler(channel string) {
	//TODO implement me
	panic("implement me")
}

func (p *Player) HandleMessage(message nsq.Message) error {
	//TODO implement me
	panic("implement me")
}

func (p *Player) PublishChatMsg(chatMsg interface{}) error {
	//TODO implement me
	panic("implement me")
}

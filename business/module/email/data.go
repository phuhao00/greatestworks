package email

import "sync"

type MailM struct {
	Id     uint64 `bson:"id"`
	ConfId uint32 `bson:"confId"`
	Status int    `json:"status"`
	Desc   string `json:"desc"`
}

type MailStatus int

const (
	MailStatusUnRead MailStatus = iota + 1
	MailStatusRead
	MailStatusDelete
)

type Data struct {
	mails sync.Map
}

type Owner interface {
	AddMail(mail IEmail)
	ReadMail()
	DelMail()
	DelAllMail()
	GetMailItem()
	GetAllMailItems()
}

func (o *Data) AddMail(mail IEmail) {
	o.mails.Store(mail.GetID, mail)
}

func (o *Data) ReadMail() {
	//TODO implement me
	panic("implement me")
}

func (o *Data) DelMail() {
	//TODO implement me
	panic("implement me")
}

func (o *Data) DelAllMail() {
	//TODO implement me
	panic("implement me")
}

func (o *Data) GetMailItem() {
	//TODO implement me
	panic("implement me")
}

func (o *Data) GetAllMailItems() {
	//TODO implement me
	panic("implement me")
}

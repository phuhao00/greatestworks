package email

import "sync"

type Owner interface {
	AddMail(mail Abstract)
	ReadMail()
	DelMail()
	DelAllMail()
	GetMailItem()
	GetAllMailItems()
}

type OwnerBase struct {
	mails sync.Map
}

func (o *OwnerBase) AddMail(mail Abstract) {
	o.mails.Store(mail.GetID, mail)
}

func (o *OwnerBase) ReadMail() {
	//TODO implement me
	panic("implement me")
}

func (o *OwnerBase) DelMail() {
	//TODO implement me
	panic("implement me")
}

func (o *OwnerBase) DelAllMail() {
	//TODO implement me
	panic("implement me")
}

func (o *OwnerBase) GetMailItem() {
	//TODO implement me
	panic("implement me")
}

func (o *OwnerBase) GetAllMailItems() {
	//TODO implement me
	panic("implement me")
}

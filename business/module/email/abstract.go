package email

type Abstract interface {
	SetStatus(status MailStatus)
	GetID() uint64
	ToPB()
	LoadFrom(*MailM)
	GetDBModel() *MailM
}

package email

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

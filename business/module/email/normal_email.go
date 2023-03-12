package email

type NormalMail struct {
	Id     uint64
	ConfId uint32
	Desc   string
	Status MailStatus
}

func (m *NormalMail) SetStatus(status MailStatus) {
	m.Status = status
}

func (m *NormalMail) GetID() uint64 {
	return m.Id
}
func (m *NormalMail) ToPB() {
	//TODO implement me
	panic("implement me")
}

func (m *NormalMail) LoadFrom(data *MailM) {
	m.Status = MailStatus(data.Status)
	m.Desc = data.Desc
	m.ConfId = data.ConfId
	m.Id = data.Id
}

func (m *NormalMail) GetDBModel() *MailM {
	return &MailM{
		Id:     m.Id,
		ConfId: m.ConfId,
		Status: int(m.Status),
		Desc:   m.Desc,
	}
}

package family

type Family struct {
	Id       uint64
	Name     string
	Desc     string
	members  map[uint64]*Member
	requests []interface{}
	ChIn     chan *MemberActionParam
	ChOut    chan interface{}
}

func (f *Family) Start() {
	//TODO implement me
	panic("implement me")
}

func (f *Family) Stop() {
	//TODO implement me
	panic("implement me")
}

func (f *Family) AddMember(member *Member) {
	f.members[member.Id] = member
}

func (f *Family) DelMember(ids []uint64) {
	for _, id := range ids {
		delete(f.members, id)
	}
}

func (f *Family) UpdateMember(member *Member) {
	f.members[member.Id] = member
}

func (f *Family) GetMember(id uint64) *Member {
	return f.members[id]
}

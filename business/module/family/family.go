package family

type Family struct {
	members map[uint64]*Member
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

func (f *Family) DelMember(id uint64) {
	delete(f.members, id)
}

func (f *Family) UpdateMember(member *Member) {
	f.members[member.Id] = member
}

func (f *Family) GetMember(id uint64) *Member {
	return f.members[id]
}

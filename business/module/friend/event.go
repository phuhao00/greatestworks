package friend

type AddOrDelFriendEvent struct {
	CurFriendCount int
}

func (e AddOrDelFriendEvent) GetDesc() string {
	return ""
}

func (s *System) PublishAddOrDelFriend() {
	e := &AddOrDelFriendEvent{}
	s.DataAsPublisher.Publish(e)
}

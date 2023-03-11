package friend

func (s *System) AddSubscriber() {
	s.DataAsPublisher.AddSubscriber(&AddOrDelFriendEvent{}, s.Player)
}

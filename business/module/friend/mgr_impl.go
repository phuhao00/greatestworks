package friend

import "time"

type System struct {
	friends   []Info
	BlackList []uint64
	requests  []Request
	Owner
}

func (s *System) SetOwner(owner Owner) {
	s.Owner = owner
}

func (s *System) Add(uId uint64) {

}

func (s *System) isFriend(uId uint64) (bool, int) {
	for index, val := range s.friends {
		if val.UId == uId {
			return true, index
		}
	}
	return false, -1
}

func (s *System) isBlackList(uId uint64) (bool, int) {
	for index, val := range s.BlackList {
		if val == uId {
			return true, index
		}
	}
	return false, -1
}

func (s *System) getRequest(uId uint64) (bool, int) {
	for index, val := range s.requests {
		if val.Userid == uId {
			return true, index
		}
	}
	return false, -1
}

func (s *System) delRequest(uId uint64) {
	if ok, index := s.getRequest(uId); ok == true {
		s.requests = append(s.requests[:index], s.requests[index+1:]...)
	}
}

func (s *System) addRequest(uId uint64, addType int32) {
	s.requests = append(s.requests, Request{
		Userid:  uId,
		OpTime:  time.Now().Unix(),
		AddType: addType,
	})
}

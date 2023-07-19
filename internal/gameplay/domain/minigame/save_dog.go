package minigame

import "time"

type SaveDog struct {
	players   []uint64
	startTime int64
	stopTime  int64
	IsSuccess bool
}

func NewSaveDog() *SaveDog {
	return &SaveDog{}
}

func (s *SaveDog) Start() {
	s.startTime = time.Now().Unix()
	s.Run()
}

func (s *SaveDog) Stop() {
	s.stopTime = time.Now().Unix()
	s.End()
}

func (s *SaveDog) End() {

}

func (s *SaveDog) Run() {

}

func (s *SaveDog) CheckSuccess() bool {
	return true
}

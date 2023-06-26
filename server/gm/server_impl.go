package main

func (s *Server) Start() {
}

func (s *Server) Reload() {
	for {
		select {
		case data := <-s.toGateWay:
			s.OnGateWayMsg(data)
		}
	}
}

func (s *Server) Init(config interface{}, processId int) {

}

func (s *Server) Stop() {

}

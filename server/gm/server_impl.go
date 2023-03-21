package main

func (s *Server) Start() {
}

func (s *Server) Loop() {
	for {
		select {
		case data := <-s.toGateWay:
			s.OnGateWayMsg(data)
		}
	}
}

func (s *Server) Monitor() {
}

func (s *Server) Stop() {
}

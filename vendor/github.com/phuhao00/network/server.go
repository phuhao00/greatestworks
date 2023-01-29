package network

import (
	"github.com/phuhao00/spoor"
	"net"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	pid            int64
	Addr           string
	MaxConnNum     int
	ln             *net.TCPListener
	connSet        map[net.Conn]interface{}
	counter        int64
	idCounter      int64
	mutexConn      sync.Mutex
	wgLn           sync.WaitGroup
	wgConn         sync.WaitGroup
	connBuffSize   int
	logger         *spoor.Spoor
	MessageHandler func(packet *Packet)
}

func NewServer(addr string, maxConnNum int, buffSize int, logger *spoor.Spoor) *Server {
	s := &Server{
		Addr:         addr,
		MaxConnNum:   maxConnNum,
		connBuffSize: buffSize,
		logger:       logger,
	}
	s.Init()
	return s
}

func (s *Server) Init() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.Addr)

	if err != nil {
		s.logger.FatalF("[net] addr resolve error", tcpAddr, err)
	}

	ln, err := net.ListenTCP("tcp6", tcpAddr)

	if err != nil {
		s.logger.FatalF("%v", err)
	}

	if s.MaxConnNum <= 0 {
		s.MaxConnNum = 100
		s.logger.InfoF("invalid MaxConnNum, reset to %v", s.MaxConnNum)
	}

	s.ln = ln
	s.connSet = make(map[net.Conn]interface{})
	s.counter = 1
	s.idCounter = 1
	s.pid = int64(os.Getpid())
	s.logger.InfoF("Server Listen %s", s.ln.Addr().String())
}

func (s *Server) Run() {
	defer func() {
		if err := recover(); err != nil {
			s.logger.ErrorF("[net] panic", err, "\n", string(debug.Stack()))
		}
	}()

	s.wgLn.Add(1)
	defer s.wgLn.Done()

	var tempDelay time.Duration
	for {
		conn, err := s.ln.AcceptTCP()

		if err != nil {
			if _, ok := err.(net.Error); ok {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				s.logger.InfoF("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		if atomic.LoadInt64(&s.counter) >= int64(s.MaxConnNum) {
			conn.Close()
			s.logger.InfoF("too many connections %v", atomic.LoadInt64(&s.counter))
			continue
		}
		tcpConnX, err := NewTcpConnX(conn, s.connBuffSize, s.logger)
		if err != nil {
			s.logger.ErrorF("%v", err)
			return
		}
		s.addConn(conn, tcpConnX)
		tcpConnX.Impl = s
		s.wgConn.Add(1)
		go func() {
			tcpConnX.Connect()
			s.removeConn(conn, tcpConnX)
			s.wgConn.Done()
		}()
	}
}

func (s *Server) Close() {
	s.ln.Close()
	s.wgLn.Wait()

	s.mutexConn.Lock()
	for conn := range s.connSet {
		conn.Close()
	}
	s.connSet = nil
	s.mutexConn.Unlock()
	s.wgConn.Wait()
}

func (s *Server) addConn(conn net.Conn, tcpConnX *TcpConnX) {
	s.mutexConn.Lock()
	atomic.AddInt64(&s.counter, 1)
	s.connSet[conn] = conn
	nowTime := time.Now().Unix()
	idCounter := atomic.AddInt64(&s.idCounter, 1)
	connId := (nowTime << 32) | (s.pid << 24) | idCounter
	tcpConnX.ConnID = connId
	s.mutexConn.Unlock()
	tcpConnX.OnConnect()
}

func (s *Server) removeConn(conn net.Conn, tcpConn *TcpConnX) {
	tcpConn.Close()
	s.mutexConn.Lock()
	atomic.AddInt64(&s.counter, -1)
	delete(s.connSet, conn)
	s.mutexConn.Unlock()
}

func (s *Server) OnMessage(message *Message, conn *TcpConnX) {
	s.MessageHandler(&Packet{
		Msg:  message,
		Conn: conn,
	})
}

func (s *Server) OnClose() {

}

func (s *Server) OnConnect() {

}

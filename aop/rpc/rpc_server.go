package rpc

import (
	"greatestworks/aop/logger"
	"net"
	"net/rpc"
	"runtime/debug"
	"sync"
	"time"
)

type Server struct {
	Addr string
	ln   *net.TCPListener
	wgLn sync.WaitGroup
}

func (srv *Server) Init(addr string) {

	srv.Addr = addr

	tcpAddr, err := net.ResolveTCPAddr("tcp4", srv.Addr)

	if err != nil {
		logger.ErrorF("[net] addr resolve error", tcpAddr, err)
		return
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		logger.ErrorF("%v", err)
		return
	}

	srv.ln = ln
	logger.InfoF("RpcServer Listen %s", srv.ln.Addr().String())
}

func (srv *Server) Run() {
	// 捕获异常
	defer func() {
		if err := recover(); err != nil {
			logger.ErrorF("[net] panic", err, "\n", string(debug.Stack()))
		}
	}()

	srv.wgLn.Add(1)
	defer srv.wgLn.Done()

	var tempDelay time.Duration
	for {
		conn, err := srv.ln.AcceptTCP()

		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				logger.InfoF("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		// Try to open keepalive for tcp.
		conn.SetKeepAlive(true)

		conn.SetKeepAlivePeriod(1 * time.Minute)

		// disable Nagle algorithm.
		conn.SetNoDelay(true)

		conn.SetWriteBuffer(128 * 1024)

		conn.SetReadBuffer(128 * 1024)

		go rpc.ServeConn(conn)

		logger.DebugF("accept a rpc conn")
	}
}

func (srv *Server) Close() {
	srv.ln.Close()
	srv.wgLn.Wait()
}

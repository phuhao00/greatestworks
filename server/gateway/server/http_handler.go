package server

import (
	"crypto/tls"
	"fmt"
	"github.com/phuhao00/network"
	"greatestworks/aop/logger"
	"net"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"strings"
)

type HTTPHandler struct {
	Router *network.HttpRouter
}

func startHTTPServer(HTTPPort int, handler *HTTPHandler, TLSCertFile *string, TLSKeyFile *string) {
	handler.Register()
	handler.RegisterProfiler()
	srv := network.HttpServer(handler)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", HTTPPort))
	if err != nil {
		logger.Error("Listen http port fail %d", HTTPPort)
	}
	l = net.Listener(network.TCPKeepAliveListener{TCPListener: l.(*net.TCPListener)})

	if TLSCertFile != nil && TLSKeyFile != nil {
		cert, err := tls.LoadX509KeyPair(*TLSCertFile, *TLSKeyFile)
		if err == nil {
			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
				NextProtos:   []string{"http/1.1"},
			}
			l = tls.NewListener(l, srv.TLSConfig)
			logger.Info("HttpServer Use Https")
		} else {
			logger.Info("LoadX509KeyPair error %s, %s %v", *TLSCertFile, *TLSKeyFile, err)
		}
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("[异常] http服务器出错", err, "\n", string(debug.Stack()))
			}
		}()
		if err := srv.Serve(l); err != nil {
			logger.Error("encountered an error while serving listener: ", err)
		}
	}()
	logger.Info("HttpServer Listening on %s", l.Addr().String())
}

func (hs *HTTPHandler) Register() {
	hs.Router.HandleFunc("GET", "/health", HealthCheck)

}

func (hs *HTTPHandler) RegisterProfiler() {
	hs.Router.HandleFunc("GET", "/debug/pprof/", pprof.Index)
	hs.Router.HandleFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	hs.Router.HandleFunc("GET", "/debug/pprof/profile", pprof.Profile)
	hs.Router.HandleFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
	hs.Router.HandleFunc("GET", "/debug/pprof/trace", pprof.Trace)

	hs.Router.Handle("GET", "/debug/pprof/goroutine", pprof.Handler("goroutine"))
	hs.Router.Handle("GET", "/debug/pprof/heap", pprof.Handler("heap"))
	hs.Router.Handle("GET", "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	hs.Router.Handle("GET", "/debug/pprof/block", pprof.Handler("block"))

	hs.Router.HandleFunc("GET", "/debug/set_log_level", setLogLevel)
}

func (hs *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hs.Router.ServeHTTP(w, r)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {

	}
}

func setLogLevel(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	err := r.ParseForm()
	if err != nil {

	}
	level := strings.ToUpper(r.Form.Get("level"))
	tips := fmt.Sprintf("setLogLevel: %s", level)
	logger.Info(tips)
	if level == "DEBUG" {
	} else if level == "INFO" {
	} else {
		logger.Warn("The server log level could not be set as %v.", level)
		return
	}
	//logger.SetLogLevel(level)
	_, err = w.Write([]byte(tips))
	if err != nil {

	}
}

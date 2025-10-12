package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// Profiler wraps a standalone HTTP server exposing Go pprof diagnostics.
type Profiler struct {
	server *http.Server
	logger logging.Logger
	once   sync.Once
}

// NewProfiler creates a profiler instance using the provided logger.
func NewProfiler(logger logging.Logger) *Profiler {
	return &Profiler{logger: logger}
}

// RegisterHandlers attaches standard pprof handlers to the supplied mux.
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}

// Start launches a dedicated pprof HTTP server on the configured host/port.
// The server runs asynchronously and logs lifecycle events through the logger.
func (p *Profiler) Start(host string, port int) error {
	if p == nil {
		return fmt.Errorf("profiler is nil")
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	mux := http.NewServeMux()
	RegisterHandlers(mux)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	p.once.Do(func() {
		p.server = server

		go func() {
			p.logger.Info("pprof server starting", logging.Fields{"addr": addr})

			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				p.logger.Error("pprof server failed", err, logging.Fields{"addr": addr})
			}
		}()
	})

	return nil
}

// Stop gracefully shuts down the pprof HTTP server.
func (p *Profiler) Stop(ctx context.Context) error {
	if p == nil || p.server == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	p.logger.Info("stopping pprof server", logging.Fields{"addr": p.server.Addr})

	if err := p.server.Shutdown(shutdownCtx); err != nil {
		p.logger.Error("pprof server shutdown failed", err, logging.Fields{"addr": p.server.Addr})
		return err
	}

	p.logger.Info("pprof server stopped", logging.Fields{"addr": p.server.Addr})
	return nil
}

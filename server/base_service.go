package server

import (
	"context"
	"fmt"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
	"greatestworks/aop/colors"
	"greatestworks/aop/envelope/conn"
	"greatestworks/aop/files"
	"greatestworks/aop/logger"
	"greatestworks/aop/logging"
	"greatestworks/aop/logtype"
	imetrics "greatestworks/aop/metrics"
	"greatestworks/aop/net/call"
	"greatestworks/aop/perfetto"
	"greatestworks/aop/protos"
	"greatestworks/aop/retry"
	"greatestworks/aop/status"
	"greatestworks/aop/traceio"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

type BaseService struct {
	Id             string
	Name           string
	DeploymentId   string
	submissionTime time.Time
	statsProcessor *imetrics.StatsProcessor // tracks and computes stats to be rendered on the /statusz page.
	traceSaver     func(spans *protos.Spans) error
	Ctx            context.Context
	mu             sync.Mutex
	Inherit        IService
}

func NewBaseService(Name, DeploymentId string) (*BaseService, error) {
	ctx := context.Background()
	traceDB, err := perfetto.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot open Perfetto database: %w", err)
	}
	traceSaver := func(spans *protos.Spans) error {
		traces := make([]sdktrace.ReadOnlySpan, len(spans.Span))
		for i, span := range spans.Span {
			traces[i] = &traceio.ReadSpan{Span: span}
		}
		return traceDB.Store(ctx, Name, DeploymentId, traces)
	}
	bs := &BaseService{
		Ctx:            ctx,
		submissionTime: time.Now(),
		statsProcessor: imetrics.NewStatsProcessor(),
		traceSaver:     traceSaver,
	}
	go bs.statsProcessor.CollectMetrics(ctx, imetrics.Snapshot)
	return bs, nil
}

type stub struct {
	client   call.Connection  // client to talk to the remote component, created lazily.
	methods  []call.MethodKey // Keys for the remote component methods.
	balancer call.Balancer    // if not nil, component load balancer
	tracer   trace.Tracer     // component tracer

}

// serveStatus runs and registers the weaver-single status server.
func (e *BaseService) serverStatus(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.DefaultServeMux)
	status.RegisterServer(mux, e, e.SystemLogger())
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}
	errs := make(chan error, 1)
	go func() {
		errs <- serveHTTP(ctx, lis, mux)
	}()

	// Wait for the status server to become active.
	client := status.NewClient(lis.Addr().String())
	for r := retry.Begin(); r.Continue(ctx); {
		_, err := client.Status(ctx)
		if err == nil {
			break
		}
		e.SystemLogger().Error("status server unavailable", err, "address", lis.Addr())
	}

	// AddHandler the deployment.
	dir, err := files.DefaultDataDir()
	if err != nil {
		return err
	}
	dir = filepath.Join(dir, "single_registry")
	registry, err := status.NewRegistry(ctx, dir)
	if err != nil {
		return nil
	}
	reg := status.Registration{
		DeploymentId: e.DeploymentId,
		App:          e.Name,
		Addr:         lis.Addr().String(),
	}
	fmt.Fprint(os.Stderr, reg.Rolodex())
	if err := registry.Register(ctx, reg); err != nil {
		return err
	}

	// Unregister the deployment if this application is killed.
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		if err := registry.Unregister(ctx, reg.DeploymentId); err != nil {
			fmt.Fprintf(os.Stderr, "unregister deployment: %v\n", err)
		}
		os.Exit(1)
	}()

	return <-errs
}

// Status implements the status.Server interface.
func (e *BaseService) Status(ctx context.Context) (*status.Status, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	pid := int64(os.Getpid())
	stats := e.statsProcessor.GetStatsStatusz()
	components := []*status.Component{{Name: "main.go", Pids: []int64{pid}}}
	c := &status.Component{
		Name:  e.Name,
		Group: "main.go",
		Pids:  []int64{pid},
	}
	components = append(components, c)

	s := stats[logging.ShortenComponent(c.Name)]
	if s == nil {
		return nil, nil
	}
	for _, methodStats := range s {
		c.Methods = append(c.Methods, &status.Method{
			Name: methodStats.Name,
			Minute: &status.MethodStats{
				NumCalls:     methodStats.Minute.NumCalls,
				AvgLatencyMs: methodStats.Minute.AvgLatencyMs,
				RecvKbPerSec: methodStats.Minute.RecvKBPerSec,
				SentKbPerSec: methodStats.Minute.SentKBPerSec,
			},
			Hour: &status.MethodStats{
				NumCalls:     methodStats.Hour.NumCalls,
				AvgLatencyMs: methodStats.Hour.AvgLatencyMs,
				RecvKbPerSec: methodStats.Hour.RecvKBPerSec,
				SentKbPerSec: methodStats.Hour.SentKBPerSec,
			},
			Total: &status.MethodStats{
				NumCalls:     methodStats.Total.NumCalls,
				AvgLatencyMs: methodStats.Total.AvgLatencyMs,
				RecvKbPerSec: methodStats.Total.RecvKBPerSec,
				SentKbPerSec: methodStats.Total.SentKBPerSec,
			},
		})
	}

	return &status.Status{
		App:            e.Name,
		DeploymentId:   e.DeploymentId,
		SubmissionTime: timestamppb.New(e.submissionTime),
		Components:     components,
	}, nil
}

// Metrics implements the status.Server interface.
func (e *BaseService) Metrics(context.Context) (*status.Metrics, error) {
	m := &status.Metrics{}
	for _, snap := range imetrics.Snapshot() {
		proto := snap.ToProto()
		if proto.Labels == nil {
			proto.Labels = map[string]string{}
		}
		proto.Labels["server_name"] = e.Name
		proto.Labels["deploymentId"] = e.DeploymentId
		proto.Labels["node"] = e.Id
		m.Metrics = append(m.Metrics, proto)
	}
	return m, nil
}

// Profile implements the status.Server interface.
func (e *BaseService) Profile(_ context.Context, req *protos.RunProfiling) (*protos.Profile, error) {
	data, err := conn.Profile(req)
	profile := &protos.Profile{
		AppName:   e.Name,
		VersionId: e.DeploymentId,
		Data:      data,
	}
	if err != nil {
		profile.Errors = []string{err.Error()}
	}
	return profile, nil
}

func (e *BaseService) CreateLogSaver(_ context.Context, component string) func(entry *protos.LogEntry) {
	pp := logging.NewPrettyPrinter(colors.Enabled())
	return func(entry *protos.LogEntry) {
		fmt.Fprintln(os.Stderr, pp.Format(entry))
	}
}

func (e *BaseService) CreateTraceExporter() (sdktrace.SpanExporter, error) {
	return traceio.NewWriter(e.traceSaver), nil
}

func (e *BaseService) SystemLogger() logtype.Logger {
	return newAttrLogger(e.Name, e.DeploymentId, e.CreateLogSaver(e.Ctx, e.Name))
}

// serveHTTP serves HTTP traffic on the provided listener using the provided
// handler. The server is shut down when then provided context is cancelled.
func serveHTTP(ctx context.Context, lis net.Listener, handler http.Handler) error {
	server := http.Server{Handler: handler}
	errs := make(chan error, 1)
	go func() { errs <- server.Serve(lis) }()
	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return server.Shutdown(ctx)
	}
}

func (s *BaseService) Start() {

	defer func() {
		if err := recover(); err != nil {
			logger.Error("[Start] ", err, "\n", string(debug.Stack()))
		}
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Debug("CUP启用数量:", runtime.NumCPU())

	s.Inherit.Start()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE)

	for sig := range ch {
		logger.Info("[Start] 进程收到信号 %s", sig)
		switch sig {
		case syscall.SIGHUP:
			s.Inherit.Reload()
		case syscall.SIGPIPE:
		default:
			logger.Info("[Start] 进程收到信号准备退出...")
			close(ch)
			break
		}
	}

	logger.Info("[Start] 进程退出前执行最后的操作...")

	s.Inherit.Stop()
}

func (s *BaseService) Reload() {
}

func (s *BaseService) Init(config interface{}, processId int) {
}

func (s *BaseService) Stop() {
}

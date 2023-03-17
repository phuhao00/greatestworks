package status

import (
	"bytes"
	"context"
	"net/http"

	"greatestworks/aop/logtype"
	"greatestworks/aop/metrics"
	imetrics "greatestworks/aop/metrics"
	"greatestworks/aop/protomsg"
	protos "greatestworks/aop/protos"
)

const (
	statusEndpoint     = "/debug/serviceweaver/status"
	metricsEndpoint    = "/debug/serviceweaver/metrics"
	prometheusEndpoint = "/debug/serviceweaver/prometheus"
	profileEndpoint    = "/debug/serviceweaver/profile"
)

// A Server returns information about a Service Weaver deployment.
type Server interface {
	// Status returns the status of the deployment.
	Status(context.Context) (*Status, error)

	// Metrics returns a snapshot of the deployment's metrics.
	Metrics(context.Context) (*Metrics, error)

	// Profile returns a profile of the deployment.
	Profile(context.Context, *protos.RunProfiling) (*protos.Profile, error)
}

// RegisterServer registers a Server's methods with the provided mux under the
// /debug/serviceweaver/ prefix. You can use a Client to interact with a Status server.
func RegisterServer(mux *http.ServeMux, server Server, logger logtype.Logger) {
	mux.Handle(statusEndpoint, protomsg.HandlerThunk(logger, server.Status))
	mux.Handle(metricsEndpoint, protomsg.HandlerThunk(logger, server.Metrics))
	mux.Handle(profileEndpoint, protomsg.HandlerFunc(logger, server.Profile))
	mux.HandleFunc(prometheusEndpoint, func(w http.ResponseWriter, r *http.Request) {
		ms, err := server.Metrics(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		snapshots := make([]*metrics.MetricSnapshot, len(ms.Metrics))
		for i, m := range ms.Metrics {
			snapshots[i] = metrics.UnProto(m)
		}
		var b bytes.Buffer
		imetrics.TranslateMetricsToPrometheusTextFormat(&b, snapshots, r.Host, prometheusEndpoint)
		w.Write(b.Bytes()) //nolint:errcheck // response write error
	})
}

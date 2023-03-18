package server

import (
	"go.opentelemetry.io/otel/trace"
	"greatestworks/aop/net/call"
)

type BaseServer struct {
}

type stub struct {
	client   call.Connection  // client to talk to the remote component, created lazily.
	methods  []call.MethodKey // Keys for the remote component methods.
	balancer call.Balancer    // if not nil, component load balancer
	tracer   trace.Tracer     // component tracer
}

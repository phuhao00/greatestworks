package status

import (
	"context"
	"net/http"

	"greatestworks/aop/protomsg"
	"greatestworks/aop/protos"
)

// Client is an HTTP client to a status server. It's assumed the status server
// registered itself with RegisterServer.
type Client struct {
	addr string // status server (e.g., "localhost:12345")
}

var _ Server = &Client{}

// NewClient returns a client to the status server on the provided address.
func NewClient(addr string) *Client {
	return &Client{addr}
}

// Status implements the Server interface.
func (c *Client) Status(ctx context.Context) (*Status, error) {
	status := &Status{}
	err := protomsg.Call(ctx, protomsg.CallArgs{
		Client:  http.DefaultClient,
		Addr:    "http://" + c.addr,
		URLPath: statusEndpoint,
		Reply:   status,
	})
	if err != nil {
		return nil, err
	}
	status.StatusAddr = c.addr
	return status, nil
}

// Metrics implements the Server interface.
func (c *Client) Metrics(ctx context.Context) (*Metrics, error) {
	metrics := &Metrics{}
	err := protomsg.Call(ctx, protomsg.CallArgs{
		Client:  http.DefaultClient,
		Addr:    "http://" + c.addr,
		URLPath: metricsEndpoint,
		Reply:   metrics,
	})
	return metrics, err
}

// Profile implements the Server interface.
func (c *Client) Profile(ctx context.Context, req *protos.RunProfiling) (*protos.Profile, error) {
	profile := &protos.Profile{}
	err := protomsg.Call(ctx, protomsg.CallArgs{
		Client:  http.DefaultClient,
		Addr:    "http://" + c.addr,
		URLPath: profileEndpoint,
		Request: req,
		Reply:   profile,
	})
	return profile, err
}

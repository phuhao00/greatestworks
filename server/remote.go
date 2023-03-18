package server

import (
	"fmt"
	"greatestworks/aop/net/call"
)

// parseEndpoints parses a list of endpoint addresses into a list of
// call.Endpoints.
func parseEndpoints(addrs []string) ([]call.Endpoint, error) {
	var endpoints []call.Endpoint
	for _, addr := range addrs {
		endpoint, err := parseEndpoint(addr)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}

// parseEndpoint parses an endpoint address into a call.Endpoint.
func parseEndpoint(endpoint string) (call.Endpoint, error) {
	net, addr, err := call.NetworkAddress(endpoint).Split()
	if err != nil {
		return nil, fmt.Errorf("bad endpoint %q: %w", endpoint, err)
	}
	return call.NetEndpoint{Net: net, Addr: addr}, nil
}

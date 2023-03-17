package call

import (
	"fmt"
	"strings"
)

// A NetworkAddress is a string of the form <network>://<address> (e.g.,
// "tcp://localhost:8000", "unix:///tmp/unix.sock").
type NetworkAddress string

// Split splits the network and address from a NetworkAddress. For example,
//
//	NetworkAddress("tcp://localhost:80").Split() // "tcp", "localhost:80"
//	NetworkAddress("unix://unix.sock").Split()   // "unix", "unix.sock"
func (na NetworkAddress) Split() (network string, address string, err error) {
	net, addr, ok := strings.Cut(string(na), "://")
	if !ok {
		return "", "", fmt.Errorf("%q does not have format <network>://<address>", na)
	}
	return net, addr, nil
}

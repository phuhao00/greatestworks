package call_test

import (
	"testing"

	"greatestworks/aop/net/call"
)

func TestSplit(t *testing.T) {
	for _, test := range []struct {
		name    string
		s       string
		network string
		address string
	}{
		{"NetworkAddress", "network://address", "network", "address"},
		{"EmptyNetwork", "://address", "", "address"},
		{"EmptyAddress", "network://", "network", ""},
		{"JustDelim", "://", "", ""},
		{"ExtraDelim", "network://a://b://c", "network", "a://b://c"},
	} {
		t.Run(test.name, func(t *testing.T) {
			network, address, err := call.NetworkAddress(test.s).Split()
			if err != nil {
				t.Fatalf("%q.Split(): unexpected error: %v", test.s, err)
			}
			if got, want := network, test.network; got != want {
				t.Fatalf("%q.Split() bad network: got %q, want %q", test.s, got, want)
			}
			if got, want := address, test.address; got != want {
				t.Fatalf("%q.Split() bad address: got %q, want %q", test.s, got, want)
			}
		})
	}
}

func TestSplitError(t *testing.T) {
	na := call.NetworkAddress("there is no delimiter here")
	_, _, err := na.Split()
	if err == nil {
		t.Fatalf("%q.Split(): unexpected success", string(na))
	}
}

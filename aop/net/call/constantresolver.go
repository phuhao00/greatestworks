package call

import (
	"context"
	"fmt"
)

// constantResolver is a trivial constant resolver that returns a fixed set of
// endponts.
type constantResolver struct {
	endpoints []Endpoint
}

var _ Resolver = &constantResolver{}

// NewConstantResolver returns a new resolver that returns the provided
// set of endpoints.
func NewConstantResolver(endpoints ...Endpoint) Resolver {
	return &constantResolver{endpoints: endpoints}
}

// IsConstant implements the Resolver interface.
func (*constantResolver) IsConstant() bool {
	return true
}

// Resolve implements the Resolver interface.
func (c *constantResolver) Resolve(_ context.Context, version *Version) ([]Endpoint, *Version, error) {
	if version != nil {
		return nil, nil, fmt.Errorf("unexpected non-nil version %v", *version)
	}
	return c.endpoints, nil, nil
}

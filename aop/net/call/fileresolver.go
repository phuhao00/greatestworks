package call

import (
	"context"
	"fmt"
	"os"
)

var (
	versionDoesntExist = Version{Opaque: "doesn't exist"}
	versionExists      = Version{Opaque: "exists"}
)

// fileResolver is a resolver that returns a given set of endpoints when
// a given file is created.
type fileResolver struct {
	eps      []Endpoint
	filename string
}

var _ Resolver = &fileResolver{}

// NewFileResolver returns a new resolver that returns a given set of endpoints
// when a given file is created.
func NewFileResolver(filename string, endpoints ...Endpoint) Resolver {
	return &fileResolver{eps: endpoints, filename: filename}
}

// IsConstant implements the Resolver interface.
func (*fileResolver) IsConstant() bool { return false }

// Resolve implements the Resolver interface.
func (f *fileResolver) Resolve(ctx context.Context, version *Version) ([]Endpoint, *Version, error) {
	switch {
	case version == nil, *version == versionDoesntExist:
		if _, err := os.Stat(f.filename); err != nil {
			return nil, &versionDoesntExist, nil
		}
		// File exists: return the endpoints and a new version.
		return f.eps, &versionExists, nil
	case *version == versionExists:
		// Endpoints already returned: block until the context is done.
		<-ctx.Done()
		return nil, nil, ctx.Err()
	default:
		return nil, nil, fmt.Errorf("unrecognized version %v", *version)
	}
}

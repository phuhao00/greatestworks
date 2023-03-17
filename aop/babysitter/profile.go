package babysitter

import (
	"context"

	"greatestworks/aop/envelope"
	"greatestworks/aop/protos"
	"greatestworks/aop/tool"
)

// runProfiling runs a profiling request on a set of processes.
func runProfiling(ctx context.Context, req *protos.RunProfiling,
	processes map[string][]*envelope.Envelope) (*protos.Profile, error) {
	// Collect together the groups we want to profile.
	groups := make([][]func() (*protos.Profile, error), 0, len(processes))
	for _, envelopes := range processes {
		group := make([]func() (*protos.Profile, error), 0, len(envelopes))
		for _, e := range envelopes {
			group = append(group, func() (*protos.Profile, error) {
				return e.RunProfiling(ctx, req)
			})
		}
		groups = append(groups, group)
	}
	return tool.ProfileGroups(groups)
}

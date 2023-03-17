package single

import (
	"context"
	"fmt"
	"path/filepath"

	"greatestworks/aop/files"
	"greatestworks/aop/status"
	"greatestworks/aop/tool"
)

var (
	dashboardSpec = &status.DashboardSpec{
		Tool:     "weaver single",
		Registry: defaultRegistry,
		Commands: func(deploymentId string) []status.Command {
			return []status.Command{
				{Label: "status", Command: "weaver single status"},
				{Label: "profile", Command: fmt.Sprintf("weaver single profile --duration=30s %s", deploymentId)},
			}
		},
	}

	Commands = map[string]*tool.Command{
		"status":    status.StatusCommand("weaver single", defaultRegistry),
		"dashboard": status.DashboardCommand(dashboardSpec),
		"metrics":   status.MetricsCommand("weaver single", defaultRegistry),
		"profile":   status.ProfileCommand("weaver single", defaultRegistry),
	}
)

func defaultRegistry(ctx context.Context) (*status.Registry, error) {
	dir, err := files.DefaultDataDir()
	if err != nil {
		return nil, err
	}
	return status.NewRegistry(ctx, filepath.Join(dir, "single_registry"))
}

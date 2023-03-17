package multi

import (
	"context"
	"fmt"
	"path/filepath"

	"greatestworks/aop/logging"
	"greatestworks/aop/status"
	"greatestworks/aop/tool"
)

var (
	// logdir is where weaver multi deployed applications store their logs.
	logdir = filepath.Join(logging.DefaultLogDir, "weaver-multi")

	dashboardSpec = &status.DashboardSpec{
		Tool:     "weaver multi",
		Registry: defaultRegistry,
		Commands: func(deploymentId string) []status.Command {
			return []status.Command{
				{Label: "status", Command: "weaver multi status"},
				{Label: "cat logs", Command: fmt.Sprintf("weaver multi logs 'version==%q'", logging.Shorten(deploymentId))},
				{Label: "follow logs", Command: fmt.Sprintf("weaver multi logs --follow 'version==%q'", logging.Shorten(deploymentId))},
				{Label: "profile", Command: fmt.Sprintf("weaver multi profile --duration=30s %s", deploymentId)},
			}
		},
	}

	Commands = map[string]*tool.Command{
		"deploy": &deployCmd,
		"logs": tool.LogsCmd(&tool.LogsSpec{
			Tool: "weaver multi",
			Source: func(context.Context) (logging.Source, error) {
				return logging.FileSource(logdir), nil
			},
		}),
		"dashboard": status.DashboardCommand(dashboardSpec),
		"status":    status.StatusCommand("weaver multi", defaultRegistry),
		"metrics":   status.MetricsCommand("weaver multi", defaultRegistry),
		"profile":   status.ProfileCommand("weaver multi", defaultRegistry),
	}
)

package ssh

import (
	"fmt"

	"greatestworks/aop/logging"
	"greatestworks/aop/status"
	"greatestworks/aop/tool/ssh/impl"
)

var dashboardSpec = &status.DashboardSpec{
	Tool:     "weaver ssh",
	Registry: impl.DefaultRegistry,
	Commands: func(deploymentId string) []status.Command {
		return []status.Command{
			{Label: "cat logs", Command: fmt.Sprintf("weaver ssh logs 'version==%q'", logging.Shorten(deploymentId))},
			{Label: "follow logs", Command: fmt.Sprintf("weaver ssh logs --follow 'version==%q'", logging.Shorten(deploymentId))},
		}
	},
}

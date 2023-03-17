package ssh

import (
	"path/filepath"

	"greatestworks/aop/logging"
	"greatestworks/aop/status"
	"greatestworks/aop/tool"
)

var (
	// logDir is where weaver ssh deployed applications store their logs.
	logDir = filepath.Join(logging.DefaultLogDir, "weaver_ssh")

	Commands = map[string]*tool.Command{
		"deploy":    &deployCmd,
		"logs":      tool.LogsCmd(&logsSpec),
		"dashboard": status.DashboardCommand(dashboardSpec),

		// Hidden commands.
		"babysitter": &babysitterCmd,
	}
)

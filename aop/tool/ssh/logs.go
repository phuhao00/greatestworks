package ssh

import (
	"context"

	"greatestworks/aop/logging"
	"greatestworks/aop/tool"
)

var logsSpec = tool.LogsSpec{
	Tool: "weaver ssh",
	Source: func(context.Context) (logging.Source, error) {
		return logging.FileSource(logDir), nil
	},
}

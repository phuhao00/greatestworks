package server

import (
	"greatestworks/aop/logging"
	"greatestworks/aop/logtype"
	"greatestworks/aop/protos"
)

type Logger interface {
	// Debug logs an entry for debugging purposes that contains msg and attributes.
	Debug(msg string, attributes ...any)

	// Info logs an informational entry that contains msg and attributes.
	Info(msg string, attributes ...any)

	// Error logs an error entry that contains msg and attributes.
	// If err is not nil, it will be used as the value of a attribute
	// named "err".
	Error(msg string, err error, attributes ...any)

	// With returns a logger that will automatically add the
	// pre-specified attributes to all logged entries.
	With(attributes ...any) Logger
}

type attrLogger struct {
	logging.FuncLogger
}

// With implements [weaver.Logger].
func (l attrLogger) With(attrs ...any) Logger {
	result := l
	result.Opts.Attrs = logtype.AppendAttrs(l.Opts.Attrs, attrs)
	return result
}

func newAttrLogger(app, version string, saver func(*protos.LogEntry)) Logger {
	return attrLogger{
		logging.FuncLogger{
			Opts: logging.Options{
				App:        app,
				Deployment: version,
			},
			Write: saver,
		},
	}
}

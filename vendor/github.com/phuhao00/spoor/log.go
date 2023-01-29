package spoor

import (
	"fmt"
	"io"
	"strings"
)

const (
	DEBUG = Level(1)
	INFO  = Level(2)
	WARN  = Level(3)
	ERROR = Level(4)
	FATAL = Level(5)
)

type AppLogFunc func(lvl Level, f string, args ...interface{})

type Logger interface {
	Output(callerSkip int, s string) error
	SetOutput(w io.Writer)
}

type Level int

func (l *Level) Get() interface{} { return *l }

func (l *Level) Set(s string) error {
	lvl, err := ParseLogLevel(s)
	if err != nil {
		return err
	}
	*l = lvl
	return nil
}

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "invalid"
}

func ParseLogLevel(levelStr string) (Level, error) {
	switch strings.ToLower(levelStr) {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	}
	return 0, fmt.Errorf("invalid log level '%s' (debug, info, warn, error, fatal)", levelStr)
}

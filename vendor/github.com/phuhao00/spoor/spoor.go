package spoor

import (
	"fmt"
	"io"
	"log"
)

type Spoor struct {
	l        Logger
	cfgLevel Level
	prefix   string
	flag     int
}

type Option func(spoor *Spoor)

func WithFileWriter(writer *FileWriter) Option {
	return func(spoor *Spoor) {
		writer.level = spoor.cfgLevel
		spoor.l.SetOutput(writer)
	}
}

func WithConsoleWriter(writer io.Writer) Option {
	return func(spoor *Spoor) {
		spoor.l.SetOutput(writer)
	}
}

func NewSpoor(cfgLevel Level, prefix string, flag int, opts ...Option) *Spoor {
	logger := log.New(io.Discard, prefix, flag)
	s := &Spoor{
		l:        logger,
		cfgLevel: cfgLevel,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

//DebugF Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg
func (s *Spoor) DebugF(f string, args ...interface{}) {
	if s.checkLevel(DEBUG) {
		return
	}
	s.l.Output(2, fmt.Sprintf(DEBUG.String()+" "+f, args...))
}

func (s *Spoor) ErrorF(f string, args ...interface{}) {
	if s.checkLevel(ERROR) {
		return
	}
	s.l.Output(2, fmt.Sprintf(ERROR.String()+" "+f, args...))
}

func (s *Spoor) InfoF(f string, args ...interface{}) {
	if s.checkLevel(INFO) {
		return
	}
	s.l.Output(2, fmt.Sprintf(INFO.String()+" "+f, args...))
}

func (s *Spoor) FatalF(f string, args ...interface{}) {
	if s.checkLevel(FATAL) {
		return
	}
	s.l.Output(2, fmt.Sprintf(FATAL.String()+" "+f, args...))
}

func (l *Spoor) checkLevel(level Level) bool {
	if level >= l.cfgLevel {
		return false
	}
	return true
}

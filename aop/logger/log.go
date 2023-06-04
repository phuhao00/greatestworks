package logger

import (
	"fmt"
	"log"
	"sync"

	"github.com/phuhao00/spoor"
)

var (
	sp             *spoor.Spoor
	onceInitLogger sync.Once
)

func GetLogger() *spoor.Spoor {
	return sp
}

type LoggingSetting struct {
	Dir          string
	Level        int
	Prefix       string
	WriterOption spoor.Option
}

func SetLogging(setting *LoggingSetting) {
	onceInitLogger.Do(func() {
		var opt spoor.Option
		if setting.WriterOption == nil {
			fileWriter := spoor.NewFileWriter(setting.Dir, 0, 0, 0)
			opt = spoor.WithFileWriter(fileWriter)
		} else {
			opt = setting.WriterOption
		}
		l := spoor.NewSpoor(spoor.Level(setting.Level), setting.Prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, opt)
		sp = l
	})
}

// DebugF Log line format: [IWEF]mmdd hh:mm:sLogger.uuuuuu threadid file:line] msg
func DebugF(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.DEBUG) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.DEBUG.String()+" "+f, args...))
}

func ErrorF(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.ERROR) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.ERROR.String()+" "+f, args...))
}

func InfoF(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.INFO) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.INFO.String()+" "+f, args...))
}

func FatalF(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.FATAL) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.FATAL.String()+" "+f, args...))
}

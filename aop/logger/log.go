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

// Debug Log line format: [IWEF]mmdd hh:mm:sLogger.uuuuuu threadid file:line] msg
func Debug(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.DEBUG) {
		return
	}
	err := sp.Output(2, fmt.Sprintf(spoor.DEBUG.String()+" "+f, args...))
	if err != nil {
		fmt.Println(err)
	}
}

func Error(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.ERROR) {
		return
	}
	err := sp.Output(2, fmt.Sprintf(spoor.ERROR.String()+" "+f, args...))
	if err != nil {
		fmt.Println(err)
	}
}

func Info(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.INFO) {
		return
	}
	err := sp.Output(2, fmt.Sprintf(spoor.INFO.String()+" "+f, args...))
	if err != nil {
		fmt.Println(err)
	}
}

func Warn(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.WARN) {
		return
	}
	err := sp.Output(2, fmt.Sprintf(spoor.WARN.String()+" "+f, args...))
	if err != nil {
		fmt.Println(err)
	}
}

func Fatal(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.FATAL) {
		return
	}
	err := sp.Output(2, fmt.Sprintf(spoor.FATAL.String()+" "+f, args...))
	if err != nil {
		fmt.Println(err)
	}
}

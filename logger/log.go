package logger

import (
	"log"
	"sync"

	"github.com/phuhao00/spoor"
)

var (
	Logger         *spoor.Spoor
	onceInitLogger sync.Once
)

func init() {
	onceInitLogger.Do(func() {
		fileWriter := spoor.NewFileWriter("log", 0, 0, 0)
		l := spoor.NewSpoor(spoor.DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, spoor.WithFileWriter(fileWriter))
		Logger = l
	})
}

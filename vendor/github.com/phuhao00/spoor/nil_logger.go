package spoor

import "io"

type NilLogger struct{}

func (l NilLogger) Output(callerSkip int, s string) error {
	return nil
}

func (l *NilLogger) SetOutput(writer io.Writer) {

}

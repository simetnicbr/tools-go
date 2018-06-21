package logger

import "fmt"

type MockLogger struct {
	MockMsg func(m ...interface{})
	Depth   int
	Context Fields
}

func (l MockLogger) SetDepth(d int) {
	l.Depth = d
}

func (l MockLogger) AddContext(f Fields) {
	l.Context = f
}

func (l MockLogger) Fatal(f string, m ...interface{}) {
	l.MockMsg(fmt.Sprintf(f, m...))
}

func (l MockLogger) Warning(f string, m ...interface{}) {
	l.MockMsg(fmt.Sprintf(f, m...))
}

func (l MockLogger) Info(f string, m ...interface{}) {
	l.MockMsg(fmt.Sprintf(f, m...))
}

func (l MockLogger) Debug(f string, m ...interface{}) {
	l.MockMsg(fmt.Sprintf(f, m...))
}

package logger

import "fmt"

type MockLogger struct {
	MockMsg func(m ...interface{})
}

func (l MockLogger) Warning(m ...interface{}) {
	l.MockMsg(m...)
}

func (l MockLogger) Info(m ...interface{}) {
	l.MockMsg(m...)
}

func (l MockLogger) Debug(m ...interface{}) {
	l.MockMsg(m...)
}

func (l MockLogger) Warningm(m map[string]interface{}) {

}

func (l MockLogger) Infom(m map[string]interface{}) {

}

func (l MockLogger) Debugm(m map[string]interface{}) {

}

func (l MockLogger) Warningf(f string, m ...interface{}) {
	l.MockMsg(fmt.Sprintf(f, m...))
}

func (l MockLogger) Infof(f string, m ...interface{}) {
	l.MockMsg(fmt.Sprintf(f, m...))
}

func (l MockLogger) Debugf(f string, m ...interface{}) {
	l.MockMsg(fmt.Sprintf(f, m...))
}

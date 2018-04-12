package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

// logger constants
const (
	LEVEL string = "SIMET_LOG_LEVEL"
	DEBUG string = "debug"
)

// Logger represents the log interface used
type Logger interface {
	SetDepth(d int)

	Warning(m ...interface{})
	Info(m ...interface{})
	Debug(m ...interface{})

	Warningm(map[string]interface{})
	Infom(map[string]interface{})
	Debugm(map[string]interface{})

	Warningf(f string, m ...interface{})
	Infof(f string, m ...interface{})
	Debugf(f string, m ...interface{})
}

// New creates a new instance of log that implements Logger Interface
func New() Logger {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}

	level := logrus.InfoLevel

	logLevel := os.Getenv(LEVEL)
	if len(logLevel) > 0 && strings.EqualFold(logLevel, DEBUG) {
		level = logrus.DebugLevel
	}

	log.SetLevel(level)

	l := logWrapper{
		Logger: log,
		Depth:  2,
	}

	return &l
}

type logWrapper struct {
	Logger *logrus.Logger
	Depth  int
}

func (l logWrapper) SetDepth(d int) {
	l.Depth = d
}

func (l logWrapper) Warning(m ...interface{}) {
	msg := l.prefix(fmt.Sprint(m...))

	l.Logger.Warning(msg)
}

func (l logWrapper) Info(m ...interface{}) {
	msg := l.prefix(fmt.Sprint(m...))

	l.Logger.Info(msg)
}

func (l logWrapper) Debug(m ...interface{}) {
	msg := l.prefix(fmt.Sprint(m...))

	l.Logger.Debug(msg)
}

func (l logWrapper) Warningm(m map[string]interface{}) {
	m["metadata"] = l.metadata()
	l.Logger.WithFields(m).Warning()
}

func (l logWrapper) Infom(m map[string]interface{}) {
	m["metadata"] = l.metadata()
	l.Logger.WithFields(m).Info()
}

func (l logWrapper) Debugm(m map[string]interface{}) {
	m["metadata"] = l.metadata()
	l.Logger.WithFields(m).Debug()
}

func (l logWrapper) Warningf(f string, m ...interface{}) {
	msg := l.prefix(fmt.Sprintf(f, m...))

	l.Logger.Warningf(msg)
}

func (l logWrapper) Infof(f string, m ...interface{}) {
	msg := l.prefix(fmt.Sprintf(f, m...))

	l.Logger.Info(msg)
}

func (l logWrapper) Debugf(f string, m ...interface{}) {
	msg := l.prefix(fmt.Sprintf(f, m...))

	l.Logger.Debug(msg)
}

func (l logWrapper) prefix(msg string) string {
	_, filepath, line, _ := runtime.Caller(l.Depth)

	tokens := strings.Split(filepath, "/")
	file := tokens[len(tokens)-1]
	if len(tokens) >= 2 {
		file = tokens[len(tokens)-2] + "/" + tokens[len(tokens)-1]
	}

	return fmt.Sprintf("{%v: %v} %v", file, line, msg)
}

func (l logWrapper) metadata() string {
	_, filepath, line, _ := runtime.Caller(l.Depth)

	tokens := strings.Split(filepath, "/")
	file := tokens[len(tokens)-1]
	if len(tokens) >= 2 {
		file = tokens[len(tokens)-2] + "/" + tokens[len(tokens)-1]
	}

	return fmt.Sprintf("[%v - %v]", file, line)
}
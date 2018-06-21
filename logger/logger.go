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

	AddFields(f Fields)

	Fatal(s string, m ...interface{})
	Warning(s string, m ...interface{})
	Info(s string, m ...interface{})
	Debug(s string, m ...interface{})
}

// Fields is an alias for logrus.Fields
type Fields logrus.Fields

func New() Logger {
	return newLogger(make(map[string]interface{}))
}

func WithFields(fields Fields) Logger {
	return newLogger(fields)
}

// New creates a new instance of log that implements Logger Interface
func newLogger(f Fields) *logWrapper {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}

	// default log level -> INFO
	level := logrus.InfoLevel

	// DEBUG or INFO
	logLevel := os.Getenv(LEVEL)
	if len(logLevel) > 0 && strings.EqualFold(logLevel, DEBUG) {
		level = logrus.DebugLevel
	}

	log.SetLevel(level)

	l := logWrapper{
		Logger:  log,
		Depth:   2,
		Context: f,
	}

	return &l
}

type logWrapper struct {
	Logger  *logrus.Logger
	Depth   int
	Context Fields
}

// Set the function Depth
func (l logWrapper) SetDepth(d int) {
	l.Depth = d
}

// AddFields add fields to logger context
func (l logWrapper) AddFields(f Fields) {
	for k, v := range f {
		l.Context[k] = v
	}
}

func (l logWrapper) remove(key string) {
	delete(l.Context, key)
}

// Log funcs
func (l logWrapper) Fatal(f string, msg ...interface{}) {
	l.AddFields(Fields{
		"file": l.file(),
	})
	l.Logger.WithFields((logrus.Fields(l.Context))).Fatal(fmt.Sprintf(f, msg...))
	l.remove("file")
}

func (l logWrapper) Warning(f string, msg ...interface{}) {
	l.AddFields(Fields{
		"file": l.file(),
	})
	l.Logger.WithFields((logrus.Fields(l.Context))).Warning(fmt.Sprintf(f, msg...))
	l.remove("file")
}

func (l logWrapper) Info(f string, msg ...interface{}) {
	l.AddFields(Fields{
		"file": l.file(),
	})
	l.Logger.WithFields((logrus.Fields(l.Context))).Info(fmt.Sprintf(f, msg...))
	l.remove("file")
}

func (l logWrapper) Debug(f string, msg ...interface{}) {
	l.AddFields(Fields{
		"file": l.file(),
	})
	l.Logger.WithFields((logrus.Fields(l.Context))).Debug(fmt.Sprintf(f, msg...))
	l.remove("file")
}

func (l logWrapper) file() string {
	_, filepath, line, _ := runtime.Caller(l.Depth)

	tokens := strings.Split(filepath, "/")
	file := tokens[len(tokens)-1]
	if len(tokens) >= 2 {
		file = tokens[len(tokens)-2] + "/" + tokens[len(tokens)-1]
	}

	return fmt.Sprintf("[%v - %v]", file, line)
}

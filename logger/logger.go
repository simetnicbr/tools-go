// Package logger implements a wrapper arount logrus.
// It mainly enriches logs with a field 'file', indicating the source code line where
// the logger function was called.
//
// Instance of the logger may be used concurrently from multiple goroutines. All access
// to shared data is synconrized. Shared data that is copied before passing on to logrus
// functions.
package logger

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	logrus "github.com/sirupsen/logrus"
)

const (
	JsonFormatter = "json"
	TextFormatter = "text"
)

// Logger represents the log interface used
type Logger interface {
	SetDepth(d int)
	SetFormatter(string)

	AddFields(f Fields)

	Fatal(s string, m ...interface{})
	Warning(s string, m ...interface{})
	Info(s string, m ...interface{})
	Debug(s string, m ...interface{})
	Error(s string, m ...interface{})
}

// Fields is an alias for logrus.Fields
type Fields logrus.Fields

func New(level string) Logger {
	return newLogger(level, make(map[string]interface{}))
}

func WithFields(level string, fields Fields) Logger {
	return newLogger(level, fields)
}

// New creates a new instance of log that implements Logger Interface
func newLogger(lvl string, f Fields) *logWrapper {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}

	// default log level -> INFO
	level := logrus.InfoLevel

	switch lvl {
	case "debug":
		level = logrus.DebugLevel
	case "warning":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	}

	log.SetLevel(level)

	l := logWrapper{
		Level:        level,
		Logger:       log,
		Depth:        2,
		Context:      f,
		ContextMutex: &sync.RWMutex{},
	}

	return &l
}

type logWrapper struct {
	Level        logrus.Level
	Logger       *logrus.Logger
	Depth        int
	Context      Fields
	ContextMutex *sync.RWMutex
}

// Set the function Depth
func (l logWrapper) SetDepth(d int) {
	l.Depth = d
}

func (l logWrapper) SetFormatter(fomatter string) {
	switch strings.ToLower(fomatter) {
	case TextFormatter:
		l.Logger.Formatter = &logrus.TextFormatter{}
	default:
		l.Logger.Formatter = &logrus.JSONFormatter{}
	}
}

// AddFields add fields to logger context
// The function is threadsafe
func (l logWrapper) AddFields(f Fields) {
	l.ContextMutex.Lock()
	defer l.ContextMutex.Unlock()
	for k, v := range f {
		l.Context[k] = v
	}
}

// Log functions
//
// All log functions are safe for concurrent invocation on multiple goroutines.
// Log functions do not operate on shared data. Instead they create a local copy of the
// logger's shared data for each invoation. Besides mutating only this local copy,
// it is also passed to logrus when delegating the logging. This way we are sure that
// logrus will never perform a concurrent read of our shared data.
// Note that our shared data may be mutated at any point in time as this package exposes
// a public function AddFields(). While function AddFields() mutation of the logger's
// shared data is protected, this is not enough, as logrus may attempt concurrent read attempt.
//
// Our log functions are a bit more expensive then the original logrus functions,
// due to synconization and copy effort. To reduce the overhead we step out of log functions
// as early as possible, when the logger's log level is less verbose then the invoced log function.
func (l logWrapper) Fatal(f string, msg ...interface{}) {
	if l.Level < logrus.FatalLevel {
		return
	}
	fields := l.getThreadsafeCopyOfFields()
	fields["file"] = l.file()
	l.Logger.WithFields(fields).Fatal(fmt.Sprintf(f, msg...))
}

func (l logWrapper) Warning(f string, msg ...interface{}) {
	if l.Level < logrus.WarnLevel {
		return
	}
	fields := l.getThreadsafeCopyOfFields()
	fields["file"] = l.file()
	l.Logger.WithFields(fields).Warning(fmt.Sprintf(f, msg...))
}

func (l logWrapper) Info(f string, msg ...interface{}) {
	if l.Level < logrus.InfoLevel {
		return
	}
	fields := l.getThreadsafeCopyOfFields()
	fields["file"] = l.file()
	l.Logger.WithFields(fields).Info(fmt.Sprintf(f, msg...))
}

func (l logWrapper) Debug(f string, msg ...interface{}) {
	if l.Level < logrus.DebugLevel {
		return
	}
	fields := l.getThreadsafeCopyOfFields()
	fields["file"] = l.file()
	l.Logger.WithFields(fields).Debug(fmt.Sprintf(f, msg...))
}

func (l logWrapper) Error(f string, msg ...interface{}) {
	if l.Level < logrus.ErrorLevel {
		return
	}
	fields := l.getThreadsafeCopyOfFields()
	fields["file"] = l.file()
	l.Logger.WithFields(fields).Error(fmt.Sprintf(f, msg...))
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

// create a complete copy of the logger's field map; the copy process is syncronized;
// this function is intended to be used when passing the field map to third party code
// (logrus) where we do not control and syncronize access to our data structures;
func (l logWrapper) getThreadsafeCopyOfFields() logrus.Fields {
	copiedFields := make(Fields)
	l.ContextMutex.RLock()
	defer l.ContextMutex.RUnlock()
	for k, v := range l.Context {
		copiedFields[k] = v
	}
	return logrus.Fields(copiedFields)
}

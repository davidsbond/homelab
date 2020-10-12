// Package logging contains utilities for writing logs.
package logging

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

var level string

// Init initialises the logger.
func Init() {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.ErrorLevel
	}

	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

// WithError adds an error to the log.
func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

// WithField adds a field to the log.
func WithField(k string, v interface{}) *logrus.Entry {
	return logrus.WithField(k, v)
}

// Info writes an info level log.
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Logger returns the global logger.
func Logger() *logrus.Logger {
	return logrus.StandardLogger()
}

type (
	jaegerLogger struct {
		*logrus.Logger
	}
)

// JaegerLogger returns the standard logger as an implementation of jaeger.Logger.
func JaegerLogger() jaeger.Logger {
	return &jaegerLogger{Logger: Logger()}
}

// Error writes an error message to the log.
func (l *jaegerLogger) Error(msg string) {
	l.Logger.Error(msg)
}

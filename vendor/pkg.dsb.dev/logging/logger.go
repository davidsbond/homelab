// Package logging contains utilities for writing logs.
package logging

import (
	"os"

	"github.com/sirupsen/logrus"
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

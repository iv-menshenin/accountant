package ete_test

import "log"

type (
	testLogger struct {
		l *log.Logger
	}
)

func (l *testLogger) Warning(format string, args ...interface{}) {
	l.l.Printf("WARNING: "+format, args...)
}

func (l *testLogger) Debug(format string, args ...interface{}) {
	l.l.Printf("DEBUG: "+format, args...)
}

func (l *testLogger) Error(format string, args ...interface{}) {
	l.l.Printf("ERROR: "+format, args...)
}

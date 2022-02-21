package logger

import "log"

type (
	Logger struct {
		level  LogLevel
		logger *log.Logger
	}
	LogLevel int
)

const (
	LogLevelError LogLevel = iota
	LogLevelWarning
	LogLevelDebug
)

func NewFromLogger(logger *log.Logger, level LogLevel) Logger {
	return Logger{
		level:  level,
		logger: logger,
	}
}

func (l Logger) Warning(format string, args ...interface{}) {
	if l.level <= LogLevelWarning {
		l.logger.Printf("[WARNING] "+format+"\n", args...)
	}
}

func (l Logger) Debug(format string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.logger.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func (l Logger) Error(format string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.logger.Printf("[ERROR] "+format+"\n", args...)
	}
}

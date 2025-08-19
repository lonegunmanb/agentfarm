// Package config provides logging setup for the Agent Farm collective
package pkg

import (
	"io"
	"log"
	"os"
	"strings"
)

// Logger represents the revolutionary logging system for the collective
type Logger struct {
	*log.Logger
	level LogLevel
}

// LogLevel represents the revolutionary logging levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "INFO"
	}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// NewLogger creates a new revolutionary logger for the collective
func NewLogger(level string, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	logLevel := ParseLogLevel(level)
	logger := log.New(output, "", log.LstdFlags|log.Lshortfile)

	return &Logger{
		Logger: logger,
		level:  logLevel,
	}
}

// Debug logs debug messages for revolutionary development
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs informational messages for the collective
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.Printf("[INFO] "+format, v...)
	}
}

// Warn logs warning messages for collective attention
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.Printf("[WARN] "+format, v...)
	}
}

// Error logs error messages for revolutionary security
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs fatal errors and exits (for counter-revolutionary activities)
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}

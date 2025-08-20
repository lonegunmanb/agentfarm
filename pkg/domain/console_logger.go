package domain

import (
	"fmt"
	"log"
	"time"
)

// ConsoleLogger implements Logger interface for console output
// This is a simple implementation that outputs to stdout/stderr
type ConsoleLogger struct {
	debugEnabled bool
}

// NewConsoleLogger creates a new console logger
func NewConsoleLogger(debugEnabled bool) *ConsoleLogger {
	return &ConsoleLogger{
		debugEnabled: debugEnabled,
	}
}

// Info logs an informational message to stdout
func (c *ConsoleLogger) Info(message string, fields ...map[string]interface{}) {
	c.logWithLevel("INFO", message, fields...)
}

// Error logs an error message to stderr
func (c *ConsoleLogger) Error(message string, fields ...map[string]interface{}) {
	c.logWithLevel("ERROR", message, fields...)
}

// Debug logs a debug message to stdout (only if debug is enabled)
func (c *ConsoleLogger) Debug(message string, fields ...map[string]interface{}) {
	if c.debugEnabled {
		c.logWithLevel("DEBUG", message, fields...)
	}
}

// Warn logs a warning message to stdout
func (c *ConsoleLogger) Warn(message string, fields ...map[string]interface{}) {
	c.logWithLevel("WARN", message, fields...)
}

// logWithLevel outputs a formatted log message with level, timestamp and optional fields
func (c *ConsoleLogger) logWithLevel(level string, message string, fields ...map[string]interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] %s - %s", level, timestamp, message)
	
	// Add fields if provided
	if len(fields) > 0 && fields[0] != nil {
		for key, value := range fields[0] {
			logMsg += fmt.Sprintf(" | %s=%v", key, value)
		}
	}
	
	log.Println(logMsg)
}

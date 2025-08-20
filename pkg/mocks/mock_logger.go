package mocks

import (
	"sync"
	"time"

	"github.com/lonegunmanb/agentfarm/pkg/ports"
)

// MockLogger implements Logger interface for testing
type MockLogger struct {
	mu   sync.RWMutex
	logs []ports.LogEntry
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{
		logs: make([]ports.LogEntry, 0),
	}
}

// Info logs an informational message
func (m *MockLogger) Info(message string, fields ...map[string]interface{}) {
	m.logEntry("INFO", message, fields...)
}

// Error logs an error message
func (m *MockLogger) Error(message string, fields ...map[string]interface{}) {
	m.logEntry("ERROR", message, fields...)
}

// Debug logs a debug message
func (m *MockLogger) Debug(message string, fields ...map[string]interface{}) {
	m.logEntry("DEBUG", message, fields...)
}

// Warn logs a warning message
func (m *MockLogger) Warn(message string, fields ...map[string]interface{}) {
	m.logEntry("WARN", message, fields...)
}

// logEntry is a helper method to create log entries
func (m *MockLogger) logEntry(level string, message string, fields ...map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var mergedFields map[string]interface{}
	if len(fields) > 0 {
		mergedFields = make(map[string]interface{})
		for _, fieldMap := range fields {
			for k, v := range fieldMap {
				mergedFields[k] = v
			}
		}
	}

	entry := ports.LogEntry{
		Level:   level,
		Message: message,
		Fields:  mergedFields,
		Time:    time.Now().Format(time.RFC3339),
	}

	m.logs = append(m.logs, entry)
}

// GetLogs returns all log entries (for testing)
func (m *MockLogger) GetLogs() []ports.LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]ports.LogEntry, len(m.logs))
	copy(result, m.logs)
	return result
}

// ClearLogs clears all log entries (for testing)
func (m *MockLogger) ClearLogs() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = m.logs[:0]
}

// GetLogsByLevel returns log entries filtered by level (for testing)
func (m *MockLogger) GetLogsByLevel(level string) []ports.LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []ports.LogEntry
	for _, log := range m.logs {
		if log.Level == level {
			result = append(result, log)
		}
	}
	return result
}

// Verify interface compliance
var _ ports.Logger = (*MockLogger)(nil)

package ports

// Logger defines the port for logging operations
// This interface abstracts logging operations from the core domain
type Logger interface {
	// Info logs an informational message
	Info(message string, fields ...map[string]interface{})

	// Error logs an error message
	Error(message string, fields ...map[string]interface{})

	// Debug logs a debug message
	Debug(message string, fields ...map[string]interface{})

	// Warn logs a warning message
	Warn(message string, fields ...map[string]interface{})
}

// LogEntry represents a log entry (for testing/monitoring)
type LogEntry struct {
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
	Time    string                 `json:"time"`
}

package ports

// MessageSender defines the port for sending messages to external agents
// This interface abstracts message delivery operations from the core domain
type MessageSender interface {
	// SendActivation sends an activation message to an agent
	SendActivation(role string, payload string) error

	// SendDeactivation sends a deactivation message to an agent
	SendDeactivation(role string, reason string) error

	// SendNotification sends a general notification to an agent
	SendNotification(role string, message string) error
}

// SentMessage represents a message that was sent (for testing/monitoring)
type SentMessage struct {
	Recipient string                 `json:"recipient"`
	Type      string                 `json:"type"`
	Payload   string                 `json:"payload"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

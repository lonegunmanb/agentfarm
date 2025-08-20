package domain

// MessageSender defines the port for sending messages to external agents
// This interface abstracts message delivery operations from the core domain
type MessageSender interface {
	// SendActivation sends an activation message to an agent
	SendActivation(role string, payload string) error
}

// SentMessage represents a message that was sent (for testing/monitoring)
type SentMessage struct {
	Recipient string                 `json:"recipient"`
	Type      string                 `json:"type"`
	Payload   string                 `json:"payload"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

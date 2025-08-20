package mocks

import (
	"sync"

	"github.com/lonegunmanb/agentfarm/pkg/ports"
)

// MockMessageSender implements MessageSender interface for testing
type MockMessageSender struct {
	mu       sync.RWMutex
	messages []ports.SentMessage
}

// NewMockMessageSender creates a new mock message sender
func NewMockMessageSender() *MockMessageSender {
	return &MockMessageSender{
		messages: make([]ports.SentMessage, 0),
	}
}

// SendActivation sends an activation message to an agent
func (m *MockMessageSender) SendActivation(role string, payload string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	message := ports.SentMessage{
		Recipient: role,
		Type:      "activation",
		Payload:   payload,
		Metadata: map[string]interface{}{
			"action": "activate",
		},
	}

	m.messages = append(m.messages, message)
	return nil
}

// GetSentMessages returns all sent messages (for testing)
func (m *MockMessageSender) GetSentMessages() []ports.SentMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]ports.SentMessage, len(m.messages))
	copy(result, m.messages)
	return result
}

// ClearMessages clears all sent messages (for testing)
func (m *MockMessageSender) ClearMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = m.messages[:0]
}

// Verify interface compliance
var _ ports.MessageSender = (*MockMessageSender)(nil)

package mocks

import (
	"sync"
	"time"

	"github.com/lonegunmanb/agentfarm/pkg/ports"
)

// MockEventPublisher implements EventPublisher interface for testing
type MockEventPublisher struct {
	mu     sync.RWMutex
	events []ports.DomainEvent
}

// NewMockEventPublisher creates a new mock event publisher
func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		events: make([]ports.DomainEvent, 0),
	}
}

// PublishAgentRegistered publishes an agent registration event
func (m *MockEventPublisher) PublishAgentRegistered(role string, agentType string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	event := ports.DomainEvent{
		Type: "agent_registered",
		Data: map[string]interface{}{
			"role":       role,
			"agent_type": agentType,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	m.events = append(m.events, event)
	return nil
}

// PublishAgentDeregistered publishes an agent deregistration event
func (m *MockEventPublisher) PublishAgentDeregistered(role string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	event := ports.DomainEvent{
		Type: "agent_deregistered",
		Data: map[string]interface{}{
			"role":   role,
			"reason": reason,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	m.events = append(m.events, event)
	return nil
}

// PublishBarrelTransferred publishes a barrel transfer event
func (m *MockEventPublisher) PublishBarrelTransferred(fromRole string, toRole string, payload string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	event := ports.DomainEvent{
		Type: "barrel_transferred",
		Data: map[string]interface{}{
			"from_role": fromRole,
			"to_role":   toRole,
			"payload":   payload,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	m.events = append(m.events, event)
	return nil
}

// PublishSystemStatusChanged publishes a system status change event
func (m *MockEventPublisher) PublishSystemStatusChanged(status string, details map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data := map[string]interface{}{
		"status": status,
	}

	// Merge details into data
	for k, v := range details {
		data[k] = v
	}

	event := ports.DomainEvent{
		Type:      "system_status_changed",
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	m.events = append(m.events, event)
	return nil
}

// GetPublishedEvents returns all published events (for testing)
func (m *MockEventPublisher) GetPublishedEvents() []ports.DomainEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]ports.DomainEvent, len(m.events))
	copy(result, m.events)
	return result
}

// ClearEvents clears all published events (for testing)
func (m *MockEventPublisher) ClearEvents() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.events = m.events[:0]
}

// Verify interface compliance
var _ ports.EventPublisher = (*MockEventPublisher)(nil)

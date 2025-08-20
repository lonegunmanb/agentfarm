package domain

import (
	"time"
)

// YieldMessage represents a yield operation in the Agent Farm collective.
// This is the only true domain message - carrying payload from one agent to another.
// All other communications (register, activate, query, etc.) are operations, not messages.
type YieldMessage struct {
	fromRole  string
	toRole    string
	payload   string
	timestamp time.Time
}

// NewYieldMessage creates a new yield message
func NewYieldMessage(fromRole, toRole, payload string) YieldMessage {
	return YieldMessage{
		fromRole:  fromRole,
		toRole:    toRole,
		payload:   payload,
		timestamp: nowFunc(),
	}
}

// FromRole returns the sender role
func (m YieldMessage) FromRole() string {
	return m.fromRole
}

// ToRole returns the recipient role
func (m YieldMessage) ToRole() string {
	return m.toRole
}

// Payload returns the message payload
func (m YieldMessage) Payload() string {
	return m.payload
}

// Timestamp returns when the message was created
func (m YieldMessage) Timestamp() time.Time {
	return m.timestamp
}

// IsValid checks if the yield message is valid
func (m YieldMessage) IsValid() bool {
	// Yield message must have from role, to role, and non-zero timestamp
	hasFromRole := m.fromRole != ""
	hasToRole := m.toRole != ""
	hasTimestamp := !m.timestamp.IsZero()
	return hasFromRole && hasToRole && hasTimestamp
}

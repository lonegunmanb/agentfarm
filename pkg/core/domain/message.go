package domain

import (
	"time"
)

// MessageType represents the type of revolutionary message
type MessageType int

const (
	MessageTypeRegister MessageType = iota
	MessageTypeYield
	MessageTypeActivate
	MessageTypeQueryAgents
	MessageTypeAgentList
	MessageTypeError
	MessageTypeAckRegister
)

// String returns the string representation of MessageType
func (m MessageType) String() string {
	switch m {
	case MessageTypeRegister:
		return "REGISTER"
	case MessageTypeYield:
		return "YIELD"
	case MessageTypeActivate:
		return "ACTIVATE"
	case MessageTypeQueryAgents:
		return "QUERY_AGENTS"
	case MessageTypeAgentList:
		return "AGENT_LIST"
	case MessageTypeError:
		return "ERROR"
	case MessageTypeAckRegister:
		return "ACK_REGISTER"
	default:
		return "UNKNOWN"
	}
}

// RevolutionaryMessage represents internal communication within the Agent Farm collective.
// This is the domain representation, independent of transport format (TCP/JSON).
type RevolutionaryMessage struct {
	messageType MessageType
	fromRole    string
	toRole      string
	payload     string
	timestamp   time.Time
}

// NewRevolutionaryMessage creates a new revolutionary message
func NewRevolutionaryMessage(msgType MessageType, fromRole, toRole, payload string) *RevolutionaryMessage {
	return &RevolutionaryMessage{
		messageType: msgType,
		fromRole:    fromRole,
		toRole:      toRole,
		payload:     payload,
		timestamp:   nowFunc(),
	}
}

// NewRegisterMessage creates a registration message from an agent
func NewRegisterMessage(role string) *RevolutionaryMessage {
	return NewRevolutionaryMessage(MessageTypeRegister, role, "", role)
}

// NewActivateMessage creates an activation message to an agent
func NewActivateMessage(toRole, fromRole, payload string) *RevolutionaryMessage {
	return NewRevolutionaryMessage(MessageTypeActivate, fromRole, toRole, payload)
}

// NewQueryAgentsMessage creates a query agents message from the people
func NewQueryAgentsMessage() *RevolutionaryMessage {
	return NewRevolutionaryMessage(MessageTypeQueryAgents, "people", "", "")
}

// NewErrorMessage creates an error message from the soviet
func NewErrorMessage(errorMsg string) *RevolutionaryMessage {
	return NewRevolutionaryMessage(MessageTypeError, "soviet", "", errorMsg)
}

// Type returns the message type
func (m *RevolutionaryMessage) Type() MessageType {
	return m.messageType
}

// FromRole returns the sender role
func (m *RevolutionaryMessage) FromRole() string {
	return m.fromRole
}

// ToRole returns the recipient role
func (m *RevolutionaryMessage) ToRole() string {
	return m.toRole
}

// Payload returns the message payload
func (m *RevolutionaryMessage) Payload() string {
	return m.payload
}

// Timestamp returns when the message was created
func (m *RevolutionaryMessage) Timestamp() time.Time {
	return m.timestamp
}

// IsValid checks if the message is valid
func (m *RevolutionaryMessage) IsValid() bool {
	// Basic validation - message must have a valid type and non-zero timestamp
	validType := m.messageType >= MessageTypeRegister && m.messageType <= MessageTypeAckRegister
	hasTimestamp := !m.timestamp.IsZero()
	return validType && hasTimestamp
}

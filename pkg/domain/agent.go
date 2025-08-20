package domain

import (
	"fmt"
	"time"
)

// AgentState represents the current state of an agent comrade
type AgentState int

const (
	AgentStateWaiting AgentState = iota
	AgentStateWorking
)

// String returns the string representation of AgentState
func (s AgentState) String() string {
	switch s {
	case AgentStateWaiting:
		return "waiting"
	case AgentStateWorking:
		return "working"
	default:
		return "unknown"
	}
}

// AgentComrade represents a worker in the Agent Farm collective.
// Each agent has a role, capabilities, and follows the disciplined lifecycle.
type AgentComrade struct {
	role            string
	agentType       string
	capabilities    []string
	state           AgentState
	connected       bool
	createdAt       time.Time
	lastConnectedAt time.Time
	lastMessage     string
	lastMessageTime time.Time
}

// NewAgentComrade creates a new agent comrade with the specified role, type, and capabilities
func NewAgentComrade(role, agentType string, capabilities []string) *AgentComrade {
	caps := make([]string, len(capabilities))
	copy(caps, capabilities)

	return &AgentComrade{
		role:         role,
		agentType:    agentType,
		capabilities: caps,
		state:        AgentStateWaiting,
		connected:    false,
		createdAt:    nowFunc(),
	}
}

// Role returns the agent's role
func (a *AgentComrade) Role() string {
	return a.role
}

// Type returns the agent's type
func (a *AgentComrade) Type() string {
	return a.agentType
}

// Capabilities returns a copy of the agent's capabilities
func (a *AgentComrade) Capabilities() []string {
	caps := make([]string, len(a.capabilities))
	copy(caps, a.capabilities)
	return caps
}

// State returns the current state of the agent
func (a *AgentComrade) State() AgentState {
	return a.state
}

// IsConnected returns true if the agent is currently connected
func (a *AgentComrade) IsConnected() bool {
	return a.connected
}

// CreatedAt returns when the agent was created
func (a *AgentComrade) CreatedAt() time.Time {
	return a.createdAt
}

// LastConnectedAt returns when the agent was last connected
func (a *AgentComrade) LastConnectedAt() time.Time {
	return a.lastConnectedAt
}

// LastMessage returns the last message received by the agent
func (a *AgentComrade) LastMessage() string {
	return a.lastMessage
}

// LastMessageTime returns when the last message was received
func (a *AgentComrade) LastMessageTime() time.Time {
	return a.lastMessageTime
}

// SetConnected updates the connection state of the agent
func (a *AgentComrade) SetConnected(connected bool) error {
	a.connected = connected
	if connected {
		a.lastConnectedAt = nowFunc()
	}
	return nil
}

// TransitionTo transitions the agent to a new state with validation
func (a *AgentComrade) TransitionTo(newState AgentState) error {
	// Validate state transitions
	if !a.isValidTransition(a.state, newState) {
		return fmt.Errorf("invalid state transition from %s to %s", a.state, newState)
	}

	a.state = newState
	return nil
}

// isValidTransition checks if a state transition is valid
func (a *AgentComrade) isValidTransition(from, to AgentState) bool {
	switch from {
	case AgentStateWaiting:
		return to == AgentStateWorking
	case AgentStateWorking:
		return to == AgentStateWaiting
	default:
		return false
	}
}

// HasCapability checks if the agent has a specific capability
func (a *AgentComrade) HasCapability(capability string) bool {
	for _, cap := range a.capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

// SetLastMessage updates the last message and timestamp
func (a *AgentComrade) SetLastMessage(message string) {
	a.lastMessage = message
	a.lastMessageTime = nowFunc()
}

// Activate transitions the agent from waiting to working state with a message
func (a *AgentComrade) Activate(message string) error {
	if a.state != AgentStateWaiting {
		return fmt.Errorf("cannot activate agent in %s state, must be waiting", a.state)
	}

	a.state = AgentStateWorking
	a.SetLastMessage(message)
	return nil
}

// Yield transitions the agent from working back to waiting state
// This represents the completion of work and voluntary yielding of the barrel
func (a *AgentComrade) Yield() error {
	if a.state != AgentStateWorking {
		return fmt.Errorf("cannot yield while in %s state, must be working", a.state)
	}

	a.state = AgentStateWaiting
	return nil
}

// IsWorking returns true if the agent is currently working
func (a *AgentComrade) IsWorking() bool {
	return a.state == AgentStateWorking
}

// IsWaiting returns true if the agent is currently waiting
func (a *AgentComrade) IsWaiting() bool {
	return a.state == AgentStateWaiting
}

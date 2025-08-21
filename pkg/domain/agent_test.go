package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAgentComrade(t *testing.T) {
	// RED: Test creation of new agent comrade
	agent := NewAgentComrade("developer", []string{"code", "debug", "review"})

	assert.NotNil(t, agent)
	assert.Equal(t, "developer", agent.Role())
	assert.Equal(t, []string{"code", "debug", "review"}, agent.Capabilities())
	assert.Equal(t, AgentStateWaiting, agent.State())
	assert.False(t, agent.IsConnected())
	assert.NotZero(t, agent.CreatedAt())
}

func TestAgentComrade_NewAgentComrade(t *testing.T) {
	// RED: Test creation of new agent comrade
	agent := NewAgentComrade("developer", []string{"code", "debug", "review"})

	assert.NotNil(t, agent)
	assert.Equal(t, "developer", agent.Role())
	assert.Equal(t, []string{"code", "debug", "review"}, agent.Capabilities())
	assert.Equal(t, AgentStateWaiting, agent.State())
	assert.False(t, agent.IsConnected())
	assert.NotZero(t, agent.CreatedAt())
}

func TestAgentComrade_SetConnected(t *testing.T) {
	// RED: Test connection state management
	agent := NewAgentComrade("tester", []string{"test", "validate"})

	// Initially disconnected
	assert.False(t, agent.IsConnected())
	assert.Zero(t, agent.LastConnectedAt())

	// Connect
	agent.SetConnected(true)
	assert.True(t, agent.IsConnected())
	assert.NotZero(t, agent.LastConnectedAt())

	// Disconnect
	agent.SetConnected(false)
	assert.False(t, agent.IsConnected())
}

func TestAgentComrade_TransitionState(t *testing.T) {
	// RED: Test state transitions
	agent := NewAgentComrade("developer", []string{"code"})

	// Initial state is Waiting
	assert.Equal(t, AgentStateWaiting, agent.State())

	// Transition to Working
	err := agent.TransitionTo(AgentStateWorking)
	assert.NoError(t, err)
	assert.Equal(t, AgentStateWorking, agent.State())

	// Back to Waiting
	err = agent.TransitionTo(AgentStateWaiting)
	assert.NoError(t, err)
	assert.Equal(t, AgentStateWaiting, agent.State())
}

func TestAgentComrade_TransitionState_InvalidTransition(t *testing.T) {
	// RED: Test invalid state transitions
	agent := NewAgentComrade("developer", []string{"code"})

	// Cannot transition from Working to Working
	agent.TransitionTo(AgentStateWorking)
	err := agent.TransitionTo(AgentStateWorking)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid state transition")
	assert.Equal(t, AgentStateWorking, agent.State()) // Should remain unchanged
}

func TestAgentComrade_HasCapability(t *testing.T) {
	// RED: Test capability checking
	agent := NewAgentComrade("developer", []string{"code", "debug", "review"})

	assert.True(t, agent.HasCapability("code"))
	assert.True(t, agent.HasCapability("debug"))
	assert.True(t, agent.HasCapability("review"))
	assert.False(t, agent.HasCapability("test"))
	assert.False(t, agent.HasCapability(""))
}

func TestAgentComrade_SetLastMessage(t *testing.T) {
	// RED: Test message tracking
	agent := NewAgentComrade("tester", []string{"test"})

	// Initially no message
	assert.Empty(t, agent.LastMessage())
	assert.Zero(t, agent.LastMessageTime())

	// Set message
	message := "Code ready for testing"
	agent.SetLastMessage(message)

	assert.Equal(t, message, agent.LastMessage())
	assert.NotZero(t, agent.LastMessageTime())
}

func TestAgentState_String(t *testing.T) {
	// RED: Test state string representation
	assert.Equal(t, "waiting", AgentStateWaiting.String())
	assert.Equal(t, "working", AgentStateWorking.String())
}

func TestAgentComrade_Activate(t *testing.T) {
	// RED: Test agent activation
	agent := NewAgentComrade("developer", []string{"code"})

	// Initially waiting
	assert.True(t, agent.IsWaiting())
	assert.False(t, agent.IsWorking())

	// Activate with message
	err := agent.Activate("Start coding feature X")
	assert.NoError(t, err)
	assert.True(t, agent.IsWorking())
	assert.False(t, agent.IsWaiting())
	assert.Equal(t, AgentStateWorking, agent.State())
	assert.Equal(t, "Start coding feature X", agent.LastMessage())
}

func TestAgentComrade_Activate_InvalidState(t *testing.T) {
	// RED: Test activation from invalid state
	agent := NewAgentComrade("developer", []string{"code"})
	agent.TransitionTo(AgentStateWorking)

	// Cannot activate when already working
	err := agent.Activate("Another task")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot activate agent in working state")
}

func TestAgentComrade_Yield(t *testing.T) {
	// RED: Test agent yielding
	agent := NewAgentComrade("developer", []string{"code"})
	agent.Activate("Work on task")

	// Yield back to waiting
	err := agent.Yield()
	assert.NoError(t, err)
	assert.True(t, agent.IsWaiting())
	assert.False(t, agent.IsWorking())
	assert.Equal(t, AgentStateWaiting, agent.State())
}

func TestAgentComrade_Yield_InvalidState(t *testing.T) {
	// RED: Test yield from invalid state
	agent := NewAgentComrade("developer", []string{"code"})

	// Cannot yield when not working
	err := agent.Yield()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot yield while in waiting state")
}

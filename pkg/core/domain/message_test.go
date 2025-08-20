package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRevolutionaryMessage_NewRevolutionaryMessage(t *testing.T) {
	// RED: Test creation of revolutionary message
	msg := NewRevolutionaryMessage(MessageTypeYield, "developer", "tester", "Code ready for testing")

	assert.NotNil(t, msg)
	assert.Equal(t, MessageTypeYield, msg.Type())
	assert.Equal(t, "developer", msg.FromRole())
	assert.Equal(t, "tester", msg.ToRole())
	assert.Equal(t, "Code ready for testing", msg.Payload())
	assert.NotZero(t, msg.Timestamp())
}

func TestRevolutionaryMessage_NewRegisterMessage(t *testing.T) {
	// RED: Test creation of register message
	msg := NewRegisterMessage("developer")

	assert.Equal(t, MessageTypeRegister, msg.Type())
	assert.Equal(t, "developer", msg.FromRole())
	assert.Empty(t, msg.ToRole())
	assert.Equal(t, "developer", msg.Payload())
}

func TestRevolutionaryMessage_NewActivateMessage(t *testing.T) {
	// RED: Test creation of activate message
	msg := NewActivateMessage("tester", "developer", "Code ready for testing")

	assert.Equal(t, MessageTypeActivate, msg.Type())
	assert.Equal(t, "developer", msg.FromRole())
	assert.Equal(t, "tester", msg.ToRole())
	assert.Equal(t, "Code ready for testing", msg.Payload())
}

func TestRevolutionaryMessage_NewQueryAgentsMessage(t *testing.T) {
	// RED: Test creation of query agents message
	msg := NewQueryAgentsMessage()

	assert.Equal(t, MessageTypeQueryAgents, msg.Type())
	assert.Equal(t, "people", msg.FromRole())
	assert.Empty(t, msg.ToRole())
	assert.Empty(t, msg.Payload())
}

func TestRevolutionaryMessage_NewErrorMessage(t *testing.T) {
	// RED: Test creation of error message
	errorMsg := "Invalid operation detected"
	msg := NewErrorMessage(errorMsg)

	assert.Equal(t, MessageTypeError, msg.Type())
	assert.Equal(t, "soviet", msg.FromRole())
	assert.Empty(t, msg.ToRole())
	assert.Equal(t, errorMsg, msg.Payload())
}

func TestRevolutionaryMessage_IsValid(t *testing.T) {
	// RED: Test message validation
	validMsg := NewRevolutionaryMessage(MessageTypeYield, "developer", "tester", "Code ready")
	assert.True(t, validMsg.IsValid())

	// Invalid message with empty type
	invalidMsg := &RevolutionaryMessage{}
	assert.False(t, invalidMsg.IsValid())
}

func TestMessageType_String(t *testing.T) {
	// RED: Test message type string representation
	assert.Equal(t, "REGISTER", MessageTypeRegister.String())
	assert.Equal(t, "YIELD", MessageTypeYield.String())
	assert.Equal(t, "ACTIVATE", MessageTypeActivate.String())
	assert.Equal(t, "QUERY_AGENTS", MessageTypeQueryAgents.String())
	assert.Equal(t, "AGENT_LIST", MessageTypeAgentList.String())
	assert.Equal(t, "ERROR", MessageTypeError.String())
	assert.Equal(t, "ACK_REGISTER", MessageTypeAckRegister.String())
}

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYieldMessage_NewYieldMessage(t *testing.T) {
	// RED: Test creation of yield message
	msg := NewYieldMessage("developer", "tester", "Code ready for testing")

	assert.NotNil(t, msg)
	assert.Equal(t, "developer", msg.FromRole())
	assert.Equal(t, "tester", msg.ToRole())
	assert.Equal(t, "Code ready for testing", msg.Payload())
	assert.NotZero(t, msg.Timestamp())
}

func TestYieldMessage_IsValid(t *testing.T) {
	// RED: Test message validation
	validMsg := NewYieldMessage("developer", "tester", "Code ready")
	assert.True(t, validMsg.IsValid())

	// Invalid message with empty from role
	invalidMsg := NewYieldMessage("", "tester", "Code ready")
	assert.False(t, invalidMsg.IsValid())

	// Invalid message with empty to role
	invalidMsg2 := NewYieldMessage("developer", "", "Code ready")
	assert.False(t, invalidMsg2.IsValid())

	// Invalid message with empty fields (zero value)
	invalidMsg3 := YieldMessage{}
	assert.False(t, invalidMsg3.IsValid())
}

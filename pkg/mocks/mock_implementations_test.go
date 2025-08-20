package mocks

import (
	"testing"

	"github.com/lonegunmanb/agentfarm/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MockImplementationsTestSuite tests all mock implementations for external ports
type MockImplementationsTestSuite struct {
	suite.Suite
}

// TestMockImplementations runs all mock implementation tests
func TestMockImplementations(t *testing.T) {
	suite.Run(t, new(MockImplementationsTestSuite))
}

// TestMockAgentRepository_StoreAndRetrieveAgent tests agent storage functionality
func (suite *MockImplementationsTestSuite) TestMockAgentRepository_StoreAndRetrieveAgent() {
	// RED: This test will fail until we implement MockAgentRepository
	repo := NewMockAgentRepository()

	// Create test agent
	agent := domain.NewAgentComrade("developer", "code-agent", []string{"coding", "testing"})

	// Store agent
	err := repo.Store(agent)
	assert.NoError(suite.T(), err)

	// Retrieve agent
	retrieved, err := repo.GetByRole("developer")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "developer", retrieved.Role())
	assert.Equal(suite.T(), "code-agent", retrieved.Type())
}

// TestMockAgentRepository_GetNonExistentAgent tests error handling
func (suite *MockImplementationsTestSuite) TestMockAgentRepository_GetNonExistentAgent() {
	repo := NewMockAgentRepository()

	// Try to get non-existent agent
	_, err := repo.GetByRole("nonexistent")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
}

// TestMockMessageSender_SendAndCapture tests message sending functionality
func (suite *MockImplementationsTestSuite) TestMockMessageSender_SendAndCapture() {
	// RED: This test will fail until we implement MockMessageSender
	sender := NewMockMessageSender()

	// Send test message
	err := sender.SendActivation("developer", "Start working on authentication module")
	assert.NoError(suite.T(), err)

	// Verify message was captured
	messages := sender.GetSentMessages()
	assert.Len(suite.T(), messages, 1)
	assert.Equal(suite.T(), "developer", messages[0].Recipient)
	assert.Equal(suite.T(), "activation", messages[0].Type)
	assert.Equal(suite.T(), "Start working on authentication module", messages[0].Payload)
}

// TestMockLogger_LogAndCapture tests logging functionality
func (suite *MockImplementationsTestSuite) TestMockLogger_LogAndCapture() {
	// RED: This test will fail until we implement MockLogger
	logger := NewMockLogger()

	// Log test messages
	logger.Info("System started")
	logger.Error("Connection failed", map[string]interface{}{
		"agent": "developer",
		"error": "timeout",
	})

	// Verify logs were captured
	logs := logger.GetLogs()
	assert.Len(suite.T(), logs, 2)
	assert.Equal(suite.T(), "INFO", logs[0].Level)
	assert.Equal(suite.T(), "System started", logs[0].Message)
	assert.Equal(suite.T(), "ERROR", logs[1].Level)
	assert.Equal(suite.T(), "Connection failed", logs[1].Message)
}

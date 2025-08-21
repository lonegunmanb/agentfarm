package tcp

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lonegunmanb/agentfarm/pkg/domain"
)

// MockSovietService for testing
type MockSovietService struct {
	mock.Mock
}

func (m *MockSovietService) RegisterAgent(agent *domain.AgentComrade) (bool, string, error) {
	args := m.Called(agent)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockSovietService) ProcessYield(message domain.YieldMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockSovietService) DeregisterAgent(role string) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockSovietService) QueryStatus() domain.StatusResponse {
	args := m.Called()
	return args.Get(0).(domain.StatusResponse)
}

// MockAgentService for testing
type MockAgentService struct {
	mock.Mock
}

func (m *MockAgentService) GetAgentState(role string) (domain.AgentState, error) {
	args := m.Called(role)
	return args.Get(0).(domain.AgentState), args.Error(1)
}

func (m *MockAgentService) GetBarrelStatus() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAgentService) GetRegisteredAgents() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// MockMessageSender for testing
type MockMessageSender struct {
	mock.Mock
}

func (m *MockMessageSender) SendActivation(role string, payload string) error {
	args := m.Called(role, payload)
	return args.Error(0)
}

// MockLogger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(message string, fields ...map[string]interface{}) {
	m.Called(message, fields)
}

func (m *MockLogger) Error(message string, fields ...map[string]interface{}) {
	m.Called(message, fields)
}

func (m *MockLogger) Debug(message string, fields ...map[string]interface{}) {
	m.Called(message, fields)
}

func (m *MockLogger) Warn(message string, fields ...map[string]interface{}) {
	m.Called(message, fields)
}

func TestTCPServer_HandleRegister(t *testing.T) {
	// Setup
	mockSoviet := &MockSovietService{}
	mockAgent := &MockAgentService{}
	mockSender := &MockMessageSender{}
	mockLogger := &MockLogger{}

	server := NewTCPServer(mockSoviet, mockAgent, mockSender, mockLogger, 0)

	// Test successful registration
	t.Run("successful registration", func(t *testing.T) {
		mockSoviet.On("RegisterAgent", mock.MatchedBy(func(agent *domain.AgentComrade) bool {
			return agent.Role() == "developer" && len(agent.Capabilities()) == 2
		})).Return(false, "", nil).Once()

		shouldActivate, payload, err := server.HandleRegister(context.Background(), "developer", []string{"coding", "testing"})

		assert.NoError(t, err)
		assert.False(t, shouldActivate)
		assert.Empty(t, payload)
		mockSoviet.AssertExpectations(t)
	})

	// Test registration with activation
	t.Run("registration with activation", func(t *testing.T) {
		mockSoviet.On("RegisterAgent", mock.MatchedBy(func(agent *domain.AgentComrade) bool {
			return agent.Role() == "tester"
		})).Return(true, "Start testing", nil).Once()

		shouldActivate, payload, err := server.HandleRegister(context.Background(), "tester", []string{"testing", "automation"})

		assert.NoError(t, err)
		assert.True(t, shouldActivate)
		assert.Equal(t, "Start testing", payload)
		mockSoviet.AssertExpectations(t)
	})
}

func TestTCPServer_HandleYield(t *testing.T) {
	// Setup
	mockSoviet := &MockSovietService{}
	mockAgent := &MockAgentService{}
	mockSender := &MockMessageSender{}
	mockLogger := &MockLogger{}

	server := NewTCPServer(mockSoviet, mockAgent, mockSender, mockLogger, 0)

	// Test successful yield
	t.Run("successful yield", func(t *testing.T) {
		mockSoviet.On("ProcessYield", mock.MatchedBy(func(msg domain.YieldMessage) bool {
			return msg.FromRole() == "developer" && msg.ToRole() == "tester" && msg.Payload() == "Code ready for testing"
		})).Return(nil).Once()

		err := server.HandleYield(context.Background(), "developer", "tester", "Code ready for testing")

		assert.NoError(t, err)
		mockSoviet.AssertExpectations(t)
	})
}

func TestTCPServer_HandleQueryAgents(t *testing.T) {
	// Setup
	mockSoviet := &MockSovietService{}
	mockAgent := &MockAgentService{}
	mockSender := &MockMessageSender{}
	mockLogger := &MockLogger{}

	server := NewTCPServer(mockSoviet, mockAgent, mockSender, mockLogger, 0)

	// Test query agents
	t.Run("query agents", func(t *testing.T) {
		expectedAgents := []string{"developer", "tester"}
		mockAgent.On("GetRegisteredAgents").Return(expectedAgents).Once()

		agents, err := server.HandleQueryAgents(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedAgents, agents)
		mockAgent.AssertExpectations(t)
	})
}

func TestTCPServer_HandleQueryStatus(t *testing.T) {
	// Setup
	mockSoviet := &MockSovietService{}
	mockAgent := &MockAgentService{}
	mockSender := &MockMessageSender{}
	mockLogger := &MockLogger{}

	server := NewTCPServer(mockSoviet, mockAgent, mockSender, mockLogger, 0)

	// Test query status
	t.Run("query status", func(t *testing.T) {
		expectedStatus := domain.StatusResponse{
			BarrelHolder:     "developer",
			RegisteredAgents: []string{"developer", "tester"},
			AgentStates: map[string]domain.AgentState{
				"developer": domain.AgentStateWorking,
				"tester":    domain.AgentStateWaiting,
			},
			ConnectedAgents: map[string]bool{
				"developer": true,
				"tester":    true,
			},
		}

		mockSoviet.On("QueryStatus").Return(expectedStatus).Once()

		status, err := server.HandleQueryStatus(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "developer", status.BarrelHolder)
		assert.Equal(t, []string{"developer", "tester"}, status.RegisteredAgents)
		assert.Equal(t, domain.AgentStateWorking, status.AgentStates["developer"])
		assert.Equal(t, domain.AgentStateWaiting, status.AgentStates["tester"])
		assert.True(t, status.ConnectedAgents["developer"])
		assert.True(t, status.ConnectedAgents["tester"])
		mockSoviet.AssertExpectations(t)
	})
}

func TestTCPMessageSender(t *testing.T) {
	sender := NewTCPMessageSender()

	t.Run("send activation to non-existent connection", func(t *testing.T) {
		err := sender.SendActivation("nonexistent", "test payload")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no connection found")
	})

	t.Run("register and send activation", func(t *testing.T) {
		// Create a test connection using pipe
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		// Register connection
		sender.RegisterConnection("developer", client)

		// Send activation in a goroutine to avoid blocking
		errChan := make(chan error, 1)
		go func() {
			errChan <- sender.SendActivation("developer", "test payload")
		}()

		// Read the message from server side
		buffer := make([]byte, 1024)
		n, err := server.Read(buffer)
		assert.NoError(t, err)

		// Parse the received message
		var msg ActivateMessage
		err = json.Unmarshal(buffer[:n-1], &msg) // -1 to remove the newline
		assert.NoError(t, err)
		assert.Equal(t, "ACTIVATE", msg.Type)
		assert.Equal(t, "test payload", msg.Payload)

		// Check that send was successful
		select {
		case err := <-errChan:
			assert.NoError(t, err)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("SendActivation timed out")
		}
	})

	t.Run("get connected roles", func(t *testing.T) {
		// Create test connections
		_, client1 := net.Pipe()
		_, client2 := net.Pipe()
		defer client1.Close()
		defer client2.Close()

		sender.RegisterConnection("developer", client1)
		sender.RegisterConnection("tester", client2)

		roles := sender.GetConnectedRoles()
		assert.Len(t, roles, 2)
		assert.Contains(t, roles, "developer")
		assert.Contains(t, roles, "tester")

		// Test IsConnected
		assert.True(t, sender.IsConnected("developer"))
		assert.True(t, sender.IsConnected("tester"))
		assert.False(t, sender.IsConnected("nonexistent"))
	})
}

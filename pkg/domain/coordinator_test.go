package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestAgentData represents test agent configuration data
type TestAgentData struct {
	Role         string
	Type         string
	Capabilities []string
}

// Common test agent configurations
var testAgents = map[string]TestAgentData{
	"developer": {
		Role:         "developer",
		Type:         "worker",
		Capabilities: []string{"coding", "testing"},
	},
	"tester": {
		Role:         "tester",
		Type:         "worker",
		Capabilities: []string{"testing", "qa"},
	},
	"reviewer": {
		Role:         "reviewer",
		Type:         "worker",
		Capabilities: []string{"code-review", "mentoring"},
	},
}

// createTestAgent creates an agent from test data
func createTestAgent(agentKey string) *AgentComrade {
	data := testAgents[agentKey]
	return NewAgentComrade(data.Role, data.Capabilities)
}

// CoordinatorTestSuite tests the coordinator functionality in SovietState
type CoordinatorTestSuite struct {
	suite.Suite
	soviet *SovietState
	barrel *BarrelOfGun
}

// SetupTest initializes test dependencies before each test
func (suite *CoordinatorTestSuite) SetupTest() {
	repo := NewMemoryAgentRepository()
	suite.soviet = NewSovietState(repo)
	suite.barrel = NewBarrelOfGun()
	suite.soviet.SetBarrel(suite.barrel)
}

// TestCoordinatorTestSuite runs the test suite
func TestCoordinatorTestSuite(t *testing.T) {
	suite.Run(t, new(CoordinatorTestSuite))
}

// Test_RegisterAgent_SuccessfulRegistration tests agent registration
func (suite *CoordinatorTestSuite) Test_RegisterAgent_SuccessfulRegistration() {
	agent := createTestAgent("developer")

	shouldResume, lastMessage, err := suite.soviet.RegisterAgent(agent)

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), shouldResume) // New agent shouldn't resume work (barrel is with people)
	assert.Empty(suite.T(), lastMessage)
	assert.True(suite.T(), suite.soviet.IsAgentRegistered("developer"))
	assert.True(suite.T(), agent.IsConnected())
	assert.Equal(suite.T(), AgentStateWaiting, agent.State())
}

// Test_RegisterAgent_DuplicateRole_ReplacesExistingAgent tests agent replacement behavior
func (suite *CoordinatorTestSuite) Test_RegisterAgent_DuplicateRole_ReplacesExistingAgent() {
	agent1 := createTestAgent("developer")
	agent2 := createTestAgent("developer") // Same role, different instance

	// First registration should succeed
	_, _, err := suite.soviet.RegisterAgent(agent1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), agent1.IsConnected())

	// Second registration should replace the first
	_, _, err = suite.soviet.RegisterAgent(agent2)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), agent1.IsConnected()) // Original agent disconnected
	assert.True(suite.T(), agent2.IsConnected())  // New agent connected
	assert.True(suite.T(), suite.soviet.IsAgentRegistered("developer"))

	// Only the new agent should be in the soviet
	retrievedAgent := suite.soviet.GetAgent("developer")
	assert.Equal(suite.T(), agent2, retrievedAgent)
	assert.NotEqual(suite.T(), agent1, retrievedAgent)
}

// Test_RegisterAgent_WithBarrel_ShouldResume tests reconnection with barrel
func (suite *CoordinatorTestSuite) Test_RegisterAgent_WithBarrel_ShouldResume() {
	agent := createTestAgent("developer")

	// Transfer barrel to developer role first
	suite.barrel.TransferTo("developer", "Test message")

	// Now register the agent
	shouldResume, lastMessage, err := suite.soviet.RegisterAgent(agent)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), shouldResume)                      // Agent should resume work
	assert.Equal(suite.T(), "Test message", lastMessage)      // Should get the last message
	assert.Equal(suite.T(), AgentStateWorking, agent.State()) // Should be working
}

// Test_DeregisterAgent_WithBarrel_ReturnsTopeople tests deregistration with barrel transfer
func (suite *CoordinatorTestSuite) Test_DeregisterAgent_WithBarrel_ReturnsTopeople() {
	agent := createTestAgent("developer")

	// Register and give barrel to agent
	suite.soviet.RegisterAgent(agent)
	suite.barrel.TransferTo("developer", "Working")

	// Deregister the agent
	err := suite.soviet.DeregisterAgent("developer")

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), suite.soviet.IsAgentRegistered("developer"))
	assert.Equal(suite.T(), "people", suite.barrel.CurrentHolder()) // Barrel should go back to people
}

// Test_ProcessYield_ValidTransfer tests yield processing
func (suite *CoordinatorTestSuite) Test_ProcessYield_ValidTransfer() {
	fromAgent := createTestAgent("developer")
	toAgent := createTestAgent("tester")

	// Register both agents
	suite.soviet.RegisterAgent(fromAgent)
	suite.soviet.RegisterAgent(toAgent)

	// Give barrel to from agent
	suite.barrel.TransferTo("developer", "Initial")
	fromAgent.TransitionTo(AgentStateWorking)

	// Create yield message
	message := NewYieldMessage("developer", "tester", "Code ready for testing")

	// Process yield
	err := suite.soviet.ProcessYield(message)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "tester", suite.barrel.CurrentHolder())
	assert.Equal(suite.T(), AgentStateWaiting, fromAgent.State())
	assert.Equal(suite.T(), AgentStateWorking, toAgent.State())
}

// Test_GetAgentState tests agent state retrieval
func (suite *CoordinatorTestSuite) Test_GetAgentState() {
	agent := createTestAgent("developer")
	suite.soviet.RegisterAgent(agent)

	state, err := suite.soviet.GetAgentState("developer")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), AgentStateWaiting, state)

	// Test non-existent agent
	_, err = suite.soviet.GetAgentState("nonexistent")
	assert.Error(suite.T(), err)
}

// Test_GetBarrelStatus tests barrel status retrieval
func (suite *CoordinatorTestSuite) Test_GetBarrelStatus() {
	// Initially should be with people
	status := suite.soviet.GetBarrelStatus()
	assert.Equal(suite.T(), "people", status)

	// Transfer to an agent
	suite.barrel.TransferTo("developer", "Test")
	status = suite.soviet.GetBarrelStatus()
	assert.Equal(suite.T(), "developer", status)
}

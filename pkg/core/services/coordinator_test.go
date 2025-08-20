package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/lonegunmanb/agentfarm/pkg/core/domain"
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
	"minimal": {
		Role:         "minimal",
		Type:         "worker",
		Capabilities: []string{"basic"},
	},
}

// createTestAgent creates an agent from test data
func createTestAgent(agentKey string) *domain.AgentComrade {
	data := testAgents[agentKey]
	return domain.NewAgentComrade(data.Role, data.Type, data.Capabilities)
}

// SovietCoordinatorTestSuite tests the SovietCoordinator service
type SovietCoordinatorTestSuite struct {
	suite.Suite
	coordinator *SovietCoordinator
	soviet      *domain.SovietState
	barrel      *domain.BarrelOfGun
}

// SetupTest initializes test dependencies before each test
func (suite *SovietCoordinatorTestSuite) SetupTest() {
	suite.soviet = domain.NewSovietState()
	suite.barrel = domain.NewBarrelOfGun()
	suite.soviet.SetBarrel(suite.barrel)
	suite.coordinator = NewSovietCoordinator(suite.soviet)
}

// TestSovietCoordinatorTestSuite runs the test suite
func TestSovietCoordinatorTestSuite(t *testing.T) {
	suite.Run(t, new(SovietCoordinatorTestSuite))
}

// Test_NewSovietCoordinator_CreatesValidInstance tests coordinator creation
func (suite *SovietCoordinatorTestSuite) Test_NewSovietCoordinator_CreatesValidInstance() {
	coordinator := NewSovietCoordinator(suite.soviet)

	assert.NotNil(suite.T(), coordinator)
	assert.Equal(suite.T(), suite.soviet, coordinator.soviet)
}

// Test_RegisterAgent_SuccessfulRegistration tests agent registration
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_SuccessfulRegistration() {
	agent := createTestAgent("developer")

	shouldResume, lastMessage, err := suite.coordinator.RegisterAgent(agent)

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), shouldResume) // New agent shouldn't resume work (barrel is with people)
	assert.Empty(suite.T(), lastMessage)
	assert.True(suite.T(), suite.soviet.IsAgentRegistered("developer"))
	assert.True(suite.T(), agent.IsConnected())
	assert.Equal(suite.T(), domain.AgentStateWaiting, agent.State())
}

// Test_RegisterAgent_DuplicateRole_ReplacesExistingAgent tests agent replacement behavior
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_DuplicateRole_ReplacesExistingAgent() {
	agent1 := createTestAgent("developer")
	agent2 := createTestAgent("developer") // Same role, different instance

	// First registration should succeed
	_, _, err := suite.coordinator.RegisterAgent(agent1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), agent1.IsConnected())

	// Second registration with same role should succeed and replace the first
	_, _, err = suite.coordinator.RegisterAgent(agent2)
	assert.NoError(suite.T(), err)

	// The new agent should be connected and registered
	assert.True(suite.T(), agent2.IsConnected())
	assert.True(suite.T(), suite.soviet.IsAgentRegistered("developer"))

	// The old agent should be disconnected (replaced)
	assert.False(suite.T(), agent1.IsConnected())

	// Only the new agent should be in the registry
	registeredAgent := suite.soviet.GetAgent("developer")
	assert.Equal(suite.T(), agent2, registeredAgent)
}

// Test_RegisterAgent_ReplacementWithBarrelTransfer_HandlesBarrelCorrectly tests barrel handling during replacement
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_ReplacementWithBarrelTransfer_HandlesBarrelCorrectly() {
	agent1 := createTestAgent("developer")

	// Initial registration and barrel assignment
	_, _, err := suite.coordinator.RegisterAgent(agent1)
	assert.NoError(suite.T(), err)
	suite.barrel.TransferTo("developer", "Working on feature")
	agent1.TransitionTo(domain.AgentStateWorking)

	// Verify initial state
	assert.Equal(suite.T(), "developer", suite.barrel.CurrentHolder())
	assert.Equal(suite.T(), domain.AgentStateWorking, agent1.State())

	// New agent with same role registers (reconnection/replacement scenario)
	agent2 := createTestAgent("developer")
	shouldResume, lastMessage, err := suite.coordinator.RegisterAgent(agent2)
	assert.NoError(suite.T(), err)

	// Since the "developer" role holds the barrel, the new agent should resume work
	assert.True(suite.T(), shouldResume)
	assert.Equal(suite.T(), "Working on feature", lastMessage)

	// Old agent should be disconnected
	assert.False(suite.T(), agent1.IsConnected())

	// New agent should be registered, connected, and working
	assert.True(suite.T(), agent2.IsConnected())
	assert.Equal(suite.T(), domain.AgentStateWorking, agent2.State())
	assert.Equal(suite.T(), agent2, suite.soviet.GetAgent("developer"))

	// Barrel should still be held by "developer" role (but new agent instance)
	assert.Equal(suite.T(), "developer", suite.barrel.CurrentHolder())

	// The agent is already working due to unified registration handling the reconnection automatically
	assert.Equal(suite.T(), domain.AgentStateWorking, agent2.State())
}

// Test_RegisterAgent_DifferentRoles_BothSucceed tests that different roles can coexist
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_DifferentRoles_BothSucceed() {
	developer := createTestAgent("developer")
	tester := createTestAgent("tester")

	// Both registrations should succeed (different roles)
	_, _, err := suite.coordinator.RegisterAgent(developer)
	assert.NoError(suite.T(), err)

	_, _, err = suite.coordinator.RegisterAgent(tester)
	assert.NoError(suite.T(), err)

	// Both agents should be connected and registered
	assert.True(suite.T(), developer.IsConnected())
	assert.True(suite.T(), tester.IsConnected())
	assert.True(suite.T(), suite.soviet.IsAgentRegistered("developer"))
	assert.True(suite.T(), suite.soviet.IsAgentRegistered("tester"))

	// Both should be in the registered agents list
	agents := suite.coordinator.GetRegisteredAgents()
	assert.Len(suite.T(), agents, 2)
	assert.Contains(suite.T(), agents, "developer")
	assert.Contains(suite.T(), agents, "tester")
}

// Test_RegisterAgent_NilAgent_ReturnsError tests nil agent registration
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_NilAgent_ReturnsError() {
	_, _, err := suite.coordinator.RegisterAgent(nil)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "agent cannot be nil")
}

// Test_DeregisterAgent_SuccessfulDeregistration tests agent deregistration
func (suite *SovietCoordinatorTestSuite) Test_DeregisterAgent_SuccessfulDeregistration() {
	agent := createTestAgent("developer")
	suite.coordinator.RegisterAgent(agent)

	err := suite.coordinator.DeregisterAgent("developer")

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), suite.soviet.IsAgentRegistered("developer"))
}

// Test_DeregisterAgent_NonExistentRole_ReturnsError tests deregistration of non-existent agent
func (suite *SovietCoordinatorTestSuite) Test_DeregisterAgent_NonExistentRole_ReturnsError() {
	err := suite.coordinator.DeregisterAgent("nonexistent")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "agent with role 'nonexistent' not found")
}

// Test_DeregisterAgent_AgentHoldsBarrel_TransfersToCommonOwnership tests deregistration of barrel holder
func (suite *SovietCoordinatorTestSuite) Test_DeregisterAgent_AgentHoldsBarrel_TransfersToCommonOwnership() {
	agent := createTestAgent("developer")
	suite.coordinator.RegisterAgent(agent)

	// Transfer barrel to agent
	suite.barrel.TransferTo("developer", "Test transfer")
	assert.Equal(suite.T(), "developer", suite.barrel.CurrentHolder())

	// Deregister agent holding barrel
	err := suite.coordinator.DeregisterAgent("developer")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "people", suite.barrel.CurrentHolder()) // Should return to people
	assert.False(suite.T(), suite.soviet.IsAgentRegistered("developer"))
}

// Test_ProcessYield_ValidYield_SuccessfulTransfer tests valid yield processing
func (suite *SovietCoordinatorTestSuite) Test_ProcessYield_ValidYield_SuccessfulTransfer() {
	// Register agents
	developer := createTestAgent("developer")
	tester := createTestAgent("tester")
	suite.coordinator.RegisterAgent(developer)
	suite.coordinator.RegisterAgent(tester)

	// Transfer barrel to developer
	suite.barrel.TransferTo("developer", "Initial transfer")
	developer.TransitionTo(domain.AgentStateWorking)

	// Create yield message
	message := domain.NewYieldMessage(
		"developer",
		"tester",
		"Code ready for testing",
	)

	err := suite.coordinator.ProcessYield(message)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "tester", suite.barrel.CurrentHolder())
	assert.Equal(suite.T(), domain.AgentStateWaiting, developer.State()) // Developer should be waiting
	assert.Equal(suite.T(), domain.AgentStateWorking, tester.State())    // Tester should be working
}

// Test_ProcessYield_InvalidSender_ReturnsError tests yield from non-barrel holder
func (suite *SovietCoordinatorTestSuite) Test_ProcessYield_InvalidSender_ReturnsError() {
	// Register agents
	developer := createTestAgent("developer")
	tester := createTestAgent("tester")
	suite.coordinator.RegisterAgent(developer)
	suite.coordinator.RegisterAgent(tester)

	// Barrel is held by "people", not "developer"
	assert.Equal(suite.T(), "people", suite.barrel.CurrentHolder())

	// Try to yield from developer (who doesn't hold barrel)
	message := domain.NewYieldMessage(
		"developer",
		"tester",
		"Unauthorized yield attempt",
	)

	err := suite.coordinator.ProcessYield(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "only current barrel holder can yield")
}

// Test_ProcessYield_NonExistentTarget_ReturnsError tests yield to non-existent agent
func (suite *SovietCoordinatorTestSuite) Test_ProcessYield_NonExistentTarget_ReturnsError() {
	// Register source agent
	developer := createTestAgent("developer")
	suite.coordinator.RegisterAgent(developer)

	// Transfer barrel to developer
	suite.barrel.TransferTo("developer", "Initial transfer")
	developer.TransitionTo(domain.AgentStateWorking)

	// Try to yield to non-existent agent
	message := domain.NewYieldMessage(
		"developer",
		"nonexistent",
		"Yield to nowhere",
	)

	err := suite.coordinator.ProcessYield(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "target agent 'nonexistent' not found")
}

// Test_ProcessYield_YieldToPeople_SuccessfulTransfer tests yielding to people
func (suite *SovietCoordinatorTestSuite) Test_ProcessYield_YieldToPeople_SuccessfulTransfer() {
	// Register agent
	developer := createTestAgent("developer")
	suite.coordinator.RegisterAgent(developer)

	// Transfer barrel to developer
	suite.barrel.TransferTo("developer", "Initial transfer")
	developer.TransitionTo(domain.AgentStateWorking)

	// Yield to people
	message := domain.NewYieldMessage(
		"developer",
		"people",
		"Task completed, returning to people",
	)

	err := suite.coordinator.ProcessYield(message)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "people", suite.barrel.CurrentHolder())
	assert.Equal(suite.T(), domain.AgentStateWaiting, developer.State())
}

// Test_GetAgentState_RegisteredAgent_ReturnsState tests getting agent state
func (suite *SovietCoordinatorTestSuite) Test_GetAgentState_RegisteredAgent_ReturnsState() {
	agent := createTestAgent("developer")
	suite.coordinator.RegisterAgent(agent)

	state, err := suite.coordinator.GetAgentState("developer")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.AgentStateWaiting, state)
}

// Test_GetAgentState_NonExistentAgent_ReturnsError tests getting state of non-existent agent
func (suite *SovietCoordinatorTestSuite) Test_GetAgentState_NonExistentAgent_ReturnsError() {
	_, err := suite.coordinator.GetAgentState("nonexistent")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "agent with role 'nonexistent' not found")
}

// Test_GetBarrelStatus_ReturnsCurrentStatus tests getting barrel status
func (suite *SovietCoordinatorTestSuite) Test_GetBarrelStatus_ReturnsCurrentStatus() {
	holder := suite.coordinator.GetBarrelStatus()

	assert.Equal(suite.T(), "people", holder) // Initial holder is people
}

// Test_GetRegisteredAgents_ReturnsAgentList tests getting registered agents
func (suite *SovietCoordinatorTestSuite) Test_GetRegisteredAgents_ReturnsAgentList() {
	// Register some agents
	developer := createTestAgent("developer")
	tester := createTestAgent("tester")
	suite.coordinator.RegisterAgent(developer)
	suite.coordinator.RegisterAgent(tester)

	agents := suite.coordinator.GetRegisteredAgents()

	assert.Len(suite.T(), agents, 2)
	assert.Contains(suite.T(), agents, "developer")
	assert.Contains(suite.T(), agents, "tester")
}

// Test_RegisterAgent_BarrelHolderReconnects_ResumesWork tests unified registration handling reconnection of barrel holder
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_BarrelHolderReconnects_ResumesWork() {
	// Register agent
	developer := createTestAgent("developer")
	_, _, _ = suite.coordinator.RegisterAgent(developer)

	// Transfer barrel to developer and make them work
	suite.barrel.TransferTo("developer", "Work assignment")
	developer.TransitionTo(domain.AgentStateWorking)

	// Simulate disconnection (agent instance lost/crashed)
	developer.SetConnected(false)

	// Simulate reconnection: new agent instance with same role registers
	reconnectedDeveloper := createTestAgent("developer")
	shouldResume, lastMessage, err := suite.coordinator.RegisterAgent(reconnectedDeveloper)

	// The registration should handle reconnection automatically
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), shouldResume) // Should resume because "developer" role holds barrel
	assert.Equal(suite.T(), "Work assignment", lastMessage)
	assert.True(suite.T(), reconnectedDeveloper.IsConnected())
	assert.Equal(suite.T(), domain.AgentStateWorking, reconnectedDeveloper.State())
	
	// Old agent should be disconnected, new agent should be in registry
	assert.False(suite.T(), developer.IsConnected())
	assert.Equal(suite.T(), reconnectedDeveloper, suite.soviet.GetAgent("developer"))
}

// Test_RegisterAgent_ReconnectionWorkflow_ProperSequence tests the correct reconnection workflow
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_ReconnectionWorkflow_ProperSequence() {
	// Step 1: Agent registers and gets barrel
	developer := createTestAgent("developer")
	_, _, _ = suite.coordinator.RegisterAgent(developer)
	suite.barrel.TransferTo("developer", "Initial work")
	developer.TransitionTo(domain.AgentStateWorking)

	// Step 2: Agent disconnects unexpectedly (network failure, crash, etc.)
	developer.SetConnected(false)

	// Step 3: Agent process restarts and registers again (same role)
	// This should replace the disconnected agent instance and automatically handle reconnection
	reconnectedAgent := createTestAgent("developer")
	shouldResume, lastMessage, err := suite.coordinator.RegisterAgent(reconnectedAgent)

	// Step 4: Registration should succeed and replace the old agent, automatically resuming work
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), shouldResume) // Should resume because "developer" role held the barrel
	assert.Equal(suite.T(), "Initial work", lastMessage)
	assert.True(suite.T(), reconnectedAgent.IsConnected())
	assert.False(suite.T(), developer.IsConnected()) // Old agent should be disconnected
	assert.Equal(suite.T(), domain.AgentStateWorking, reconnectedAgent.State())
}

// Test_RegisterAgent_MultipleReplacements_OnlyLatestAgentActive tests multiple rapid replacements
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_MultipleReplacements_OnlyLatestAgentActive() {
	// Register first agent
	agent1 := createTestAgent("developer")
	_, _, err := suite.coordinator.RegisterAgent(agent1)
	assert.NoError(suite.T(), err)

	// Register second agent (should replace first)
	agent2 := createTestAgent("developer")
	_, _, err = suite.coordinator.RegisterAgent(agent2)
	assert.NoError(suite.T(), err)

	// Register third agent (should replace second)
	agent3 := createTestAgent("developer")
	_, _, err = suite.coordinator.RegisterAgent(agent3)
	assert.NoError(suite.T(), err)

	// Only the latest agent should be connected and registered
	assert.False(suite.T(), agent1.IsConnected())
	assert.False(suite.T(), agent2.IsConnected())
	assert.True(suite.T(), agent3.IsConnected())
	assert.Equal(suite.T(), agent3, suite.soviet.GetAgent("developer"))
}

// Test_RegisterAgent_NonBarrelHolderReconnects_NoResume tests unified registration for non-barrel holder
func (suite *SovietCoordinatorTestSuite) Test_RegisterAgent_NonBarrelHolderReconnects_NoResume() {
	// Register agent
	developer := createTestAgent("developer")
	_, _, _ = suite.coordinator.RegisterAgent(developer)

	// Barrel remains with people, developer is just waiting
	assert.Equal(suite.T(), "people", suite.barrel.CurrentHolder())

	// Simulate disconnection
	developer.SetConnected(false)

	// Simulate reconnection: new agent instance registers with same role
	reconnectedDeveloper := createTestAgent("developer")
	shouldResume, lastMessage, err := suite.coordinator.RegisterAgent(reconnectedDeveloper)

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), shouldResume) // Should not resume (barrel is with people)
	assert.Empty(suite.T(), lastMessage)
	assert.True(suite.T(), reconnectedDeveloper.IsConnected())
	assert.Equal(suite.T(), domain.AgentStateWaiting, reconnectedDeveloper.State())
	
	// Old agent should be disconnected, new agent should be in registry
	assert.False(suite.T(), developer.IsConnected())
	assert.Equal(suite.T(), reconnectedDeveloper, suite.soviet.GetAgent("developer"))
}

package mocks

import (
	"fmt"
	"testing"

	"github.com/lonegunmanb/agentfarm/pkg/core/domain"
	"github.com/lonegunmanb/agentfarm/pkg/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// WorkflowIntegrationTestSuite tests complete revolutionary workflows using mock implementations
type WorkflowIntegrationTestSuite struct {
	suite.Suite
	soviet        *domain.SovietState
	coordinator   *services.SovietCoordinator
	sovietService *CoordinatorAdapter
	mockRepo      *MockAgentRepository
	mockSender    *MockMessageSender
	mockLogger    *MockLogger
}

// SetupTest sets up test environment for each test
func (suite *WorkflowIntegrationTestSuite) SetupTest() {
	// Create mock implementations
	suite.mockRepo = NewMockAgentRepository()
	suite.mockSender = NewMockMessageSender()
	suite.mockLogger = NewMockLogger()

	// Create pure domain objects
	suite.soviet = domain.NewSovietState()

	// Initialize barrel (required for tests)
	barrel := domain.NewBarrelOfGun()
	err := suite.soviet.SetBarrel(barrel)
	if err != nil {
		panic(fmt.Sprintf("Failed to set barrel: %v", err))
	}

	// Create coordinator with dependencies injected
	suite.coordinator = services.NewSovietCoordinatorWithDependencies(
		suite.soviet,
		suite.mockRepo,
		suite.mockSender,
		suite.mockLogger,
	)
	suite.sovietService = NewCoordinatorAdapter(suite.coordinator, suite.soviet)
}

// TestWorkflowIntegrationTests runs all workflow integration tests
func TestWorkflowIntegrationTests(t *testing.T) {
	suite.Run(t, new(WorkflowIntegrationTestSuite))
}

// TestCompleteRevolutionaryWorkflow tests the complete agent registration -> yield -> transfer cycle
func (suite *WorkflowIntegrationTestSuite) TestCompleteRevolutionaryWorkflow() {
	// Phase 1: Register developer agent (SovietState now handles all external operations)
	developerAgent := domain.NewAgentComrade("developer", "code-agent", []string{"coding", "testing"})

	shouldResume, lastMessage, err := suite.sovietService.RegisterAgent(developerAgent)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), shouldResume) // Should not resume initially (barrel with people)
	assert.Empty(suite.T(), lastMessage)

	// Phase 2: People yield barrel to developer (SovietState handles messaging and events)
	yieldToDeveloper := domain.NewYieldMessage("people", "developer", "Implement authentication module")
	err = suite.sovietService.ProcessYield(yieldToDeveloper)
	assert.NoError(suite.T(), err)

	// Verify developer is now working
	developerState, err := suite.coordinator.GetAgentState("developer")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.AgentStateWorking, developerState)

	// Phase 3: Register tester agent while developer is working
	testerAgent := domain.NewAgentComrade("tester", "test-agent", []string{"testing", "validation"})

	shouldResume, lastMessage, err = suite.sovietService.RegisterAgent(testerAgent)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), shouldResume) // Should not resume (developer holds barrel)
	assert.Empty(suite.T(), lastMessage)

	// Phase 4: Developer yields barrel to tester (SovietState handles all external operations)
	yieldToTester := domain.NewYieldMessage("developer", "tester", "Test the authentication module")
	err = suite.sovietService.ProcessYield(yieldToTester)
	assert.NoError(suite.T(), err)

	// Verify states
	developerState, err = suite.coordinator.GetAgentState("developer")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.AgentStateWaiting, developerState)

	testerState, err := suite.coordinator.GetAgentState("tester")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.AgentStateWorking, testerState)

	// Phase 5: Tester yields barrel back to people
	yieldToPeople := domain.NewYieldMessage("tester", "people", "Testing completed successfully")
	err = suite.sovietService.ProcessYield(yieldToPeople)
	assert.NoError(suite.T(), err)

	// Verify final state
	barrelHolder := suite.coordinator.GetBarrelStatus()
	assert.Equal(suite.T(), "people", barrelHolder)

	testerState, err = suite.coordinator.GetAgentState("tester")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.AgentStateWaiting, testerState)
}

// TestPeoplesInterventionAndStatusQuery tests People's intervention and status querying
func (suite *WorkflowIntegrationTestSuite) TestPeoplesInterventionAndStatusQuery() {
	// Register multiple agents
	developerAgent := domain.NewAgentComrade("developer", "code-agent", []string{"coding"})
	testerAgent := domain.NewAgentComrade("tester", "test-agent", []string{"testing"})

	_, _, err := suite.sovietService.RegisterAgent(developerAgent)
	assert.NoError(suite.T(), err)

	_, _, err = suite.sovietService.RegisterAgent(testerAgent)
	assert.NoError(suite.T(), err)

	// People yield to developer
	yieldToDeveloper := domain.NewYieldMessage("people", "developer", "Start development phase")
	err = suite.sovietService.ProcessYield(yieldToDeveloper)
	assert.NoError(suite.T(), err)

	// Query status (simulating People's representative checking system)
	status := suite.sovietService.QueryStatus()

	assert.Equal(suite.T(), "developer", status.BarrelHolder)
	assert.Len(suite.T(), status.RegisteredAgents, 2)
	assert.Contains(suite.T(), status.RegisteredAgents, "developer")
	assert.Contains(suite.T(), status.RegisteredAgents, "tester")
	assert.Equal(suite.T(), domain.AgentStateWorking, status.AgentStates["developer"])
	assert.Equal(suite.T(), domain.AgentStateWaiting, status.AgentStates["tester"])

	// People intervene and take back the barrel
	yieldToPeople := domain.NewYieldMessage("developer", "people", "People's intervention required")
	err = suite.sovietService.ProcessYield(yieldToPeople)
	assert.NoError(suite.T(), err)

	// Verify intervention succeeded
	barrelHolder := suite.coordinator.GetBarrelStatus()
	assert.Equal(suite.T(), "people", barrelHolder)

	developerState, err := suite.coordinator.GetAgentState("developer")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.AgentStateWaiting, developerState)
}

// TestDisconnectionRecoveryWithMocks tests disconnection/reconnection recovery with mock failure simulation
func (suite *WorkflowIntegrationTestSuite) TestDisconnectionRecoveryWithMocks() {
	// Phase 1: Register developer and give them the barrel
	developerAgent := domain.NewAgentComrade("developer", "code-agent", []string{"coding"})

	_, _, err := suite.sovietService.RegisterAgent(developerAgent)
	assert.NoError(suite.T(), err)

	// People yield to developer
	yieldToDeveloper := domain.NewYieldMessage("people", "developer", "Work on critical feature")
	err = suite.sovietService.ProcessYield(yieldToDeveloper)
	assert.NoError(suite.T(), err)

	// Verify developer is working
	developerState, err := suite.coordinator.GetAgentState("developer")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.AgentStateWorking, developerState)

	// Phase 2: Simulate reconnection - developer reconnects
	newDeveloperAgent := domain.NewAgentComrade("developer", "code-agent", []string{"coding"})

	shouldResume, lastMessage, err := suite.sovietService.RegisterAgent(newDeveloperAgent)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), shouldResume) // Should resume work since they hold the barrel
	assert.Equal(suite.T(), "Work on critical feature", lastMessage)

	// Verify developer resumed working state
	developerState, err = suite.coordinator.GetAgentState("developer")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.AgentStateWorking, developerState)

	// Verify barrel status
	barrelHolder := suite.coordinator.GetBarrelStatus()
	assert.Equal(suite.T(), "developer", barrelHolder)
}

// TestMockVerificationAndAssertion tests that all mocks captured expected interactions
func (suite *WorkflowIntegrationTestSuite) TestMockVerificationAndAssertion() {
	// Register an agent and perform a complete workflow (SovietState handles all external operations)
	developerAgent := domain.NewAgentComrade("developer", "code-agent", []string{"coding"})

	_, _, err := suite.sovietService.RegisterAgent(developerAgent)
	assert.NoError(suite.T(), err)

	// Yield barrel to agent (triggers messaging and events)
	yieldMessage := domain.NewYieldMessage("people", "developer", "Start coding")
	err = suite.sovietService.ProcessYield(yieldMessage)
	assert.NoError(suite.T(), err)

	// Verify all mock interactions captured by SovietState

	// Repository verification - agent should be stored automatically
	storedAgent, err := suite.mockRepo.GetByRole("developer")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "developer", storedAgent.Role())
	assert.True(suite.T(), suite.mockRepo.Exists("developer"))

	// Message sender verification - activation should be sent automatically
	sentMessages := suite.mockSender.GetSentMessages()
	assert.GreaterOrEqual(suite.T(), len(sentMessages), 1)

	// Find activation message
	var activationFound bool
	for _, msg := range sentMessages {
		if msg.Type == "activation" && msg.Recipient == "developer" {
			activationFound = true
			assert.Equal(suite.T(), "Start coding", msg.Payload)
			break
		}
	}
	assert.True(suite.T(), activationFound, "Activation message should be sent automatically")

	// Logger verification - logs should be created automatically
	logs := suite.mockLogger.GetLogs()
	assert.GreaterOrEqual(suite.T(), len(logs), 1)

	// Find info logs (successful operations)
	infoLogs := suite.mockLogger.GetLogsByLevel("INFO")
	assert.GreaterOrEqual(suite.T(), len(infoLogs), 1)
}

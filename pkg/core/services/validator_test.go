package services

import (
	"testing"

	"github.com/lonegunmanb/agentfarm/pkg/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ProtocolValidatorTestSuite provides comprehensive test coverage for protocol validation
type ProtocolValidatorTestSuite struct {
	suite.Suite
	validator  *ProtocolValidator
	soviet     *domain.SovietState
	testAgents map[string]*domain.AgentComrade
	testBarrel *domain.BarrelOfGun
}

func (suite *ProtocolValidatorTestSuite) SetupTest() {
	// Create test soviet state
	suite.soviet = domain.NewSovietState()

	// Create validator
	suite.validator = NewProtocolValidator(suite.soviet)

	// Create test barrel
	suite.testBarrel = domain.NewBarrelOfGun()
	suite.soviet.SetBarrel(suite.testBarrel)

	// Create test agents
	suite.testAgents = map[string]*domain.AgentComrade{
		"developer": domain.NewAgentComrade("developer", "developer", []string{"code", "test"}),
		"tester":    domain.NewAgentComrade("tester", "tester", []string{"test", "validate"}),
		"reviewer":  domain.NewAgentComrade("reviewer", "reviewer", []string{"review", "approve"}),
	}

	// Set agents as connected (they start as disconnected by default)
	for _, agent := range suite.testAgents {
		agent.SetConnected(true)
	}

	// Register test agents
	for _, agent := range suite.testAgents {
		err := suite.soviet.RegisterAgent(agent)
		suite.Require().NoError(err)
	}
}

// Test ValidateYieldMessage - Message Structure Validation
func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_ValidMessage() {
	message := domain.NewYieldMessage("developer", "tester", "Code ready for testing")

	err := suite.validator.ValidateYieldMessage(message)

	assert.NoError(suite.T(), err)
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_EmptyFromRole() {
	message := domain.NewYieldMessage("", "tester", "Code ready for testing")

	err := suite.validator.ValidateYieldMessage(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "from_role cannot be empty")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_EmptyToRole() {
	message := domain.NewYieldMessage("developer", "", "Code ready for testing")

	err := suite.validator.ValidateYieldMessage(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "to_role cannot be empty")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_SelfYield() {
	message := domain.NewYieldMessage("developer", "developer", "Self-yield not allowed")

	err := suite.validator.ValidateYieldMessage(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "agent cannot yield to itself")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_InvalidMessage() {
	// Create a message with roles but invalid timestamp (manually construct to bypass NewYieldMessage validation)
	// This requires creating a message that has from_role and to_role but fails IsValid due to timestamp
	message := domain.YieldMessage{} // Zero value has empty timestamp
	// We can't access private fields directly, so this test needs to be adjusted
	// Instead, test with an empty message that will fail on from_role check first

	err := suite.validator.ValidateYieldMessage(message)

	assert.Error(suite.T(), err)
	// Since empty message fails on from_role first, we expect that error
	assert.Contains(suite.T(), err.Error(), "from_role cannot be empty")
}

// Test ValidateBarrelHolderRights - Revolutionary Rule Enforcement
func (suite *ProtocolValidatorTestSuite) TestValidateBarrelHolderRights_ValidHolder() {
	// Set barrel holder
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)

	err = suite.validator.ValidateBarrelHolderRights("developer")

	assert.NoError(suite.T(), err)
}

func (suite *ProtocolValidatorTestSuite) TestValidateBarrelHolderRights_NotCurrentHolder() {
	// Set barrel holder to developer
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)

	// Try to validate rights for tester
	err = suite.validator.ValidateBarrelHolderRights("tester")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "only current barrel holder can yield")
	assert.Contains(suite.T(), err.Error(), "current holder: developer")
	assert.Contains(suite.T(), err.Error(), "requester: tester")
}

func (suite *ProtocolValidatorTestSuite) TestValidateBarrelHolderRights_PeopleCanAlwaysYield() {
	// Set barrel holder to developer
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)

	// People should always be able to yield regardless of current holder
	err = suite.validator.ValidateBarrelHolderRights("people")

	assert.NoError(suite.T(), err)
}

func (suite *ProtocolValidatorTestSuite) TestValidateBarrelHolderRights_NoBarrel() {
	// Create a new soviet without a barrel to test edge case
	emptyValidator := NewProtocolValidator(domain.NewSovietState())

	err := emptyValidator.ValidateBarrelHolderRights("developer")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "no barrel available")
}

// Test ValidateTargetAgent - Target Validation
func (suite *ProtocolValidatorTestSuite) TestValidateTargetAgent_ValidAgent() {
	err := suite.validator.ValidateTargetAgent("tester")

	assert.NoError(suite.T(), err)
}

func (suite *ProtocolValidatorTestSuite) TestValidateTargetAgent_AgentNotFound() {
	err := suite.validator.ValidateTargetAgent("nonexistent")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "target agent 'nonexistent' not found")
}

func (suite *ProtocolValidatorTestSuite) TestValidateTargetAgent_PeopleAlwaysValid() {
	err := suite.validator.ValidateTargetAgent("people")

	assert.NoError(suite.T(), err)
}

func (suite *ProtocolValidatorTestSuite) TestValidateTargetAgent_DisconnectedAgent() {
	// Disconnect the tester agent
	suite.testAgents["tester"].SetConnected(false)

	err := suite.validator.ValidateTargetAgent("tester")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "target agent 'tester' is not connected")
}

// Test ValidateAgentStateConsistency - State Consistency Validation
func (suite *ProtocolValidatorTestSuite) TestValidateAgentStateConsistency_ConsistentState() {
	// Set up consistent state: developer has barrel and is working
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)
	suite.testAgents["developer"].Activate("Working on feature")

	err = suite.validator.ValidateAgentStateConsistency("developer")

	assert.NoError(suite.T(), err)
}

func (suite *ProtocolValidatorTestSuite) TestValidateAgentStateConsistency_InconsistentState() {
	// Set up inconsistent state: developer has barrel but is waiting
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)
	// Agent remains in waiting state (not activated)

	err = suite.validator.ValidateAgentStateConsistency("developer")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "agent state inconsistency")
	assert.Contains(suite.T(), err.Error(), "has barrel but is waiting")
}

func (suite *ProtocolValidatorTestSuite) TestValidateAgentStateConsistency_AgentNotFound() {
	err := suite.validator.ValidateAgentStateConsistency("nonexistent")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "agent 'nonexistent' not found")
}

// Test ValidateYieldWorkflow - Complete Workflow Validation
func (suite *ProtocolValidatorTestSuite) TestValidateYieldWorkflow_CompleteValidWorkflow() {
	// Set up valid state
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)
	suite.testAgents["developer"].Activate("Working on feature")

	message := domain.NewYieldMessage("developer", "tester", "Code ready for testing")

	err = suite.validator.ValidateYieldWorkflow(message)

	assert.NoError(suite.T(), err)
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldWorkflow_InvalidMessage() {
	// Invalid message structure (empty message)
	var message domain.YieldMessage

	err := suite.validator.ValidateYieldWorkflow(message)

	assert.Error(suite.T(), err)
	// Since empty message fails on from_role first, we expect that error
	assert.Contains(suite.T(), err.Error(), "from_role cannot be empty")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldWorkflow_UnauthorizedYield() {
	// Set up state where tester tries to yield but doesn't have barrel
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)

	message := domain.NewYieldMessage("tester", "reviewer", "Unauthorized yield")

	err = suite.validator.ValidateYieldWorkflow(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "only current barrel holder can yield")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldWorkflow_InvalidTarget() {
	// Set up valid holder but invalid target
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)
	suite.testAgents["developer"].Activate("Working on feature")

	message := domain.NewYieldMessage("developer", "nonexistent", "Invalid target")

	err = suite.validator.ValidateYieldWorkflow(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "target agent 'nonexistent' not found")
}

// Test Edge Cases
func (suite *ProtocolValidatorTestSuite) TestValidateYieldWorkflow_YieldToPeople() {
	// Yield to people should always be valid if holder is correct
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)
	suite.testAgents["developer"].Activate("Working on feature")

	message := domain.NewYieldMessage("developer", "people", "Need guidance")

	err = suite.validator.ValidateYieldWorkflow(message)

	assert.NoError(suite.T(), err)
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldWorkflow_PeopleYieldToAgent() {
	// People should already have the barrel initially, so no need to transfer
	// Just validate that people can yield to an agent
	message := domain.NewYieldMessage("people", "developer", "Start working on feature X")

	err := suite.validator.ValidateYieldWorkflow(message)

	assert.NoError(suite.T(), err)
}

// Test GetValidationErrors - Comprehensive Error Collection
func (suite *ProtocolValidatorTestSuite) TestGetValidationErrors_MultipleErrors() {
	// Create a scenario with multiple validation errors
	var message domain.YieldMessage // Invalid message

	errors := suite.validator.GetValidationErrors(message)

	assert.NotEmpty(suite.T(), errors)
	assert.True(suite.T(), len(errors) > 0)
}

func (suite *ProtocolValidatorTestSuite) TestGetValidationErrors_NoErrors() {
	// Set up valid state
	err := suite.testBarrel.TransferTo("developer", "Working on feature")
	suite.Require().NoError(err)
	suite.testAgents["developer"].Activate("Working on feature")

	message := domain.NewYieldMessage("developer", "tester", "Code ready for testing")

	errors := suite.validator.GetValidationErrors(message)

	assert.Empty(suite.T(), errors)
}

// Run the test suite
func TestProtocolValidatorSuite(t *testing.T) {
	suite.Run(t, new(ProtocolValidatorTestSuite))
}

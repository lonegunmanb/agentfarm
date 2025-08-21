package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ProtocolValidatorTestSuite provides comprehensive test coverage for protocol validation
type ProtocolValidatorTestSuite struct {
	suite.Suite
	validator  *ProtocolValidator
	soviet     *SovietState
	testAgents map[string]*AgentComrade
	testBarrel *BarrelOfGun
}

func (suite *ProtocolValidatorTestSuite) SetupTest() {
	// Create test soviet state
	repo := NewMemoryAgentRepository()
	suite.soviet = NewSovietState(repo)

	// Create validator
	suite.validator = NewProtocolValidator(suite.soviet)

	// Create test barrel
	suite.testBarrel = NewBarrelOfGun()
	suite.soviet.SetBarrel(suite.testBarrel)

	// Create test agents
	suite.testAgents = map[string]*AgentComrade{
		"developer": NewAgentComrade("developer", []string{"code", "test"}),
		"tester":    NewAgentComrade("tester", []string{"test", "validate"}),
		"reviewer":  NewAgentComrade("reviewer", []string{"review", "approve"}),
	}

	// Set agents as connected (they start as disconnected by default)
	for _, agent := range suite.testAgents {
		agent.SetConnected(true)
	}

	// Register test agents
	for _, agent := range suite.testAgents {
		err := suite.soviet.SimpleRegisterAgent(agent)
		suite.Require().NoError(err)
	}
}

func TestProtocolValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(ProtocolValidatorTestSuite))
}

// Test ValidateYieldMessage - Message Structure Validation
func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_ValidMessage() {
	message := NewYieldMessage("developer", "tester", "Valid test message")

	err := suite.validator.ValidateYieldMessage(message)

	assert.NoError(suite.T(), err, "Valid message should pass validation")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_EmptyFromRole() {
	message := NewYieldMessage("", "tester", "Test message")

	err := suite.validator.ValidateYieldMessage(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "from_role cannot be empty")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_EmptyToRole() {
	message := NewYieldMessage("developer", "", "Test message")

	err := suite.validator.ValidateYieldMessage(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "to_role cannot be empty")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldMessage_SelfYield() {
	message := NewYieldMessage("developer", "developer", "Test message")

	err := suite.validator.ValidateYieldMessage(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "agent cannot yield to itself")
}

// Test ValidateBarrelHolderRights - Barrel Ownership Validation
func (suite *ProtocolValidatorTestSuite) TestValidateBarrelHolderRights_ValidHolder() {
	// Give barrel to developer
	suite.testBarrel.TransferTo("developer", "Test payload")

	err := suite.validator.ValidateBarrelHolderRights("developer")

	assert.NoError(suite.T(), err, "Current barrel holder should be able to yield")
}

func (suite *ProtocolValidatorTestSuite) TestValidateBarrelHolderRights_NotCurrentHolder() {
	// Give barrel to developer
	suite.testBarrel.TransferTo("developer", "Test payload")

	// Tester tries to yield
	err := suite.validator.ValidateBarrelHolderRights("tester")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "only current barrel holder can yield")
}

func (suite *ProtocolValidatorTestSuite) TestValidateBarrelHolderRights_PeopleCanAlwaysYield() {
	// Give barrel to developer
	suite.testBarrel.TransferTo("developer", "Test payload")

	// People should always be able to yield
	err := suite.validator.ValidateBarrelHolderRights("people")

	assert.NoError(suite.T(), err, "People should always be able to yield")
}

// Test ValidateTargetAgent - Target Agent Validation
func (suite *ProtocolValidatorTestSuite) TestValidateTargetAgent_ValidAgent() {
	err := suite.validator.ValidateTargetAgent("developer")

	assert.NoError(suite.T(), err, "Registered and connected agent should be valid target")
}

func (suite *ProtocolValidatorTestSuite) TestValidateTargetAgent_AgentNotFound() {
	err := suite.validator.ValidateTargetAgent("nonexistent")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "target agent 'nonexistent' not found")
}

func (suite *ProtocolValidatorTestSuite) TestValidateTargetAgent_PeopleAlwaysValid() {
	err := suite.validator.ValidateTargetAgent("people")

	assert.NoError(suite.T(), err, "People should always be a valid target")
}

func (suite *ProtocolValidatorTestSuite) TestValidateTargetAgent_DisconnectedAgent() {
	// Disconnect the agent
	suite.testAgents["developer"].SetConnected(false)

	err := suite.validator.ValidateTargetAgent("developer")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "target agent 'developer' is not connected")
}

// Test ValidateYieldWorkflow - Complete Workflow Validation
func (suite *ProtocolValidatorTestSuite) TestValidateYieldWorkflow_CompleteValidWorkflow() {
	// Set up valid scenario: developer has barrel and wants to yield to tester
	suite.testBarrel.TransferTo("developer", "Initial work")
	suite.testAgents["developer"].TransitionTo(AgentStateWorking)

	message := NewYieldMessage("developer", "tester", "Work completed")

	err := suite.validator.ValidateYieldWorkflow(message)

	assert.NoError(suite.T(), err, "Complete valid workflow should pass all validation")
}

func (suite *ProtocolValidatorTestSuite) TestValidateYieldWorkflow_InvalidMessage() {
	message := NewYieldMessage("", "tester", "Invalid message")

	err := suite.validator.ValidateYieldWorkflow(message)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "from_role cannot be empty")
}

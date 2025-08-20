package domain

import (
	"fmt"
)

// ProtocolValidator enforces revolutionary discipline and validation rules
// It provides comprehensive validation for yield messages and agent states
type ProtocolValidator struct {
	soviet *SovietState
}

// NewProtocolValidator creates a new protocol validator with the given soviet state
func NewProtocolValidator(soviet *SovietState) *ProtocolValidator {
	if soviet == nil {
		panic("soviet state cannot be nil")
	}
	return &ProtocolValidator{
		soviet: soviet,
	}
}

// ValidateYieldMessage validates the structure and content of a yield message
func (v *ProtocolValidator) ValidateYieldMessage(message YieldMessage) error {
	fromRole := message.FromRole()
	toRole := message.ToRole()

	// Check for empty roles first for specific error messages
	if fromRole == "" {
		return fmt.Errorf("from_role cannot be empty")
	}

	if toRole == "" {
		return fmt.Errorf("to_role cannot be empty")
	}

	// Check if message is valid (uses the domain's IsValid method)
	if !message.IsValid() {
		return fmt.Errorf("invalid yield message: missing required fields")
	}

	// Check for self-yield
	if fromRole == toRole {
		return fmt.Errorf("agent cannot yield to itself: %s", fromRole)
	}

	return nil
}

// ValidateBarrelHolderRights validates that the requester has the right to yield the barrel
func (v *ProtocolValidator) ValidateBarrelHolderRights(requesterRole string) error {
	// People always have the right to yield
	if requesterRole == "people" {
		return nil
	}

	// Get the barrel
	barrel := v.soviet.GetBarrel()
	if barrel == nil {
		return fmt.Errorf("no barrel available in soviet")
	}

	// Check if the requester is the current barrel holder
	if !barrel.IsHeldBy(requesterRole) {
		return fmt.Errorf("only current barrel holder can yield (current holder: %s, requester: %s)",
			barrel.CurrentHolder(), requesterRole)
	}

	return nil
}

// ValidateTargetAgent validates that the target agent exists and can receive the barrel
func (v *ProtocolValidator) ValidateTargetAgent(targetRole string) error {
	// People is always a valid target
	if targetRole == "people" {
		return nil
	}

	// Check if agent exists
	if !v.soviet.IsAgentRegistered(targetRole) {
		return fmt.Errorf("target agent '%s' not found", targetRole)
	}

	// Check if agent is connected
	agent := v.soviet.GetAgent(targetRole)
	if agent != nil && !agent.IsConnected() {
		return fmt.Errorf("target agent '%s' is not connected", targetRole)
	}

	return nil
}

// ValidateAgentStateConsistency validates that agent state is consistent with barrel ownership
func (v *ProtocolValidator) ValidateAgentStateConsistency(agentRole string) error {
	// Get the agent
	agent := v.soviet.GetAgent(agentRole)
	if agent == nil {
		return fmt.Errorf("agent '%s' not found", agentRole)
	}

	// Get the barrel
	barrel := v.soviet.GetBarrel()
	if barrel == nil {
		return fmt.Errorf("no barrel available in soviet")
	}

	// Check consistency: if agent has barrel, they should be working
	hasBarrel := barrel.IsHeldBy(agentRole)
	isWorking := agent.State() == AgentStateWorking

	if hasBarrel && !isWorking {
		return fmt.Errorf("agent state inconsistency: agent '%s' has barrel but is waiting", agentRole)
	}

	if !hasBarrel && isWorking {
		return fmt.Errorf("agent state inconsistency: agent '%s' is working but doesn't have barrel", agentRole)
	}

	return nil
}

// ValidateYieldWorkflow performs comprehensive validation of the entire yield workflow
func (v *ProtocolValidator) ValidateYieldWorkflow(message YieldMessage) error {
	// 1. Validate message structure
	if err := v.ValidateYieldMessage(message); err != nil {
		return err
	}

	// 2. Validate barrel holder rights
	if err := v.ValidateBarrelHolderRights(message.FromRole()); err != nil {
		return err
	}

	// 3. Validate target agent
	if err := v.ValidateTargetAgent(message.ToRole()); err != nil {
		return err
	}

	// 4. Validate state consistency (only for non-people agents)
	if message.FromRole() != "people" {
		if err := v.ValidateAgentStateConsistency(message.FromRole()); err != nil {
			return err
		}
	}

	return nil
}

// GetValidationErrors collects all validation errors for a yield message
func (v *ProtocolValidator) GetValidationErrors(message YieldMessage) []error {
	var errors []error

	// Collect all validation errors without short-circuiting
	if err := v.ValidateYieldMessage(message); err != nil {
		errors = append(errors, err)
	}

	if err := v.ValidateBarrelHolderRights(message.FromRole()); err != nil {
		errors = append(errors, err)
	}

	if err := v.ValidateTargetAgent(message.ToRole()); err != nil {
		errors = append(errors, err)
	}

	if message.FromRole() != "people" {
		if err := v.ValidateAgentStateConsistency(message.FromRole()); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

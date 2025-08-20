package services

import (
	"fmt"

	"github.com/lonegunmanb/agentfarm/pkg/core/domain"
)

// SovietCoordinator implements the core business logic for managing agents and barrel transfers
// It serves as the central coordination service for the Agent Farm collective
type SovietCoordinator struct {
	soviet    *domain.SovietState
	validator *ProtocolValidator
}

// NewSovietCoordinator creates a new coordinator with the given soviet state
func NewSovietCoordinator(soviet *domain.SovietState) *SovietCoordinator {
	if soviet == nil {
		panic("soviet state cannot be nil")
	}
	return &SovietCoordinator{
		soviet:    soviet,
		validator: NewProtocolValidator(soviet),
	}
}

// RegisterAgent registers a new agent or handles reconnection intelligently
// This unified method handles both new registrations and reconnections automatically
// Returns: (shouldResume, lastMessage, error) where shouldResume indicates if agent should start working
func (c *SovietCoordinator) RegisterAgent(agent *domain.AgentComrade) (bool, string, error) {
	if agent == nil {
		return false, "", fmt.Errorf("agent cannot be nil")
	}

	role := agent.Role()

	// Check if an agent with this role already exists
	if existingAgent := c.soviet.GetAgent(role); existingAgent != nil {
		// Disconnect the existing agent (replacement behavior)
		existingAgent.SetConnected(false)

		// Unregister the existing agent to make room for the new one
		err := c.soviet.UnregisterAgent(role)
		if err != nil {
			return false, "", fmt.Errorf("failed to unregister existing agent: %w", err)
		}
	}

	// Register the new agent
	err := c.soviet.RegisterAgent(agent)
	if err != nil {
		return false, "", fmt.Errorf("failed to register agent: %w", err)
	}

	// Set the agent as connected and in waiting state initially
	agent.SetConnected(true)
	agent.TransitionTo(domain.AgentStateWaiting)

	// Check if this agent role should resume work (if they hold the barrel)
	barrel := c.soviet.GetBarrel()
	if barrel != nil && barrel.IsHeldBy(role) {
		// Agent should resume work - activate them
		lastMessage := barrel.LastMessage()
		agent.TransitionTo(domain.AgentStateWorking)
		return true, lastMessage, nil
	}

	// Agent doesn't hold barrel, remains in waiting state
	return false, "", nil
}

// DeregisterAgent removes an agent from the collective
// If the agent holds the barrel, it's transferred back to the people
func (c *SovietCoordinator) DeregisterAgent(role string) error {
	if !c.soviet.IsAgentRegistered(role) {
		return fmt.Errorf("agent with role '%s' not found", role)
	}

	// Check if this agent holds the barrel
	if c.soviet.IsBarrelHeldBy(role) {
		// Transfer barrel back to the people
		barrel := c.soviet.GetBarrel()
		if barrel != nil {
			barrel.TransferTo("people", fmt.Sprintf("Agent '%s' deregistered, returning barrel to people", role))
		}
	}

	// Remove the agent from the soviet
	err := c.soviet.UnregisterAgent(role)
	if err != nil {
		return fmt.Errorf("failed to deregister agent: %w", err)
	}

	return nil
}

// ProcessYield handles yield requests and manages barrel transfers
func (c *SovietCoordinator) ProcessYield(message domain.YieldMessage) error {
	// Use the protocol validator for comprehensive validation
	if err := c.validator.ValidateYieldWorkflow(message); err != nil {
		return err
	}

	fromRole := message.FromRole()
	toRole := message.ToRole()
	payload := message.Payload()

	// Get the source agent and transition it to waiting
	sourceAgent := c.soviet.GetAgent(fromRole)
	if sourceAgent != nil {
		sourceAgent.Yield() // This transitions the agent to waiting state
	}

	// Transfer the barrel
	barrel := c.soviet.GetBarrel()
	err := barrel.TransferTo(toRole, payload)
	if err != nil {
		return fmt.Errorf("failed to transfer barrel: %w", err)
	}

	// If transferring to an agent, activate them
	if toRole != "people" {
		targetAgent := c.soviet.GetAgent(toRole)
		if targetAgent != nil {
			targetAgent.Activate(payload) // This transitions the agent to working state
		}
	}

	return nil
}

// GetAgentState returns the current state of an agent
func (c *SovietCoordinator) GetAgentState(role string) (domain.AgentState, error) {
	agent := c.soviet.GetAgent(role)
	if agent == nil {
		return domain.AgentStateWaiting, fmt.Errorf("agent with role '%s' not found", role)
	}

	return agent.State(), nil
}

// GetBarrelStatus returns the role that currently holds the barrel
func (c *SovietCoordinator) GetBarrelStatus() string {
	barrel := c.soviet.GetBarrel()
	if barrel == nil {
		return "people" // Default to people if no barrel
	}
	return barrel.CurrentHolder()
}

// GetRegisteredAgents returns a list of all registered agent roles
func (c *SovietCoordinator) GetRegisteredAgents() []string {
	return c.soviet.GetAgentRoles()
}

package ports

import "github.com/lonegunmanb/agentfarm/pkg/core/domain"

// AgentService defines the primary port for querying agent and barrel information
// This interface represents the query use cases for the Agent Farm application
// External adapters (TCP, CLI, etc.) will call these methods to query system state
type AgentService interface {
	// GetAgentState returns the current state of a specific agent
	// Returns the agent's current state (waiting/working) or an error if agent not found
	GetAgentState(role string) (domain.AgentState, error)
	
	// GetBarrelStatus returns the role that currently holds the barrel of gun
	// Returns "people" if no agent currently holds the barrel
	GetBarrelStatus() string
	
	// GetRegisteredAgents returns a list of all currently registered agent roles
	// This provides an overview of all agents known to the collective
	GetRegisteredAgents() []string
}

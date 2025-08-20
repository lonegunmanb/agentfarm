package ports

import "github.com/lonegunmanb/agentfarm/pkg/core/domain"

// SovietService defines the primary port for commanding the Soviet coordinator
// This interface represents the use cases that drive the Agent Farm application
// External adapters (TCP, CLI, etc.) will call these methods to interact with the core domain
type SovietService interface {
	// RegisterAgent registers a new agent comrade or replaces an existing one with the same role
	// This is called when an agent connects and declares service to the collective
	RegisterAgent(agent *domain.AgentComrade) error
	
	// ProcessYield handles yield requests and manages barrel transfers
	// This is called when an agent comrade yields the barrel to another agent or to the people
	ProcessYield(message domain.YieldMessage) error
	
	// DeregisterAgent removes an agent from the collective
	// This is called when an agent disconnects or is manually removed
	DeregisterAgent(role string) error
	
	// HandleReconnection handles the reconnection of an agent and determines if they should resume work
	// Returns: (shouldResume, lastMessage, error)
	// This is called when a previously registered agent reconnects to the collective
	HandleReconnection(role string) (bool, string, error)
	
	// QueryStatus returns the current status of the collective including all agents and barrel state
	// This is called by People's representatives to inspect the collective
	QueryStatus() StatusResponse
}

// StatusResponse represents the current status of the Agent Farm collective
// This is returned by QueryStatus to provide a complete view of the system state
type StatusResponse struct {
	// BarrelHolder indicates which role currently holds the barrel of gun
	BarrelHolder string `json:"barrel_holder"`
	
	// RegisteredAgents contains all currently registered agent roles
	RegisteredAgents []string `json:"registered_agents"`
	
	// AgentStates maps agent roles to their current states
	AgentStates map[string]domain.AgentState `json:"agent_states"`
	
	// ConnectedAgents indicates which agents are currently connected
	ConnectedAgents map[string]bool `json:"connected_agents"`
}

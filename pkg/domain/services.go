package domain

// AgentDetails represents detailed information about an agent comrade
type AgentDetails struct {
	Role         string      `json:"role"`
	Capabilities []string    `json:"capabilities"`
	State        AgentState  `json:"state"`
	Connected    bool        `json:"connected"`
}

// SovietService defines the primary port for commanding the Soviet coordinator
// This interface represents the use cases that drive the Agent Farm application
// External adapters (TCP, CLI, etc.) will call these methods to interact with the core domain
type SovietService interface {
	// RegisterAgent registers a new agent comrade or handles reconnection intelligently
	// When an agent connects (new or reconnection), they simply register their role.
	// The system automatically handles both cases:
	// - For new agents: registers and places in waiting state
	// - For reconnections: replaces existing agent and resumes work if role holds barrel
	// Returns: (shouldResume, lastMessage, error) where shouldResume indicates if agent should start working
	RegisterAgent(agent *AgentComrade) (bool, string, error)

	// ProcessYield handles yield requests and manages barrel transfers
	// This is called when an agent comrade yields the barrel to another agent or to the people
	ProcessYield(message YieldMessage) error

	// DeregisterAgent removes an agent from the collective
	// This is called when an agent disconnects or is manually removed
	DeregisterAgent(role string) error

	// QueryStatus returns the current status of the collective including all agents and barrel state
	// This is called by People's representatives to inspect the collective
	QueryStatus() StatusResponse
}

// AgentService defines the primary port for querying agent and barrel information
// This interface represents the query use cases for the Agent Farm application
// External adapters (TCP, CLI, etc.) will call these methods to query system state
type AgentService interface {
	// GetAgentState returns the current state of a specific agent
	// Returns the agent's current state (waiting/working) or an error if agent not found
	GetAgentState(role string) (AgentState, error)

	// GetBarrelStatus returns the role that currently holds the barrel of gun
	// Returns "people" if no agent currently holds the barrel
	GetBarrelStatus() string

	// GetRegisteredAgents returns a list of all currently registered agent roles
	// This provides an overview of all agents known to the collective
	GetRegisteredAgents() []string

	// GetAgentDetails returns detailed information about all registered agents including capabilities
	// This provides a comprehensive view of all agents and their capabilities for the collective
	GetAgentDetails() []AgentDetails
}

// StatusResponse represents the current status of the Agent Farm collective
// This is returned by QueryStatus to provide a complete view of the system state
type StatusResponse struct {
	// BarrelHolder indicates which role currently holds the barrel of gun
	BarrelHolder string `json:"barrel_holder"`

	// RegisteredAgents contains all currently registered agent roles
	RegisteredAgents []string `json:"registered_agents"`

	// AgentStates maps agent roles to their current states
	AgentStates map[string]AgentState `json:"agent_states"`

	// ConnectedAgents indicates which agents are currently connected
	ConnectedAgents map[string]bool `json:"connected_agents"`
}

// CommandHandler defines the port for handling incoming commands from external sources
// This interface represents how external adapters (TCP, CLI, etc.) send commands to the system
type CommandHandler interface {
	// HandleCommand processes an incoming command message
	// Returns a response message or error
	HandleCommand(command map[string]interface{}) (map[string]interface{}, error)
}

package ports

import "github.com/lonegunmanb/agentfarm/pkg/core/domain"

// AgentRepository defines the port for agent persistence
// This interface abstracts agent storage operations from the core domain
type AgentRepository interface {
	// Store persists an agent to the repository
	Store(agent *domain.AgentComrade) error

	// GetByRole retrieves an agent by their role
	GetByRole(role string) (*domain.AgentComrade, error)

	// GetAll retrieves all agents
	GetAll() ([]*domain.AgentComrade, error)

	// Delete removes an agent from the repository
	Delete(role string) error

	// Exists checks if an agent with the given role exists
	Exists(role string) bool
}

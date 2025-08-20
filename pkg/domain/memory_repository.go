package domain

import (
	"fmt"
	"sync"
)

// MemoryAgentRepository implements AgentRepository interface with in-memory storage
// This is a simple implementation for testing and development purposes
type MemoryAgentRepository struct {
	agents map[string]*AgentComrade
	mutex  sync.RWMutex
}

// NewMemoryAgentRepository creates a new in-memory agent repository
func NewMemoryAgentRepository() *MemoryAgentRepository {
	return &MemoryAgentRepository{
		agents: make(map[string]*AgentComrade),
	}
}

// Store persists an agent to the repository
func (m *MemoryAgentRepository) Store(agent *AgentComrade) error {
	if agent == nil {
		return fmt.Errorf("agent cannot be nil")
	}

	role := agent.Role()
	if role == "" {
		return fmt.Errorf("agent role cannot be empty")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.agents[role] = agent
	return nil
}

// GetByRole retrieves an agent by their role
func (m *MemoryAgentRepository) GetByRole(role string) (*AgentComrade, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	agent, exists := m.agents[role]
	if !exists {
		return nil, fmt.Errorf("agent with role '%s' not found", role)
	}
	return agent, nil
}

// GetAll retrieves all agents
func (m *MemoryAgentRepository) GetAll() ([]*AgentComrade, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	agents := make([]*AgentComrade, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}
	return agents, nil
}

// Delete removes an agent from the repository
func (m *MemoryAgentRepository) Delete(role string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.agents[role]; !exists {
		return fmt.Errorf("agent with role '%s' not found", role)
	}

	delete(m.agents, role)
	return nil
}

// Exists checks if an agent with the given role exists
func (m *MemoryAgentRepository) Exists(role string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.agents[role]
	return exists
}

// Ensure MemoryAgentRepository implements AgentRepository
var _ AgentRepository = (*MemoryAgentRepository)(nil)

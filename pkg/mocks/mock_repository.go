package mocks

import (
	"fmt"
	"sync"

	"github.com/lonegunmanb/agentfarm/pkg/core/domain"
	"github.com/lonegunmanb/agentfarm/pkg/ports"
)

// MockAgentRepository implements AgentRepository interface for testing
type MockAgentRepository struct {
	mu     sync.RWMutex
	agents map[string]*domain.AgentComrade
}

// NewMockAgentRepository creates a new mock agent repository
func NewMockAgentRepository() *MockAgentRepository {
	return &MockAgentRepository{
		agents: make(map[string]*domain.AgentComrade),
	}
}

// Store persists an agent to the repository
func (m *MockAgentRepository) Store(agent *domain.AgentComrade) error {
	if agent == nil {
		return fmt.Errorf("agent cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.agents[agent.Role()] = agent
	return nil
}

// GetByRole retrieves an agent by their role
func (m *MockAgentRepository) GetByRole(role string) (*domain.AgentComrade, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	agent, exists := m.agents[role]
	if !exists {
		return nil, fmt.Errorf("agent with role '%s' not found", role)
	}

	return agent, nil
}

// GetAll retrieves all agents
func (m *MockAgentRepository) GetAll() ([]*domain.AgentComrade, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	agents := make([]*domain.AgentComrade, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}

	return agents, nil
}

// Delete removes an agent from the repository
func (m *MockAgentRepository) Delete(role string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[role]; !exists {
		return fmt.Errorf("agent with role '%s' not found", role)
	}

	delete(m.agents, role)
	return nil
}

// Exists checks if an agent with the given role exists
func (m *MockAgentRepository) Exists(role string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.agents[role]
	return exists
}

// Verify interface compliance
var _ ports.AgentRepository = (*MockAgentRepository)(nil)

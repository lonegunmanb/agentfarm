package domain

import (
	"fmt"
	"time"
)

// SovietStats represents statistics about the soviet state
type SovietStats struct {
	TotalAgents         int       `json:"total_agents"`
	ConnectedAgents     int       `json:"connected_agents"`
	CurrentBarrelHolder string    `json:"current_barrel_holder"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	DeactivatedAt       time.Time `json:"deactivated_at,omitempty"`
}

// SovietState represents the state of the collective, managing all agents and the barrel
// It remains a pure domain object without external dependencies
type SovietState struct {
	agents        map[string]*AgentComrade
	barrel        *BarrelOfGun
	active        bool
	createdAt     time.Time
	deactivatedAt time.Time
}

// NewSovietState creates a new soviet state managing the collective
func NewSovietState() *SovietState {
	return &SovietState{
		agents:    make(map[string]*AgentComrade),
		active:    true,
		createdAt: nowFunc(),
	}
}

// CreatedAt returns when the soviet was created
func (s *SovietState) CreatedAt() time.Time {
	return s.createdAt
}

// IsActive returns whether the soviet is currently active
func (s *SovietState) IsActive() bool {
	return s.active
}

// Activate sets the soviet to active state
func (s *SovietState) Activate() {
	s.active = true
	s.deactivatedAt = time.Time{}
}

// Deactivate sets the soviet to inactive state
func (s *SovietState) Deactivate() {
	s.active = false
	s.deactivatedAt = nowFunc()
}

// DeactivatedAt returns when the soviet was deactivated (zero time if active)
func (s *SovietState) DeactivatedAt() time.Time {
	return s.deactivatedAt
}

// SetBarrel sets the barrel of gun for the soviet to manage
func (s *SovietState) SetBarrel(barrel *BarrelOfGun) error {
	if barrel == nil {
		return fmt.Errorf("barrel cannot be nil")
	}
	s.barrel = barrel
	return nil
}

// GetBarrel returns the current barrel of gun
func (s *SovietState) GetBarrel() *BarrelOfGun {
	return s.barrel
}

// RegisterAgent registers a new agent with the soviet
func (s *SovietState) RegisterAgent(agent *AgentComrade) error {
	if agent == nil {
		return fmt.Errorf("agent cannot be nil")
	}

	role := agent.Role()
	if role == "" {
		return fmt.Errorf("agent role cannot be empty")
	}

	if _, exists := s.agents[role]; exists {
		return fmt.Errorf("agent with role '%s' is already registered", role)
	}

	s.agents[role] = agent
	return nil
}

// UnregisterAgent removes an agent from the soviet
func (s *SovietState) UnregisterAgent(role string) error {
	if role == "" {
		return fmt.Errorf("role cannot be empty")
	}

	_, exists := s.agents[role]
	if !exists {
		return fmt.Errorf("agent with role '%s' is not registered", role)
	}

	delete(s.agents, role)
	return nil
}

// IsAgentRegistered checks if an agent with the given role is registered
func (s *SovietState) IsAgentRegistered(role string) bool {
	_, exists := s.agents[role]
	return exists
}

// GetAgent returns the agent with the specified role
func (s *SovietState) GetAgent(role string) *AgentComrade {
	return s.agents[role]
}

// RegisteredAgents returns a copy of all registered agents
func (s *SovietState) RegisteredAgents() map[string]*AgentComrade {
	result := make(map[string]*AgentComrade)
	for role, agent := range s.agents {
		result[role] = agent
	}
	return result
}

// GetAgentRoles returns a slice of all registered agent roles
func (s *SovietState) GetAgentRoles() []string {
	roles := make([]string, 0, len(s.agents))
	for role := range s.agents {
		roles = append(roles, role)
	}
	return roles
}

// CurrentBarrelHolder returns the role that currently holds the barrel
func (s *SovietState) CurrentBarrelHolder() string {
	if s.barrel == nil {
		return ""
	}
	return s.barrel.CurrentHolder()
}

// IsBarrelHeldBy checks if the barrel is currently held by the specified role
func (s *SovietState) IsBarrelHeldBy(role string) bool {
	if s.barrel == nil {
		return false
	}
	return s.barrel.IsHeldBy(role)
}

// ProcessBarrelTransfer handles barrel transfer
func (s *SovietState) ProcessBarrelTransfer(fromRole, toRole, payload string) error {
	if s.barrel == nil {
		return fmt.Errorf("no barrel available for transfer")
	}

	return s.barrel.TransferTo(toRole, payload)
}

// GetStats returns statistics about the current soviet state
func (s *SovietState) GetStats() *SovietStats {
	totalAgents := len(s.agents)
	connectedAgents := 0

	for _, agent := range s.agents {
		if agent.IsConnected() {
			connectedAgents++
		}
	}

	currentHolder := ""
	if s.barrel != nil {
		currentHolder = s.barrel.CurrentHolder()
	}

	return &SovietStats{
		TotalAgents:         totalAgents,
		ConnectedAgents:     connectedAgents,
		CurrentBarrelHolder: currentHolder,
		IsActive:            s.active,
		CreatedAt:           s.createdAt,
		DeactivatedAt:       s.deactivatedAt,
	}
}

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
	validator     *ProtocolValidator

	// External dependencies for side effects (optional, can be nil)
	repo   AgentRepository
	sender MessageSender
	logger Logger
}

// NewSovietState creates a new soviet state managing the collective
func NewSovietState() *SovietState {
	soviet := &SovietState{
		agents:    make(map[string]*AgentComrade),
		active:    true,
		createdAt: nowFunc(),
	}
	soviet.validator = NewProtocolValidator(soviet)
	return soviet
}

// NewSovietStateWithDependencies creates a new soviet state with external dependencies
func NewSovietStateWithDependencies(
	repo AgentRepository,
	sender MessageSender,
	logger Logger,
) *SovietState {
	soviet := &SovietState{
		agents:    make(map[string]*AgentComrade),
		active:    true,
		createdAt: nowFunc(),
		repo:      repo,
		sender:    sender,
		logger:    logger,
	}
	soviet.validator = NewProtocolValidator(soviet)
	return soviet
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

// Coordinator methods (moved from services/coordinator.go)

// RegisterAgent registers a new agent or handles reconnection intelligently
// This unified method handles both new registrations and reconnections automatically
// Returns: (shouldResume, lastMessage, error) where shouldResume indicates if agent should start working
func (s *SovietState) RegisterAgent(agent *AgentComrade) (bool, string, error) {
	if agent == nil {
		return false, "", fmt.Errorf("agent cannot be nil")
	}

	role := agent.Role()

	// Check if an agent with this role already exists
	if existingAgent := s.GetAgent(role); existingAgent != nil {
		// Disconnect the existing agent (replacement behavior)
		existingAgent.SetConnected(false)

		// Unregister the existing agent to make room for the new one
		err := s.UnregisterAgent(role)
		if err != nil {
			return false, "", fmt.Errorf("failed to unregister existing agent: %w", err)
		}
	}

	// Register the new agent
	err := s.registerAgent(agent)
	if err != nil {
		return false, "", fmt.Errorf("failed to register agent: %w", err)
	}

	// Handle external operations if dependencies are available
	if s.repo != nil {
		if err := s.repo.Store(agent); err != nil {
			// Try to rollback domain operation
			s.UnregisterAgent(role)
			if s.logger != nil {
				s.logger.Error("Failed to persist agent", map[string]interface{}{
					"role":  role,
					"error": err.Error(),
				})
			}
			return false, "", fmt.Errorf("failed to persist agent: %w", err)
		}
	}

	if s.logger != nil {
		s.logger.Info("Agent registered successfully", map[string]interface{}{
			"role":       role,
			"agent_type": agent.Type(),
		})
	}

	// Set the agent as connected and in waiting state initially
	agent.SetConnected(true)
	agent.TransitionTo(AgentStateWaiting)

	// Check if this agent role should resume work (if they hold the barrel)
	barrel := s.GetBarrel()
	if barrel != nil && barrel.IsHeldBy(role) {
		// Agent should resume work - activate them
		lastMessage := barrel.LastMessage()
		agent.TransitionTo(AgentStateWorking)
		return true, lastMessage, nil
	}

	// Agent doesn't hold barrel, remains in waiting state
	return false, "", nil
}

// registerAgent is the internal registration method (renamed to avoid conflict)
func (s *SovietState) registerAgent(agent *AgentComrade) error {
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

// SimpleRegisterAgent provides the original simple registration for tests
// This is the domain-only registration without coordinator logic
func (s *SovietState) SimpleRegisterAgent(agent *AgentComrade) error {
	return s.registerAgent(agent)
}

// DeregisterAgent removes an agent from the collective
// If the agent holds the barrel, it's transferred back to the people
func (s *SovietState) DeregisterAgent(role string) error {
	if !s.IsAgentRegistered(role) {
		return fmt.Errorf("agent with role '%s' not found", role)
	}

	// Check if this agent holds the barrel
	if s.IsBarrelHeldBy(role) {
		// Transfer barrel back to the people
		barrel := s.GetBarrel()
		if barrel != nil {
			barrel.TransferTo("people", fmt.Sprintf("Agent '%s' deregistered, returning barrel to people", role))
		}
	}

	// Remove the agent from the soviet
	err := s.UnregisterAgent(role)
	if err != nil {
		return fmt.Errorf("failed to deregister agent: %w", err)
	}

	// Handle external operations if dependencies are available
	if s.repo != nil {
		if err := s.repo.Delete(role); err != nil {
			if s.logger != nil {
				s.logger.Error("Failed to delete agent from persistence", map[string]interface{}{
					"role":  role,
					"error": err.Error(),
				})
			}
			// Note: We don't rollback domain operation for persistence failures
		}
	}

	if s.logger != nil {
		s.logger.Info("Agent deregistered successfully", map[string]interface{}{
			"role": role,
		})
	}

	return nil
}

// ProcessYield handles yield requests and manages barrel transfers
func (s *SovietState) ProcessYield(message YieldMessage) error {
	// Use the protocol validator for comprehensive validation
	if err := s.validator.ValidateYieldWorkflow(message); err != nil {
		return err
	}

	fromRole := message.FromRole()
	toRole := message.ToRole()
	payload := message.Payload()

	// Get the source agent and transition it to waiting
	sourceAgent := s.GetAgent(fromRole)
	if sourceAgent != nil {
		sourceAgent.Yield() // This transitions the agent to waiting state
	}

	// Use SovietState to handle barrel transfer
	err := s.ProcessBarrelTransfer(fromRole, toRole, payload)
	if err != nil {
		return err
	}

	// Handle external operations if dependencies are available

	// Send activation to target agent (if not people)
	if toRole != "people" && s.sender != nil {
		if err := s.sender.SendActivation(toRole, payload); err != nil {
			if s.logger != nil {
				s.logger.Error("Failed to send activation message", map[string]interface{}{
					"role":  toRole,
					"error": err.Error(),
				})
			}
		}
	}

	// Log successful transfer
	if s.logger != nil {
		s.logger.Info("Barrel transferred successfully", map[string]interface{}{
			"from_role": fromRole,
			"to_role":   toRole,
			"payload":   payload,
		})
	}

	// If transferring to an agent, activate them
	if toRole != "people" {
		targetAgent := s.GetAgent(toRole)
		if targetAgent != nil {
			targetAgent.Activate(payload) // This transitions the agent to working state
		}
	}

	return nil
}

// GetAgentState returns the current state of an agent
func (s *SovietState) GetAgentState(role string) (AgentState, error) {
	agent := s.GetAgent(role)
	if agent == nil {
		return AgentStateWaiting, fmt.Errorf("agent with role '%s' not found", role)
	}

	return agent.State(), nil
}

// GetBarrelStatus returns the role that currently holds the barrel
func (s *SovietState) GetBarrelStatus() string {
	barrel := s.GetBarrel()
	if barrel == nil {
		return "people" // Default to people if no barrel
	}
	return barrel.CurrentHolder()
}

// QueryStatus returns the current status of the collective including all agents and barrel state
func (s *SovietState) QueryStatus() StatusResponse {
	agentStates := make(map[string]AgentState)
	connectedAgents := make(map[string]bool)

	for role, agent := range s.agents {
		agentStates[role] = agent.State()
		connectedAgents[role] = agent.IsConnected()
	}

	return StatusResponse{
		BarrelHolder:     s.GetBarrelStatus(),
		RegisteredAgents: s.GetAgentRoles(),
		AgentStates:      agentStates,
		ConnectedAgents:  connectedAgents,
	}
}

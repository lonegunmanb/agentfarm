package mocks

import (
	"github.com/lonegunmanb/agentfarm/pkg/core/domain"
	"github.com/lonegunmanb/agentfarm/pkg/core/services"
	"github.com/lonegunmanb/agentfarm/pkg/ports"
)

// CoordinatorAdapter adapts SovietCoordinator to implement SovietService interface
// This adapter demonstrates how the coordinator would be used through the ports interface
type CoordinatorAdapter struct {
	coordinator *services.SovietCoordinator
	soviet      *domain.SovietState
}

// NewCoordinatorAdapter creates a new adapter wrapping the coordinator
func NewCoordinatorAdapter(coordinator *services.SovietCoordinator, soviet *domain.SovietState) *CoordinatorAdapter {
	return &CoordinatorAdapter{
		coordinator: coordinator,
		soviet:      soviet,
	}
}

// RegisterAgent implements SovietService.RegisterAgent
func (a *CoordinatorAdapter) RegisterAgent(agent *domain.AgentComrade) (bool, string, error) {
	return a.coordinator.RegisterAgent(agent)
}

// ProcessYield implements SovietService.ProcessYield
func (a *CoordinatorAdapter) ProcessYield(message domain.YieldMessage) error {
	return a.coordinator.ProcessYield(message)
}

// DeregisterAgent implements SovietService.DeregisterAgent
func (a *CoordinatorAdapter) DeregisterAgent(role string) error {
	return a.coordinator.DeregisterAgent(role)
}

// QueryStatus implements SovietService.QueryStatus
func (a *CoordinatorAdapter) QueryStatus() ports.StatusResponse {
	barrelHolder := a.coordinator.GetBarrelStatus()
	registeredAgents := a.coordinator.GetRegisteredAgents()
	
	// Build agent states map
	agentStates := make(map[string]domain.AgentState)
	connectedAgents := make(map[string]bool)
	
	for _, role := range registeredAgents {
		state, err := a.coordinator.GetAgentState(role)
		if err == nil {
			agentStates[role] = state
		}
		
		// Check if agent is connected
		agent := a.soviet.GetAgent(role)
		if agent != nil {
			connectedAgents[role] = agent.IsConnected()
		}
	}
	
	return ports.StatusResponse{
		BarrelHolder:     barrelHolder,
		RegisteredAgents: registeredAgents,
		AgentStates:      agentStates,
		ConnectedAgents:  connectedAgents,
	}
}

// Verify interface compliance
var _ ports.SovietService = (*CoordinatorAdapter)(nil)

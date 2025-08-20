package mocks

import (
	"github.com/lonegunmanb/agentfarm/pkg/domain"
)

// CoordinatorAdapter adapts SovietState to implement SovietService interface
// This adapter demonstrates how the soviet state would be used through the ports interface
type CoordinatorAdapter struct {
	soviet *domain.SovietState
}

// NewCoordinatorAdapter creates a new adapter wrapping the soviet state
func NewCoordinatorAdapter(soviet *domain.SovietState) *CoordinatorAdapter {
	return &CoordinatorAdapter{
		soviet: soviet,
	}
}

// RegisterAgent implements SovietService.RegisterAgent
func (a *CoordinatorAdapter) RegisterAgent(agent *domain.AgentComrade) (bool, string, error) {
	return a.soviet.RegisterAgent(agent)
}

// ProcessYield implements SovietService.ProcessYield
func (a *CoordinatorAdapter) ProcessYield(message domain.YieldMessage) error {
	return a.soviet.ProcessYield(message)
}

// DeregisterAgent implements SovietService.DeregisterAgent
func (a *CoordinatorAdapter) DeregisterAgent(role string) error {
	return a.soviet.DeregisterAgent(role)
}

// QueryStatus implements SovietService.QueryStatus
func (a *CoordinatorAdapter) QueryStatus() domain.StatusResponse {
	return a.soviet.QueryStatus()
}

// Verify interface compliance
var _ domain.SovietService = (*CoordinatorAdapter)(nil)

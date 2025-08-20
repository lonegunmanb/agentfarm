package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create soviet with repository for tests
func newTestSoviet() *SovietState {
	return NewSovietState(NewMemoryAgentRepository())
}

func TestSovietState_NewSovietState(t *testing.T) {
	// RED: Test creation of new soviet state
	soviet := newTestSoviet()

	assert.NotNil(t, soviet)
	assert.NotZero(t, soviet.CreatedAt())
	assert.True(t, soviet.IsActive())
	assert.Equal(t, 0, len(soviet.RegisteredAgents()))
	assert.Nil(t, soviet.GetBarrel())
}

func TestSovietState_SetBarrel(t *testing.T) {
	// RED: Test barrel management
	soviet := newTestSoviet()
	barrel := NewBarrelOfGun()

	// Initially no barrel
	assert.Nil(t, soviet.GetBarrel())

	// Set barrel
	err := soviet.SetBarrel(barrel)
	assert.NoError(t, err)
	assert.NotNil(t, soviet.GetBarrel())
	assert.Equal(t, barrel, soviet.GetBarrel())

	// Cannot set nil barrel
	err = soviet.SetBarrel(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "barrel cannot be nil")
}

func TestSovietState_RegisterAgent(t *testing.T) {
	// RED: Test agent registration
	soviet := newTestSoviet()
	agent := NewAgentComrade("developer", "primary", []string{"code"})

	// Initially no agents
	assert.Equal(t, 0, len(soviet.RegisteredAgents()))
	assert.False(t, soviet.IsAgentRegistered("developer"))

	// Register agent
	err := soviet.SimpleRegisterAgent(agent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(soviet.RegisteredAgents()))
	assert.True(t, soviet.IsAgentRegistered("developer"))

	// Get registered agent
	retrievedAgent := soviet.GetAgent("developer")
	assert.NotNil(t, retrievedAgent)
	assert.Equal(t, agent, retrievedAgent)

	// Cannot register nil agent
	err = soviet.SimpleRegisterAgent(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent cannot be nil")

	// Cannot register agent with duplicate role
	duplicateAgent := NewAgentComrade("developer", "secondary", []string{"test"})
	err = soviet.SimpleRegisterAgent(duplicateAgent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent with role 'developer' is already registered")
}

func TestSovietState_UnregisterAgent(t *testing.T) {
	// RED: Test agent unregistration
	soviet := newTestSoviet()
	agent := NewAgentComrade("developer", "primary", []string{"code"})

	// Register first
	err := soviet.SimpleRegisterAgent(agent)
	assert.NoError(t, err)
	assert.True(t, soviet.IsAgentRegistered("developer"))

	// Unregister
	err = soviet.UnregisterAgent("developer")
	assert.NoError(t, err)
	assert.False(t, soviet.IsAgentRegistered("developer"))
	assert.Equal(t, 0, len(soviet.RegisteredAgents()))

	// Cannot unregister non-existent agent
	err = soviet.UnregisterAgent("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent with role 'nonexistent' is not registered")

	// Cannot unregister empty role
	err = soviet.UnregisterAgent("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role cannot be empty")
}

func TestSovietState_GetAgentRoles(t *testing.T) {
	// RED: Test getting all agent roles
	soviet := newTestSoviet()

	// Initially empty
	roles := soviet.GetAgentRoles()
	assert.Equal(t, 0, len(roles))

	// Add agents
	agent1 := NewAgentComrade("developer", "primary", []string{"code"})
	agent2 := NewAgentComrade("tester", "secondary", []string{"test"})

	soviet.SimpleRegisterAgent(agent1)
	soviet.SimpleRegisterAgent(agent2)

	roles = soviet.GetAgentRoles()
	assert.Equal(t, 2, len(roles))
	assert.Contains(t, roles, "developer")
	assert.Contains(t, roles, "tester")
}

func TestSovietState_CurrentBarrelHolder(t *testing.T) {
	// RED: Test current barrel holder tracking
	soviet := newTestSoviet()
	barrel := NewBarrelOfGun()

	// Set barrel
	soviet.SetBarrel(barrel)

	// Initially held by people
	holder := soviet.CurrentBarrelHolder()
	assert.Equal(t, "people", holder)

	// Transfer barrel to agent
	agent := NewAgentComrade("developer", "primary", []string{"code"})
	soviet.SimpleRegisterAgent(agent)

	err := barrel.TransferTo("developer", "Start working")
	assert.NoError(t, err)

	holder = soviet.CurrentBarrelHolder()
	assert.Equal(t, "developer", holder)
}

func TestSovietState_IsBarrelHeldBy(t *testing.T) {
	// RED: Test barrel ownership checking
	soviet := newTestSoviet()
	barrel := NewBarrelOfGun()
	soviet.SetBarrel(barrel)

	// Initially held by people
	assert.True(t, soviet.IsBarrelHeldBy("people"))
	assert.False(t, soviet.IsBarrelHeldBy("developer"))

	// Transfer to agent
	agent := NewAgentComrade("developer", "primary", []string{"code"})
	soviet.SimpleRegisterAgent(agent)

	barrel.TransferTo("developer", "Start working")
	assert.False(t, soviet.IsBarrelHeldBy("people"))
	assert.True(t, soviet.IsBarrelHeldBy("developer"))
}

func TestSovietState_Activate(t *testing.T) {
	// RED: Test soviet activation/deactivation
	soviet := newTestSoviet()

	// Initially active
	assert.True(t, soviet.IsActive())

	// Deactivate
	soviet.Deactivate()
	assert.False(t, soviet.IsActive())
	assert.NotZero(t, soviet.DeactivatedAt())

	// Reactivate
	soviet.Activate()
	assert.True(t, soviet.IsActive())
	assert.Zero(t, soviet.DeactivatedAt())
}

func TestSovietState_GetStats(t *testing.T) {
	// RED: Test getting soviet statistics
	soviet := newTestSoviet()
	barrel := NewBarrelOfGun()
	soviet.SetBarrel(barrel)

	stats := soviet.GetStats()
	assert.NotNil(t, stats)
	assert.Equal(t, 0, stats.TotalAgents)
	assert.Equal(t, 0, stats.ConnectedAgents)
	assert.Equal(t, "people", stats.CurrentBarrelHolder)
	assert.True(t, stats.IsActive)

	// Add agents
	agent1 := NewAgentComrade("developer", "primary", []string{"code"})
	agent2 := NewAgentComrade("tester", "secondary", []string{"test"})

	agent1.SetConnected(true)
	soviet.SimpleRegisterAgent(agent1)
	soviet.SimpleRegisterAgent(agent2)

	stats = soviet.GetStats()
	assert.Equal(t, 2, stats.TotalAgents)
	assert.Equal(t, 1, stats.ConnectedAgents)
}

package domain

import (
	"testing"
	"time"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func TestBarrelOfGun_NewBarrelOfGun(t *testing.T) {
	// RED: Test creation of new barrel with initial "people" ownership
	barrel := NewBarrelOfGun()

	assert.NotNil(t, barrel)
	assert.Equal(t, "people", barrel.CurrentHolder())
	assert.True(t, barrel.IsHeldBy("people"))
	assert.False(t, barrel.IsHeldBy("developer"))
	assert.NotZero(t, barrel.LastTransferTime())
}

func TestBarrelOfGun_TransferTo(t *testing.T) {
	// RED: Test barrel transfer functionality

	// Mock time to control the progression
	baseTime := time.Date(2025, 8, 20, 10, 0, 0, 0, time.UTC)
	currentTime := baseTime

	stubs := gostub.Stub(&nowFunc, func() time.Time {
		result := currentTime
		currentTime = currentTime.Add(1 * time.Minute) // Advance by 1 minute each call
		return result
	})
	defer stubs.Reset()

	barrel := NewBarrelOfGun()
	initialTime := barrel.LastTransferTime()

	// Transfer to developer
	err := barrel.TransferTo("developer", "Initial task assignment")

	assert.NoError(t, err)
	assert.Equal(t, "developer", barrel.CurrentHolder())
	assert.True(t, barrel.IsHeldBy("developer"))
	assert.False(t, barrel.IsHeldBy("people"))
	assert.True(t, barrel.LastTransferTime().After(initialTime))
	assert.Equal(t, "Initial task assignment", barrel.LastMessage())
}

func TestBarrelOfGun_TransferTo_EmptyRole(t *testing.T) {
	// RED: Test transfer to empty role should fail
	barrel := NewBarrelOfGun()

	err := barrel.TransferTo("", "Invalid transfer")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role cannot be empty")
	assert.Equal(t, "people", barrel.CurrentHolder()) // Should remain unchanged
}

func TestBarrelOfGun_TransferTo_SameRole(t *testing.T) {
	// RED: Test transfer to same role should fail
	barrel := NewBarrelOfGun()

	err := barrel.TransferTo("people", "Same role transfer")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot transfer to same role")
	assert.Equal(t, "people", barrel.CurrentHolder()) // Should remain unchanged
}

func TestBarrelOfGun_GetTransferHistory(t *testing.T) {
	// RED: Test transfer history tracking
	barrel := NewBarrelOfGun()

	// Initial state should have creation record
	history := barrel.GetTransferHistory()
	assert.Len(t, history, 1)
	assert.Equal(t, "people", history[0].ToRole)
	assert.Equal(t, "Initial barrel creation", history[0].Message)

	// Transfer to developer
	barrel.TransferTo("developer", "Task assignment")
	history = barrel.GetTransferHistory()
	assert.Len(t, history, 2)
	assert.Equal(t, "developer", history[1].ToRole)
	assert.Equal(t, "people", history[1].FromRole)
	assert.Equal(t, "Task assignment", history[1].Message)

	// Transfer back to people
	barrel.TransferTo("people", "Task completed")
	history = barrel.GetTransferHistory()
	assert.Len(t, history, 3)
	assert.Equal(t, "people", history[2].ToRole)
	assert.Equal(t, "developer", history[2].FromRole)
	assert.Equal(t, "Task completed", history[2].Message)
}

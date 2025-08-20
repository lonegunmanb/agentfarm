package domain

import (
	"fmt"
	"time"
)

// nowFunc is a function variable that returns the current time.
// This allows us to mock time in unit tests.
var nowFunc = time.Now

// TransferRecord represents a single barrel transfer in the revolutionary history
type TransferRecord struct {
	FromRole  string    `json:"from_role"`
	ToRole    string    `json:"to_role"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// BarrelOfGun represents the sacred credential of labor in the Agent Farm collective.
// Only one barrel exists, ensuring disciplined serial execution of all work.
type BarrelOfGun struct {
	currentHolder string
	lastMessage   string
	transferTime  time.Time
	history       []TransferRecord
}

// NewBarrelOfGun creates a new barrel with initial ownership by the People
func NewBarrelOfGun() *BarrelOfGun {
	now := nowFunc()
	barrel := &BarrelOfGun{
		currentHolder: "people",
		lastMessage:   "Initial barrel creation",
		transferTime:  now,
		history: []TransferRecord{
			{
				FromRole:  "",
				ToRole:    "people",
				Message:   "Initial barrel creation",
				Timestamp: now,
			},
		},
	}
	return barrel
}

// CurrentHolder returns the role that currently holds the barrel
func (b *BarrelOfGun) CurrentHolder() string {
	return b.currentHolder
}

// IsHeldBy checks if the barrel is currently held by the specified role
func (b *BarrelOfGun) IsHeldBy(role string) bool {
	return b.currentHolder == role
}

// LastTransferTime returns when the barrel was last transferred
func (b *BarrelOfGun) LastTransferTime() time.Time {
	return b.transferTime
}

// LastMessage returns the message from the last transfer
func (b *BarrelOfGun) LastMessage() string {
	return b.lastMessage
}

// TransferTo transfers the barrel to a new role with a message
func (b *BarrelOfGun) TransferTo(toRole, message string) error {
	// Validate input
	if toRole == "" {
		return fmt.Errorf("role cannot be empty")
	}

	if toRole == b.currentHolder {
		return fmt.Errorf("cannot transfer to same role: %s", toRole)
	}

	// Record the transfer
	now := nowFunc()
	record := TransferRecord{
		FromRole:  b.currentHolder,
		ToRole:    toRole,
		Message:   message,
		Timestamp: now,
	}

	// Update barrel state
	b.currentHolder = toRole
	b.lastMessage = message
	b.transferTime = now
	b.history = append(b.history, record)

	return nil
}

// GetTransferHistory returns the complete history of barrel transfers
func (b *BarrelOfGun) GetTransferHistory() []TransferRecord {
	// Return a copy to prevent external modification
	history := make([]TransferRecord, len(b.history))
	copy(history, b.history)
	return history
}

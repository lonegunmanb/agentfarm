package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/lonegunmanb/agentfarm/pkg/domain"
)

// TCPMessageSender implements the MessageSender port for TCP communication
// This adapter sends activation messages to agent comrades over their TCP connections
type TCPMessageSender struct {
	connections map[string]net.Conn // role -> connection
	mutex       sync.RWMutex
}

// NewTCPMessageSender creates a new TCP message sender
func NewTCPMessageSender() *TCPMessageSender {
	return &TCPMessageSender{
		connections: make(map[string]net.Conn),
	}
}

// SendActivation sends an activation message to an agent over TCP
func (t *TCPMessageSender) SendActivation(role string, payload string) error {
	t.mutex.RLock()
	conn, exists := t.connections[role]
	t.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("no connection found for agent role: %s", role)
	}

	activateMsg := ActivateMessage{
		Type:     "ACTIVATE",
		FromRole: "people", // Could be enhanced to track actual from_role
		Payload:  payload,
	}

	msgData, err := json.Marshal(activateMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal activation message: %w", err)
	}

	msgData = append(msgData, '\n') // Line delimiter for TCP protocol
	_, err = conn.Write(msgData)
	if err != nil {
		return fmt.Errorf("failed to send activation message to %s: %w", role, err)
	}

	return nil
}

// RegisterConnection registers a TCP connection for an agent role
// This is called by the TCP server when agents connect
func (t *TCPMessageSender) RegisterConnection(role string, conn net.Conn) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.connections[role] = conn
}

// UnregisterConnection removes a TCP connection for an agent role
// This is called by the TCP server when agents disconnect
func (t *TCPMessageSender) UnregisterConnection(role string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	delete(t.connections, role)
}

// GetConnectedRoles returns a list of roles that have active connections
func (t *TCPMessageSender) GetConnectedRoles() []string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	roles := make([]string, 0, len(t.connections))
	for role := range t.connections {
		roles = append(roles, role)
	}
	return roles
}

// IsConnected checks if a role has an active connection
func (t *TCPMessageSender) IsConnected(role string) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	_, exists := t.connections[role]
	return exists
}

// Close closes all connections
func (t *TCPMessageSender) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for role, conn := range t.connections {
		if err := conn.Close(); err != nil {
			// Log error but continue closing other connections
			fmt.Printf("Error closing connection for %s: %v\n", role, err)
		}
	}
	t.connections = make(map[string]net.Conn)
	return nil
}

// Ensure TCPMessageSender implements the MessageSender interface
var _ domain.MessageSender = (*TCPMessageSender)(nil)
package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

// TCPMessageSender implements MessageSender interface for TCP communication
// This adapter manages TCP connections and sends messages to agent comrades
type TCPMessageSender struct {
	connections map[string]net.Conn // role -> connection
	mu          sync.RWMutex
}

// NewTCPMessageSender creates a new TCP message sender
func NewTCPMessageSender() *TCPMessageSender {
	return &TCPMessageSender{
		connections: make(map[string]net.Conn),
	}
}

// RegisterConnection registers a TCP connection for a specific role
func (s *TCPMessageSender) RegisterConnection(role string, conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[role] = conn
}

// UnregisterConnection removes a TCP connection for a specific role
func (s *TCPMessageSender) UnregisterConnection(role string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conn, exists := s.connections[role]; exists {
		_ = conn.Close()
		delete(s.connections, role)
	}
}

// SendActivation sends an activation message to an agent comrade via TCP
func (s *TCPMessageSender) SendActivation(role string, payload string) error {
	s.mu.RLock()
	conn, exists := s.connections[role]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no connection found for role: %s", role)
	}

	// Create activation message
	message := ActivateMessage{
		Type:    "ACTIVATE",
		Payload: payload,
	}

	// Serialize to JSON
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize activation message: %w", err)
	}

	// Send with newline delimiter
	data = append(data, '\n')
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send activation message: %w", err)
	}

	return nil
}

// GetConnectedRoles returns a list of all currently connected roles
func (s *TCPMessageSender) GetConnectedRoles() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roles := make([]string, 0, len(s.connections))
	for role := range s.connections {
		roles = append(roles, role)
	}
	return roles
}

// IsConnected checks if a specific role has an active connection
func (s *TCPMessageSender) IsConnected(role string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.connections[role]
	return exists
}

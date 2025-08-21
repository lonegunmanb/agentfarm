package tcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/lonegunmanb/agentfarm/pkg/domain"
)

// TCPServer implements the CommandHandler port for TCP communication
// This adapter handles incoming TCP connections and translates them to domain operations
type TCPServer struct {
	sovietService domain.SovietService
	agentService  domain.AgentService
	sender        domain.MessageSender
	logger        domain.Logger
	connections   map[string]net.Conn // role -> connection
	mu            sync.RWMutex
	port          int
	listener      net.Listener
}

// NewTCPServer creates a new TCP server adapter
func NewTCPServer(
	sovietService domain.SovietService,
	agentService domain.AgentService,
	sender domain.MessageSender,
	logger domain.Logger,
	port int,
) *TCPServer {
	return &TCPServer{
		sovietService: sovietService,
		agentService:  agentService,
		sender:        sender,
		logger:        logger,
		connections:   make(map[string]net.Conn),
		port:          port,
	}
}

// Start starts the TCP server and begins accepting connections
func (s *TCPServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}

	s.listener = listener
	s.logger.Info("TCP Server started", map[string]interface{}{
		"port": s.port,
	})

	go s.acceptConnections(ctx)
	return nil
}

// Stop stops the TCP server
func (s *TCPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// acceptConnections accepts incoming connections and handles them
func (s *TCPServer) acceptConnections(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					s.logger.Error("Failed to accept connection", map[string]interface{}{
						"error": err.Error(),
					})
				}
				continue
			}

			go s.handleConnection(ctx, conn)
		}
	}
}

// handleConnection handles a single TCP connection
func (s *TCPServer) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		s.processMessage(ctx, conn, line)
	}

	if err := scanner.Err(); err != nil {
		s.logger.Error("Connection scan error", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// processMessage processes a single JSON message from a connection
func (s *TCPServer) processMessage(ctx context.Context, conn net.Conn, messageData string) {
	s.logger.Debug("Received message", map[string]interface{}{
		"message": messageData,
	})

	// Parse base message to determine type
	var baseMsg TCPMessage
	if err := json.Unmarshal([]byte(messageData), &baseMsg); err != nil {
		s.sendError(conn, "Invalid JSON format")
		return
	}

	switch baseMsg.Type {
	case "REGISTER":
		s.handleRegisterMessage(ctx, conn, messageData)
	case "YIELD":
		s.handleYieldMessage(ctx, conn, messageData)
	case "QUERY_AGENTS":
		s.handleQueryAgentsMessage(ctx, conn)
	case "QUERY_STATUS":
		s.handleQueryStatusMessage(ctx, conn)
	default:
		s.sendError(conn, fmt.Sprintf("Unknown message type: %s", baseMsg.Type))
	}
}

// Implementation of CommandHandler interface methods

// HandleRegister processes agent registration requests
func (s *TCPServer) HandleRegister(ctx context.Context, role string, capabilities []string) (bool, string, error) {
	agent := domain.NewAgentComrade(role, capabilities)
	return s.sovietService.RegisterAgent(agent)
}

// HandleYield processes yield requests from agents or people
func (s *TCPServer) HandleYield(ctx context.Context, fromRole, toRole, payload string) error {
	yieldMsg := domain.NewYieldMessage(fromRole, toRole, payload)
	return s.sovietService.ProcessYield(yieldMsg)
}

// HandleQueryAgents processes status query requests
func (s *TCPServer) HandleQueryAgents(ctx context.Context) ([]string, error) {
	return s.agentService.GetRegisteredAgents(), nil
}

// HandleQueryStatus processes detailed status query requests
func (s *TCPServer) HandleQueryStatus(ctx context.Context) (domain.StatusResponse, error) {
	status := s.sovietService.QueryStatus()
	return status, nil
}

// TCP-specific message handlers

func (s *TCPServer) handleRegisterMessage(ctx context.Context, conn net.Conn, messageData string) {
	var msg RegisterMessage
	if err := json.Unmarshal([]byte(messageData), &msg); err != nil {
		s.sendError(conn, "Invalid REGISTER message format")
		return
	}

	if msg.Role == "" {
		s.sendError(conn, "Role is required for registration")
		return
	}

	capabilities := msg.Capabilities
	if capabilities == nil {
		capabilities = []string{}
	}

	// Store connection for this role
	s.mu.Lock()
	s.connections[msg.Role] = conn
	s.mu.Unlock()

	shouldActivate, payload, err := s.HandleRegister(ctx, msg.Role, capabilities)
	if err != nil {
		s.sendError(conn, err.Error())
		return
	}

	// Send registration acknowledgment
	ackMsg := AckRegisterMessage{
		Type:    "ACK_REGISTER",
		Status:  "success",
		Message: fmt.Sprintf("Comrade '%s' successfully enlisted in the collective.", msg.Role),
	}
	s.sendMessage(conn, ackMsg)

	// If should activate, send activation message
	if shouldActivate {
		activateMsg := ActivateMessage{
			Type:     "ACTIVATE",
			FromRole: "soviet", // Will be set properly based on actual from role
			Payload:  payload,
		}
		s.sendMessage(conn, activateMsg)
	}
}

func (s *TCPServer) handleYieldMessage(ctx context.Context, conn net.Conn, messageData string) {
	var msg YieldMessage
	if err := json.Unmarshal([]byte(messageData), &msg); err != nil {
		s.sendError(conn, "Invalid YIELD message format")
		return
	}

	if msg.FromRole == "" || msg.ToRole == "" {
		s.sendError(conn, "FromRole and ToRole are required for yield")
		return
	}

	err := s.HandleYield(ctx, msg.FromRole, msg.ToRole, msg.Payload)
	if err != nil {
		s.sendError(conn, err.Error())
		return
	}

	// If yielding to an agent, send activation message
	if msg.ToRole != "people" {
		s.mu.RLock()
		targetConn, exists := s.connections[msg.ToRole]
		s.mu.RUnlock()

		if exists {
			activateMsg := ActivateMessage{
				Type:     "ACTIVATE",
				FromRole: msg.FromRole,
				Payload:  msg.Payload,
			}
			s.sendMessage(targetConn, activateMsg)
		}
	}
}

func (s *TCPServer) handleQueryAgentsMessage(ctx context.Context, conn net.Conn) {
	details := s.agentService.GetAgentDetails()
	
	// Convert domain.AgentDetails to TCP protocol format
	agentDetails := make([]AgentDetailInfo, len(details))
	for i, detail := range details {
		agentDetails[i] = AgentDetailInfo{
			Role:         detail.Role,
			Capabilities: detail.Capabilities,
			State:        detail.State.String(),
			Connected:    detail.Connected,
		}
	}

	response := AgentDetailsMessage{
		Type:         "AGENT_DETAILS",
		AgentDetails: agentDetails,
	}
	s.sendMessage(conn, response)
}

func (s *TCPServer) handleQueryStatusMessage(ctx context.Context, conn net.Conn) {
	status, err := s.HandleQueryStatus(ctx)
	if err != nil {
		s.sendError(conn, err.Error())
		return
	}

	// Convert domain.AgentState to string for TCP protocol
	agentStates := make(map[string]string)
	for role, state := range status.AgentStates {
		agentStates[role] = state.String()
	}

	response := StatusMessage{
		Type:             "STATUS",
		BarrelHolder:     status.BarrelHolder,
		RegisteredAgents: status.RegisteredAgents,
		AgentStates:      agentStates,
		ConnectedAgents:  status.ConnectedAgents,
	}
	s.sendMessage(conn, response)
}

func (s *TCPServer) sendError(conn net.Conn, message string) {
	errorMsg := ErrorMessage{
		Type:    "ERROR",
		Message: message,
	}
	s.sendMessage(conn, errorMsg)
}

func (s *TCPServer) sendMessage(conn net.Conn, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lonegunmanb/agentfarm/pkg/adapters/tcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

const (
	defaultServerAddr = "localhost:53646"
	connectionTimeout = 10 * time.Second
)

// AgentFarmMCPServer provides MCP tools for Agent Farm operations
type AgentFarmMCPServer struct {
	serverAddr string
	mu         sync.RWMutex
	connections map[string]*AgentConnection // role -> connection
}

// AgentConnection represents a connection to the Agent Farm server for a specific role
type AgentConnection struct {
	conn     net.Conn
	role     string
	active   bool
	doneChan chan bool
	mu       sync.RWMutex
}

func main() {
	serverAddr := flag.String("server", defaultServerAddr, "Agent Farm server address")
	mode := flag.String("mode", getenv("TRANSPORT_MODE", "stdio"), "transport mode: stdio or streamable-http")
	host := flag.String("host", getenv("TRANSPORT_HOST", "127.0.0.1"), "host for streamable-http server")
	port := flag.String("port", getenv("TRANSPORT_PORT", "8080"), "port for streamable-http server")
	flag.Parse()

	agentFarmServer := &AgentFarmMCPServer{
		serverAddr:  *serverAddr,
		connections: make(map[string]*AgentConnection),
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "agent-farm-mcp",
		Version: "1.0.0",
	}, nil)

	registerMCPTools(server, agentFarmServer)

	switch *mode {
	case "stdio":
		if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
			log.Fatal(err)
		}
	case "streamable-http":
		addr := fmt.Sprintf("%s:%s", *host, *port)
		log.Printf("MCP server serving at %s", addr)
		handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
			return server
		})
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatalf("failed to start streamable-http server: %v", err)
		}
	default:
		log.Fatalf("unknown mode: %s", *mode)
	}
}

// Input/Output types for MCP tools
type RegisterAgentInput struct {
	Role string `json:"role" jsonschema:"You role to register (e.g., 'developer', 'tester', 'reviewer'). This identifies the agent's specialty in the revolutionary collective.,required"`
}

type RegisterAgentOutput struct {
	Status   string `json:"status" jsonschema:"The status of the registration operation, typically 'activated' when the barrel is received"`
	Role     string `json:"role" jsonschema:"Your role"`
	FromRole string `json:"from_role,omitempty" jsonschema:"The role of the agent or 'people' who yielded the barrel to you"`
	Message  string `json:"message" jsonschema:"The command that was sent to you, read it carefully."`
}

type YieldBarrelInput struct {
	FromRole string `json:"from_role" jsonschema:"Your role, so the next agent could know who send the command to. Required."`
	ToRole   string `json:"to_role" jsonschema:"The target role to yield the barrel to. Can be another agent role (e.g., 'tester') or 'people' to return control to the People's representatives,required"`
	Message  string `json:"message" jsonschema:"Message to send with the yield operation. This should explain what work was completed and what the next agent should do,required"`
}

type YieldBarrelOutput struct {
	Status   string `json:"status" jsonschema:"The status of the yield operation, typically 'yielded' when successful"`
	FromRole string `json:"from_role" jsonschema:"Your role."`
	ToRole   string `json:"to_role" jsonschema:"The role of the agent or 'people' who received the barrel"`
	Message  string `json:"message" jsonschema:"Confirmation message about the barrel transfer"`
}

func registerMCPTools(server *mcp.Server, agentFarmServer *AgentFarmMCPServer) {
	// Register Agent Tool
	mcp.AddTool(server, &mcp.Tool{
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: p(false),
			IdempotentHint:  false,
			OpenWorldHint:   p(false),
			ReadOnlyHint:    false,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"role": {
					Type:        "string",
					Description: "The role of the agent comrade to register (e.g., 'developer', 'tester', 'reviewer'). This identifies the agent's specialty in the revolutionary collective.",
				},
			},
			Required: []string{"role"},
		},
		Description: "Register an agent comrade with the Central Committee and wait for barrel assignment. This tool blocks until the barrel is yielded to this agent, then returns the activation message. The agent will automatically return to waiting state after the tool completes.",
		Name:        "register_agent",
	}, func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[RegisterAgentInput]) (*mcp.CallToolResultFor[RegisterAgentOutput], error) {
		return agentFarmServer.RegisterAgentHandler(ctx, session, params)
	})

	// Yield Tool
	mcp.AddTool(server, &mcp.Tool{
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: p(false),
			IdempotentHint:  false,
			OpenWorldHint:   p(false),
			ReadOnlyHint:    false,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"from_role": {
					Type:        "string",
					Description: "The role of the agent comrade who currently holds the barrel and wants to yield it.",
				},
				"to_role": {
					Type:        "string",
					Description: "The target role to yield the barrel to. Can be another agent role (e.g., 'tester') or 'people' to return control to the People's representatives.",
				},
				"message": {
					Type:        "string",
					Description: "Message to send with the yield operation. This should explain what work was completed and what the next agent should do.",
				},
			},
			Required: []string{"from_role", "to_role", "message"},
		},
		Description: "Yield the barrel of gun from one agent to another. This tool blocks until the barrel has been successfully transferred and the from_role agent returns to waiting state. Use this to coordinate workflow between agent comrades.",
		Name:        "yield_barrel",
	}, func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[YieldBarrelInput]) (*mcp.CallToolResultFor[YieldBarrelOutput], error) {
		return agentFarmServer.YieldBarrelHandler(ctx, session, params)
	})
}

// RegisterAgentHandler handles the register agent MCP tool
func (afs *AgentFarmMCPServer) RegisterAgentHandler(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[RegisterAgentInput]) (*mcp.CallToolResultFor[RegisterAgentOutput], error) {
	role := params.Arguments.Role
	if role == "" {
		return nil, fmt.Errorf("role is required and must be a non-empty string")
	}

	// Create connection to Agent Farm server
	conn, err := afs.connectToServer()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Agent Farm server: %w", err)
	}

	agentConn := &AgentConnection{
		conn:     conn,
		role:     role,
		active:   false,
		doneChan: make(chan bool),
	}

	// Store connection
	afs.mu.Lock()
	afs.connections[role] = agentConn
	afs.mu.Unlock()

	// Cleanup on exit
	defer func() {
		afs.mu.Lock()
		delete(afs.connections, role)
		afs.mu.Unlock()
		conn.Close()
	}()

	// Register with the Central Committee
	registerMsg := tcp.RegisterMessage{
		Type: "REGISTER",
		Role: role,
	}

	if err := afs.sendMessage(conn, registerMsg); err != nil {
		return nil, fmt.Errorf("failed to register agent %s: %w", role, err)
	}

	// Listen for messages and wait for activation
	result, err := afs.waitForActivationTyped(ctx, agentConn)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResultFor[RegisterAgentOutput]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("🔥 Agent '%s' received barrel! From: %s, Message: %s", result.Role, result.FromRole, result.Message),
			},
		},
	}, nil
}

// YieldBarrelHandler handles the yield barrel MCP tool
func (afs *AgentFarmMCPServer) YieldBarrelHandler(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[YieldBarrelInput]) (*mcp.CallToolResultFor[YieldBarrelOutput], error) {
	fromRole := params.Arguments.FromRole
	toRole := params.Arguments.ToRole
	message := params.Arguments.Message

	if fromRole == "" {
		return nil, fmt.Errorf("from_role is required and must be a non-empty string")
	}

	if toRole == "" {
		return nil, fmt.Errorf("to_role is required and must be a non-empty string")
	}

	// Get the connection for the from_role
	afs.mu.RLock()
	agentConn, exists := afs.connections[fromRole]
	afs.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("agent %s is not registered", fromRole)
	}

	agentConn.mu.RLock()
	if !agentConn.active {
		agentConn.mu.RUnlock()
		return nil, fmt.Errorf("agent %s is not currently active (does not hold the barrel)", fromRole)
	}
	agentConn.mu.RUnlock()

	// Send yield message
	yieldMsg := tcp.YieldMessage{
		Type:     "YIELD",
		FromRole: fromRole,
		ToRole:   toRole,
		Payload:  message,
	}

	if err := afs.sendMessage(agentConn.conn, yieldMsg); err != nil {
		return nil, fmt.Errorf("failed to yield barrel from %s to %s: %w", fromRole, toRole, err)
	}

	// Wait for the agent to become inactive (barrel yielded)
	result, err := afs.waitForDeactivationTyped(ctx, agentConn, toRole)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResultFor[YieldBarrelOutput]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("✅ Barrel successfully yielded from '%s' to '%s'", result.FromRole, result.ToRole),
			},
		},
	}, nil
}

// waitForActivationTyped waits for activation and returns typed result
func (afs *AgentFarmMCPServer) waitForActivationTyped(ctx context.Context, agentConn *AgentConnection) (*RegisterAgentOutput, error) {
	// Start message handler goroutine
	messageChan := make(chan interface{}, 10)
	errorChan := make(chan error, 1)

	go afs.handleMessages(agentConn, messageChan, errorChan)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errorChan:
			return nil, err
		case msg := <-messageChan:
			switch m := msg.(type) {
			case tcp.ActivateMessage:
				agentConn.mu.Lock()
				agentConn.active = true
				agentConn.mu.Unlock()
				
				result := &RegisterAgentOutput{
					Status:   "activated",
					Role:     agentConn.role,
					FromRole: m.FromRole,
					Message:  m.Payload,
				}
				return result, nil
			case tcp.AckRegisterMessage:
				if m.Status != "success" {
					return nil, fmt.Errorf("registration failed: %s", m.Message)
				}
				// Continue waiting for activation
			case tcp.ErrorMessage:
				return nil, fmt.Errorf("error from Central Committee: %s", m.Message)
			}
		}
	}
}

// waitForDeactivationTyped waits for deactivation and returns typed result
func (afs *AgentFarmMCPServer) waitForDeactivationTyped(ctx context.Context, agentConn *AgentConnection, toRole string) (*YieldBarrelOutput, error) {
	// Since yielding is immediate and agents return to waiting state,
	// we just need to wait a moment and check the state
	time.Sleep(100 * time.Millisecond)
	
	agentConn.mu.Lock()
	agentConn.active = false // Agent returns to waiting state after yield
	agentConn.mu.Unlock()

	result := &YieldBarrelOutput{
		Status:   "yielded",
		FromRole: agentConn.role,
		ToRole:   toRole,
		Message:  "Barrel successfully yielded",
	}
	return result, nil
}

// RegisterAgent registers an agent and blocks until it receives the barrel
func (afs *AgentFarmMCPServer) RegisterAgent(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	role, ok := args["role"].(string)
	if !ok || role == "" {
		return nil, fmt.Errorf("role is required and must be a string")
	}

	// Create connection to Agent Farm server
	conn, err := afs.connectToServer()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Agent Farm server: %w", err)
	}

	agentConn := &AgentConnection{
		conn:     conn,
		role:     role,
		active:   false,
		doneChan: make(chan bool),
	}

	// Store connection
	afs.mu.Lock()
	afs.connections[role] = agentConn
	afs.mu.Unlock()

	// Cleanup on exit
	defer func() {
		afs.mu.Lock()
		delete(afs.connections, role)
		afs.mu.Unlock()
		conn.Close()
	}()

	// Register with the Central Committee
	registerMsg := tcp.RegisterMessage{
		Type: "REGISTER",
		Role: role,
	}

	if err := afs.sendMessage(conn, registerMsg); err != nil {
		return nil, fmt.Errorf("failed to register agent %s: %w", role, err)
	}

	// Listen for messages and wait for activation
	return afs.waitForActivation(ctx, agentConn)
}

// YieldBarrel yields the barrel from one agent to another
func (afs *AgentFarmMCPServer) YieldBarrel(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	fromRole, ok := args["from_role"].(string)
	if !ok || fromRole == "" {
		return nil, fmt.Errorf("from_role is required and must be a string")
	}

	toRole, ok := args["to_role"].(string)
	if !ok || toRole == "" {
		return nil, fmt.Errorf("to_role is required and must be a string")
	}

	message, ok := args["message"].(string)
	if !ok {
		message = ""
	}

	// Get the connection for the from_role
	afs.mu.RLock()
	agentConn, exists := afs.connections[fromRole]
	afs.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("agent %s is not registered", fromRole)
	}

	agentConn.mu.RLock()
	if !agentConn.active {
		agentConn.mu.RUnlock()
		return nil, fmt.Errorf("agent %s is not currently active (does not hold the barrel)", fromRole)
	}
	agentConn.mu.RUnlock()

	// Send yield message
	yieldMsg := tcp.YieldMessage{
		Type:     "YIELD",
		FromRole: fromRole,
		ToRole:   toRole,
		Payload:  message,
	}

	if err := afs.sendMessage(agentConn.conn, yieldMsg); err != nil {
		return nil, fmt.Errorf("failed to yield barrel from %s to %s: %w", fromRole, toRole, err)
	}

	// Wait for the agent to become inactive (barrel yielded)
	return afs.waitForDeactivation(ctx, agentConn, toRole)
}

// connectToServer establishes a connection to the Agent Farm server
func (afs *AgentFarmMCPServer) connectToServer() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", afs.serverAddr, connectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", afs.serverAddr, err)
	}
	return conn, nil
}

// sendMessage sends a JSON message over the connection
func (afs *AgentFarmMCPServer) sendMessage(conn net.Conn, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	data = append(data, '\n')
	_, err = conn.Write(data)
	return err
}

// waitForActivation waits for the agent to receive the barrel and become active
func (afs *AgentFarmMCPServer) waitForActivation(ctx context.Context, agentConn *AgentConnection) (interface{}, error) {
	// Start message handler goroutine
	messageChan := make(chan interface{}, 10)
	errorChan := make(chan error, 1)

	go afs.handleMessages(agentConn, messageChan, errorChan)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errorChan:
			return nil, err
		case msg := <-messageChan:
			switch m := msg.(type) {
			case tcp.ActivateMessage:
				agentConn.mu.Lock()
				agentConn.active = true
				agentConn.mu.Unlock()
				
				result := map[string]interface{}{
					"status":    "activated",
					"role":      agentConn.role,
					"from_role": m.FromRole,
					"message":   m.Payload,
				}
				return result, nil
			case tcp.AckRegisterMessage:
				if m.Status != "success" {
					return nil, fmt.Errorf("registration failed: %s", m.Message)
				}
				// Continue waiting for activation
			case tcp.ErrorMessage:
				return nil, fmt.Errorf("error from Central Committee: %s", m.Message)
			}
		}
	}
}

// waitForDeactivation waits for the agent to yield the barrel and become inactive
func (afs *AgentFarmMCPServer) waitForDeactivation(ctx context.Context, agentConn *AgentConnection, toRole string) (interface{}, error) {
	// Since yielding is immediate and agents return to waiting state,
	// we just need to wait a moment and check the state
	time.Sleep(100 * time.Millisecond)
	
	agentConn.mu.Lock()
	agentConn.active = false // Agent returns to waiting state after yield
	agentConn.mu.Unlock()

	result := map[string]interface{}{
		"status":    "yielded",
		"from_role": agentConn.role,
		"to_role":   toRole,
		"message":   "Barrel successfully yielded",
	}
	return result, nil
}

// handleMessages processes incoming messages from the Agent Farm server
func (afs *AgentFarmMCPServer) handleMessages(agentConn *AgentConnection, messageChan chan interface{}, errorChan chan error) {
	defer close(messageChan)
	defer close(errorChan)

	scanner := make(chan string)
	go func() {
		defer close(scanner)
		buffer := make([]byte, 4096)
		remainder := ""
		
		for {
			n, err := agentConn.conn.Read(buffer)
			if err != nil {
				errorChan <- fmt.Errorf("connection read error: %w", err)
				return
			}
			
			data := remainder + string(buffer[:n])
			lines := strings.Split(data, "\n")
			remainder = lines[len(lines)-1]
			
			for _, line := range lines[:len(lines)-1] {
				line = strings.TrimSpace(line)
				if line != "" {
					scanner <- line
				}
			}
		}
	}()

	for line := range scanner {
		if err := afs.parseAndSendMessage(line, messageChan); err != nil {
			errorChan <- err
			return
		}
	}
}

// parseAndSendMessage parses a JSON message and sends it to the message channel
func (afs *AgentFarmMCPServer) parseAndSendMessage(line string, messageChan chan interface{}) error {
	// Parse the message to determine type
	var baseMsg struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal([]byte(line), &baseMsg); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	switch baseMsg.Type {
	case "ACTIVATE":
		var activateMsg tcp.ActivateMessage
		if err := json.Unmarshal([]byte(line), &activateMsg); err != nil {
			return fmt.Errorf("failed to parse ACTIVATE message: %w", err)
		}
		messageChan <- activateMsg
	case "ERROR":
		var errorMsg tcp.ErrorMessage
		if err := json.Unmarshal([]byte(line), &errorMsg); err != nil {
			return fmt.Errorf("failed to parse ERROR message: %w", err)
		}
		messageChan <- errorMsg
	case "ACK_REGISTER":
		var ackMsg tcp.AckRegisterMessage
		if err := json.Unmarshal([]byte(line), &ackMsg); err != nil {
			return fmt.Errorf("failed to parse ACK_REGISTER message: %w", err)
		}
		messageChan <- ackMsg
	}

	return nil
}

// Helper function to get pointer to value
func p[T any](input T) *T {
	return &input
}

// getenv returns environment variable value or fallback
func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

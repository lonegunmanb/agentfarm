package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lonegunmanb/agentfarm/pkg/adapters/tcp"
)

const (
	defaultServerAddr = "localhost:53646"
	connectionTimeout = 10 * time.Second
	reconnectDelay    = 5 * time.Second
)

// AgentClient represents an Agent Comrade connection to the Central Committee
type AgentClient struct {
	role       string
	serverAddr string
	yieldTo    string
	yieldMsg   string
	conn       net.Conn
	done       chan bool
}

func main() {
	var (
		role       = flag.String("role", "", "Agent comrade role (required)")
		serverAddr = flag.String("server", defaultServerAddr, "Soviet server address")
		yieldTo    = flag.String("yield-to", "", "Target role to yield barrel to after activation")
		yieldMsg   = flag.String("yield-msg", "", "Message to send with yield")
		help       = flag.Bool("help", false, "Show help")
		version    = flag.Bool("version", false, "Show version")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *version {
		showVersion()
		return
	}

	if *role == "" {
		fmt.Fprintf(os.Stderr, "Error: --role is required\n")
		showHelp()
		os.Exit(1)
	}

	client := &AgentClient{
		role:       *role,
		serverAddr: *serverAddr,
		yieldTo:    *yieldTo,
		yieldMsg:   *yieldMsg,
		done:       make(chan bool),
	}

	if err := client.Run(); err != nil {
		log.Fatalf("Agent comrade %s failed: %v", *role, err)
	}
}

func (ac *AgentClient) Run() error {
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Printf("\nAgent comrade %s received shutdown signal, disconnecting...\n", ac.role)
		ac.done <- true
	}()

	for {
		select {
		case <-ac.done:
			if ac.conn != nil {
				ac.conn.Close()
			}
			return nil
		default:
			if err := ac.connectAndServe(); err != nil {
				fmt.Printf("Connection lost: %v. Reconnecting in %v...\n", err, reconnectDelay)
				time.Sleep(reconnectDelay)
				continue
			}
		}
	}
}

func (ac *AgentClient) connectAndServe() error {
	// Establish connection to Central Committee
	var err error
	ac.conn, err = net.DialTimeout("tcp", ac.serverAddr, connectionTimeout)
	if err != nil {
		return fmt.Errorf("failed to connect to Soviet server at %s: %w", ac.serverAddr, err)
	}
	defer ac.conn.Close()

	fmt.Printf("Agent comrade %s connected to Central Committee at %s\n", ac.role, ac.serverAddr)

	// Send registration message
	registerMsg := tcp.RegisterMessage{
		Type: "REGISTER",
		Role: ac.role,
	}

	if err := ac.sendMessage(registerMsg); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	fmt.Printf("Agent comrade %s registered successfully. Waiting for barrel assignment...\n", ac.role)

	// Listen for messages from Central Committee
	scanner := bufio.NewScanner(ac.conn)
	for scanner.Scan() {
		select {
		case <-ac.done:
			return nil
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if err := ac.handleMessage(line); err != nil {
			fmt.Printf("Error handling message: %v\n", err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	return fmt.Errorf("connection closed by server")
}

func (ac *AgentClient) handleMessage(line string) error {
	// Parse the message to determine type
	var baseMsg struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal([]byte(line), &baseMsg); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	switch baseMsg.Type {
	case "ACTIVATE":
		return ac.handleActivateMessage(line)
	case "ERROR":
		return ac.handleErrorMessage(line)
	case "ACK_REGISTER":
		return ac.handleAckRegisterMessage(line)
	default:
		fmt.Printf("Received unknown message type: %s\n", baseMsg.Type)
	}

	return nil
}

func (ac *AgentClient) handleActivateMessage(line string) error {
	var activateMsg tcp.ActivateMessage
	if err := json.Unmarshal([]byte(line), &activateMsg); err != nil {
		return fmt.Errorf("failed to parse ACTIVATE message: %w", err)
	}

	fmt.Printf("\n🔥 BARREL RECEIVED! Agent comrade %s is now active!\n", ac.role)
	if activateMsg.Payload != "" {
		fmt.Printf("📜 Message: %s\n", activateMsg.Payload)
	}

	// If yield-to is specified, yield the barrel before exiting
	if ac.yieldTo != "" {
		fmt.Printf("⚡ Auto-yielding barrel to: %s\n", ac.yieldTo)
		if err := ac.yieldBarrel(); err != nil {
			fmt.Printf("❌ Failed to yield barrel: %v\n", err)
		}
	}

	// Exit immediately after receiving activation
	fmt.Printf("✅ Agent comrade %s task completed. Exiting...\n", ac.role)
	os.Exit(0)
	return nil // This line will never be reached, but satisfies the function signature
}

func (ac *AgentClient) handleErrorMessage(line string) error {
	var errorMsg tcp.ErrorMessage
	if err := json.Unmarshal([]byte(line), &errorMsg); err != nil {
		return fmt.Errorf("failed to parse ERROR message: %w", err)
	}

	fmt.Printf("❌ Error from Central Committee: %s\n", errorMsg.Message)
	return nil
}

func (ac *AgentClient) handleAckRegisterMessage(line string) error {
	var ackMsg tcp.AckRegisterMessage
	if err := json.Unmarshal([]byte(line), &ackMsg); err != nil {
		return fmt.Errorf("failed to parse ACK_REGISTER message: %w", err)
	}

	fmt.Printf("📋 Registration acknowledged: %s\n", ackMsg.Message)
	if ackMsg.Status == "success" {
		fmt.Printf("✅ Agent comrade %s successfully enrolled in the collective\n", ac.role)
	} else {
		fmt.Printf("⚠️  Registration status: %s\n", ackMsg.Status)
	}
	return nil
}

func (ac *AgentClient) yieldBarrel() error {
	yieldMsg := tcp.YieldMessage{
		Type:     "YIELD",
		FromRole: ac.role,
		ToRole:   ac.yieldTo,
		Payload:  ac.yieldMsg,
	}

	if err := ac.sendMessage(yieldMsg); err != nil {
		return fmt.Errorf("failed to yield barrel: %w", err)
	}

	fmt.Printf("✅ Barrel successfully yielded to %s\n", ac.yieldTo)
	fmt.Printf("⏳ Agent comrade %s returned to waiting state.\n", ac.role)
	return nil
}

func (ac *AgentClient) sendMessage(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	data = append(data, '\n')
	_, err = ac.conn.Write(data)
	return err
}

func showHelp() {
	fmt.Printf(`Agent Farm - Agent Comrade CLI

USAGE:
    agent [OPTIONS]

OPTIONS:
    --role <role>           Agent comrade role (required)
    --server <address>      Soviet server address (default: %s)
    --yield-to <role>       Target role to yield barrel to after activation
    --yield-msg <message>   Message to send with yield
    --help                  Show this help
    --version               Show version

EXAMPLES:
    # Register as developer and wait for barrel
    agent --role=developer

    # Register as developer and auto-yield to tester with message
    agent --role=developer --yield-to=tester --yield-msg="Code ready for testing"

    # Connect to custom server
    agent --role=developer --server=localhost:8080

REVOLUTIONARY WORKFLOW:
    1. Agent comrade connects to Central Committee
    2. Registers with specified role
    3. Waits in disciplined formation for barrel assignment
    4. When barrel is received, prints activation message and exits immediately
    5. If --yield-to specified, yields barrel to target before exiting
    6. Process completes its revolutionary duty and terminates

The agent will automatically reconnect if connection is lost before activation.
Use Ctrl+C to gracefully disconnect while waiting for barrel assignment.
`, defaultServerAddr)
}

func showVersion() {
	fmt.Println("Agent Farm - Agent Comrade CLI v1.0")
	fmt.Println("Revolutionary Multi-agent Control Protocol")
	fmt.Println("Part of the Agent Farm collective")
}

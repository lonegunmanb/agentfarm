package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/lonegunmanb/agentfarm/pkg/adapters/tcp"
)

const (
	defaultServerAddr = "localhost:53646"
	connectionTimeout = 10 * time.Second
)

// PeopleClient represents the People's Representatives interface to the Central Committee
type PeopleClient struct {
	serverAddr string
	conn       net.Conn
}

func main() {
	var (
		serverAddr = flag.String("server", defaultServerAddr, "Soviet server address")
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

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Command is required\n")
		showHelp()
		os.Exit(1)
	}

	client := &PeopleClient{
		serverAddr: *serverAddr,
	}

	if err := client.ExecuteCommand(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (pc *PeopleClient) ExecuteCommand(args []string) error {
	command := args[0]

	switch command {
	case "yield":
		return pc.executeYield(args[1:])
	case "status":
		return pc.executeStatus()
	case "query-agents":
		return pc.executeQueryAgents()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (pc *PeopleClient) executeYield(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("yield command requires: yield <to_role> \"<message>\"")
	}

	toRole := args[0]
	message := strings.Join(args[1:], " ")
	
	// Remove quotes if present
	message = strings.Trim(message, `"'`)

	if err := pc.connect(); err != nil {
		return err
	}
	defer pc.conn.Close()

	yieldMsg := tcp.YieldMessage{
		Type:     "YIELD",
		FromRole: "people",
		ToRole:   toRole,
		Payload:  message,
	}

	if err := pc.sendMessage(yieldMsg); err != nil {
		return fmt.Errorf("failed to send yield command: %w", err)
	}

	fmt.Printf("✅ The People have yielded the barrel to comrade %s\n", toRole)
	if message != "" {
		fmt.Printf("📜 Message: %s\n", message)
	}

	return nil
}

func (pc *PeopleClient) executeStatus() error {
	if err := pc.connect(); err != nil {
		return err
	}
	defer pc.conn.Close()

	queryMsg := tcp.QueryMessage{
		Type: "QUERY_STATUS",
	}

	if err := pc.sendMessage(queryMsg); err != nil {
		return fmt.Errorf("failed to send status query: %w", err)
	}

	// Read the response
	scanner := bufio.NewScanner(pc.conn)
	if !scanner.Scan() {
		return fmt.Errorf("no response from server")
	}

	line := strings.TrimSpace(scanner.Text())
	if line == "" {
		return fmt.Errorf("empty response from server")
	}

	return pc.handleStatusResponse(line)
}

func (pc *PeopleClient) executeQueryAgents() error {
	if err := pc.connect(); err != nil {
		return err
	}
	defer pc.conn.Close()

	queryMsg := tcp.QueryMessage{
		Type: "QUERY_AGENTS",
	}

	if err := pc.sendMessage(queryMsg); err != nil {
		return fmt.Errorf("failed to send agent query: %w", err)
	}

	// Read the response
	scanner := bufio.NewScanner(pc.conn)
	if !scanner.Scan() {
		return fmt.Errorf("no response from server")
	}

	line := strings.TrimSpace(scanner.Text())
	if line == "" {
		return fmt.Errorf("empty response from server")
	}

	return pc.handleAgentListResponse(line)
}

func (pc *PeopleClient) connect() error {
	var err error
	pc.conn, err = net.DialTimeout("tcp", pc.serverAddr, connectionTimeout)
	if err != nil {
		return fmt.Errorf("failed to connect to Soviet server at %s: %w", pc.serverAddr, err)
	}
	return nil
}

func (pc *PeopleClient) sendMessage(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	data = append(data, '\n')
	_, err = pc.conn.Write(data)
	return err
}

func (pc *PeopleClient) handleStatusResponse(line string) error {
	var statusMsg tcp.StatusMessage
	if err := json.Unmarshal([]byte(line), &statusMsg); err != nil {
		// Try error message format
		var errorMsg tcp.ErrorMessage
		if errParse := json.Unmarshal([]byte(line), &errorMsg); errParse == nil {
			return fmt.Errorf("server error: %s", errorMsg.Message)
		}
		return fmt.Errorf("failed to parse status response: %w", err)
	}

	fmt.Println("🏛️  REVOLUTIONARY COLLECTIVE STATUS")
	fmt.Println("====================================")
	fmt.Printf("🔫 Barrel Holder: %s\n", statusMsg.BarrelHolder)
	fmt.Printf("👥 Registered Agents: %d\n", len(statusMsg.RegisteredAgents))

	if len(statusMsg.RegisteredAgents) > 0 {
		fmt.Println("\n📋 AGENT COMRADES:")
		for _, agent := range statusMsg.RegisteredAgents {
			state := "unknown"
			if s, exists := statusMsg.AgentStates[agent]; exists {
				state = s
			}
			
			connected := "❌ offline"
			if c, exists := statusMsg.ConnectedAgents[agent]; exists && c {
				connected = "✅ online"
			}
			
			icon := "⏳"
			if agent == statusMsg.BarrelHolder {
				icon = "🔥"
			}
			
			fmt.Printf("  %s %s - %s (%s)\n", icon, agent, state, connected)
		}
	} else {
		fmt.Println("\n📋 No agents registered in the collective")
	}

	fmt.Println("")
	return nil
}

func (pc *PeopleClient) handleAgentListResponse(line string) error {
	// Try to parse as detailed agent response first
	var agentDetailsMsg tcp.AgentDetailsMessage
	if err := json.Unmarshal([]byte(line), &agentDetailsMsg); err == nil {
		return pc.displayAgentDetails(agentDetailsMsg)
	}

	// Fallback to simple agent list (for backward compatibility)
	var agentListMsg tcp.AgentListMessage
	if err := json.Unmarshal([]byte(line), &agentListMsg); err == nil {
		return pc.displaySimpleAgentList(agentListMsg)
	}

	// Try error message format
	var errorMsg tcp.ErrorMessage
	if err := json.Unmarshal([]byte(line), &errorMsg); err == nil {
		return fmt.Errorf("server error: %s", errorMsg.Message)
	}

	return fmt.Errorf("failed to parse agent list response")
}

func (pc *PeopleClient) displayAgentDetails(msg tcp.AgentDetailsMessage) error {
	fmt.Println("👥 REGISTERED AGENT COMRADES")
	fmt.Println("============================")
	
	if len(msg.AgentDetails) > 0 {
		for i, agent := range msg.AgentDetails {
			icon := "⏳"
			if agent.State == "working" {
				icon = "🔥"
			}
			
			connected := "❌ offline"
			if agent.Connected {
				connected = "✅ online"
			}
			
			fmt.Printf("%d. %s %s - %s (%s)\n", i+1, icon, agent.Role, agent.State, connected)
			
			if len(agent.Capabilities) > 0 {
				fmt.Printf("   🛠️  Capabilities: %s\n", strings.Join(agent.Capabilities, ", "))
			} else {
				fmt.Printf("   🛠️  Capabilities: none specified\n")
			}
			fmt.Println()
		}
	} else {
		fmt.Println("No agents registered in the collective")
	}

	fmt.Printf("Total: %d comrades serving the People\n", len(msg.AgentDetails))
	return nil
}

func (pc *PeopleClient) displaySimpleAgentList(msg tcp.AgentListMessage) error {
	fmt.Println("👥 REGISTERED AGENT COMRADES")
	fmt.Println("============================")
	
	if len(msg.Agents) > 0 {
		for i, agent := range msg.Agents {
			fmt.Printf("%d. %s\n", i+1, agent)
		}
	} else {
		fmt.Println("No agents registered in the collective")
	}

	fmt.Printf("\nTotal: %d comrades serving the People\n", len(msg.Agents))
	return nil
}

func showHelp() {
	fmt.Printf(`Agent Farm - People's Representatives CLI

USAGE:
    people [OPTIONS] <command> [arguments...]

OPTIONS:
    --server <address>      Soviet server address (default: %s)
    --help                  Show this help
    --version               Show version

COMMANDS:
    yield <to_role> "<message>"     Transfer the barrel to specified agent comrade
    status                          Query comprehensive system status
    query-agents                    List all registered agent comrades

EXAMPLES:
    # Transfer barrel to developer with instructions
    people yield developer "Implement the authentication module"

    # Transfer barrel to tester
    people yield tester "Code ready for revolutionary testing"

    # Check complete system status
    people status

    # List all registered agents
    people query-agents

    # Connect to custom server
    people --server=localhost:8080 status

REVOLUTIONARY AUTHORITY:
    As representatives of the People, you have supreme authority over the collective.
    The barrel of gun serves the People's will, and all agent comrades await your guidance.
    Use this power responsibly to coordinate the revolutionary workflow.

PROTOCOL NOTES:
    - Commands execute immediately and disconnect
    - All agents respect the People's decisions unconditionally
    - Status commands provide transparency for the collective
    - The People can intervene at any point in the workflow
`, defaultServerAddr)
}

func showVersion() {
	fmt.Println("Agent Farm - People's Representatives CLI v1.0")
	fmt.Println("Revolutionary Multi-agent Control Protocol")
	fmt.Println("Supreme interface for the People's will")
}

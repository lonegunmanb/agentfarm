package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lonegunmanb/agentfarm/pkg/adapters/tcp"
	"github.com/lonegunmanb/agentfarm/pkg/domain"
)

const (
	defaultPort = 53646
)

func main() {
	// Parse command line flags
	var (
		port        = flag.Int("port", defaultPort, "TCP port for the Soviet server")
		debugMode   = flag.Bool("debug", false, "Enable debug logging")
		showHelp    = flag.Bool("help", false, "Show help message")
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *showHelp {
		showUsage()
		os.Exit(0)
	}

	if *showVersion {
		showVersionInfo()
		os.Exit(0)
	}

	// Create logger
	logger := domain.NewConsoleLogger(*debugMode)
	logger.Info("Starting Agent Farm Soviet Server", map[string]interface{}{
		"port":  *port,
		"debug": *debugMode,
	})

	// Create core domain components
	repository := domain.NewMemoryAgentRepository()
	barrel := domain.NewBarrelOfGun() // Initially held by the people
	soviet := domain.NewSovietState(repository)
	
	// Set the barrel in the soviet state
	if err := soviet.SetBarrel(barrel); err != nil {
		logger.Error("Failed to set barrel in soviet state", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	// Create message sender
	sender := tcp.NewTCPMessageSender()

	// Create TCP server adapter
	server := tcp.NewTCPServer(soviet, soviet, sender, logger, *port)

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	if err := server.Start(ctx); err != nil {
		logger.Error("Failed to start server", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	logger.Info("Agent Farm Soviet Server is running", map[string]interface{}{
		"port": *port,
		"status": "ready_for_agents",
	})
	logger.Info("Connect Agent Comrades via TCP", map[string]interface{}{
		"instructions": "Agents should connect to this port and register with their role",
	})
	logger.Info("People's representatives can connect via netcat", map[string]interface{}{
		"example": fmt.Sprintf("nc localhost %d", *port),
	})

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Received shutdown signal, gracefully stopping server...")

	// Stop the server
	if err := server.Stop(); err != nil {
		logger.Error("Error stopping server", map[string]interface{}{
			"error": err.Error(),
		})
	}

	logger.Info("Agent Farm Soviet Server stopped", map[string]interface{}{
		"status": "shutdown_complete",
	})
}

func showUsage() {
	fmt.Println("Agent Farm Soviet Server - Central Committee for Multi-agent Control Protocol")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Printf("  %s [options]\n", os.Args[0])
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Printf("  -port int\n\tTCP port for the Soviet server (default: %d)\n", defaultPort)
	fmt.Println("  -debug")
	fmt.Println("\tEnable debug logging")
	fmt.Println("  -help")
	fmt.Println("\tShow this help message")
	fmt.Println("  -version")
	fmt.Println("\tShow version information")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  The Soviet Server acts as the Central Committee managing the barrel of gun")
	fmt.Println("  and coordinating Agent Comrades in the revolutionary collective.")
	fmt.Println()
	fmt.Println("  Agent Comrades connect via TCP and register with their role.")
	fmt.Println("  People's representatives can connect via netcat or telnet for direct control.")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Printf("  # Start server on default port %d\n", defaultPort)
	fmt.Printf("  %s\n", os.Args[0])
	fmt.Println()
	fmt.Printf("  # Start server on custom port with debug logging\n")
	fmt.Printf("  %s -port 8080 -debug\n", os.Args[0])
	fmt.Println()
	fmt.Printf("  # Connect as People's representative\n")
	fmt.Printf("  nc localhost %d\n", defaultPort)
}

func showVersionInfo() {
	fmt.Println("Agent Farm Soviet Server")
	fmt.Println("Version: 4.0")
	fmt.Println("Date: August 20, 2025")
	fmt.Println("Revolutionary Multi-agent Control Protocol")
}

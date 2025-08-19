// Package config provides configuration management for the Agent Farm collective
package pkg

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config represents the revolutionary configuration for the Agent Farm collective
type Config struct {
	// Soviet (Central Committee) configuration
	SovietPort int
	SovietHost string

	// Logging configuration for the People's oversight
	LogLevel string

	// Agent Comrade configuration
	AgentReconnectTimeout int // seconds
}

// DefaultConfig returns the default revolutionary configuration
func DefaultConfig() *Config {
	return &Config{
		SovietPort:            53646, // Sacred port for the collective
		SovietHost:            "localhost",
		LogLevel:              "INFO",
		AgentReconnectTimeout: 30,
	}
}

// LoadConfig loads configuration from environment variables, falling back to defaults
func LoadConfig() *Config {
	config := DefaultConfig()

	// Soviet port configuration
	if portStr := os.Getenv("AGENT_FARM_SOVIET_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.SovietPort = port
		} else {
			log.Printf("Invalid AGENT_FARM_SOVIET_PORT: %s, using default %d", portStr, config.SovietPort)
		}
	}

	// Soviet host configuration
	if host := os.Getenv("AGENT_FARM_SOVIET_HOST"); host != "" {
		config.SovietHost = host
	}

	// Log level configuration
	if logLevel := os.Getenv("AGENT_FARM_LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	// Agent reconnect timeout
	if timeoutStr := os.Getenv("AGENT_FARM_RECONNECT_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			config.AgentReconnectTimeout = timeout
		} else {
			log.Printf("Invalid AGENT_FARM_RECONNECT_TIMEOUT: %s, using default %d", timeoutStr, config.AgentReconnectTimeout)
		}
	}

	return config
}

// GetSovietAddress returns the full address for the Soviet (Central Committee)
func (c *Config) GetSovietAddress() string {
	return fmt.Sprintf("%s:%d", c.SovietHost, c.SovietPort)
}

// Validate ensures the configuration serves the collective properly
func (c *Config) Validate() error {
	if c.SovietPort <= 0 || c.SovietPort > 65535 {
		return fmt.Errorf("invalid soviet port: %d (must be 1-65535)", c.SovietPort)
	}

	if c.SovietHost == "" {
		return fmt.Errorf("soviet host cannot be empty")
	}

	if c.AgentReconnectTimeout <= 0 {
		return fmt.Errorf("agent reconnect timeout must be positive: %d", c.AgentReconnectTimeout)
	}

	validLogLevels := map[string]bool{
		"DEBUG": true,
		"INFO":  true,
		"WARN":  true,
		"ERROR": true,
	}

	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s (must be DEBUG, INFO, WARN, or ERROR)", c.LogLevel)
	}

	return nil
}

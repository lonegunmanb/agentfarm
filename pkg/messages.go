// Package protocol defines the revolutionary communication messages for the Agent Farm collective
package pkg

import (
	"encoding/json"
	"fmt"
)

// MessageType represents the type of revolutionary messages in the collective
type MessageType string

const (
	// Agent Comrade -> Central Committee messages
	MsgTypeRegister MessageType = "REGISTER"
	MsgTypeYield    MessageType = "YIELD"

	// Central Committee -> Agent Comrade messages
	MsgTypeActivate    MessageType = "ACTIVATE"
	MsgTypeError       MessageType = "ERROR"
	MsgTypeAckRegister MessageType = "ACK_REGISTER"
)

// RegisterMessage represents an Agent Comrade's registration for collective service
type RegisterMessage struct {
	Type        MessageType `json:"type"`
	Role        string      `json:"role"`
	Description string      `json:"description"`
}

// YieldMessage represents the sacred transfer of the barrel of gun
type YieldMessage struct {
	Type    MessageType `json:"type"`
	ToRole  string      `json:"to_role"`
	Payload string      `json:"payload"`
}

// ActivateMessage represents orders from the Central Committee to begin work
type ActivateMessage struct {
	Type     MessageType `json:"type"`
	FromRole string      `json:"from_role"`
	Payload  string      `json:"payload"`
}

// ErrorMessage represents revolutionary error reporting
type ErrorMessage struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

// AckRegisterMessage represents confirmation of successful enlistment
type AckRegisterMessage struct {
	Type MessageType `json:"type"`
}

// ParseMessage parses a JSON message into the appropriate revolutionary message type
func ParseMessage(data []byte) (interface{}, error) {
	// First parse just the type field to determine message structure
	var typeParser struct {
		Type MessageType `json:"type"`
	}
	if err := json.Unmarshal(data, &typeParser); err != nil {
		return nil, fmt.Errorf("failed to parse message type: %w", err)
	}

	switch typeParser.Type {
	case MsgTypeRegister:
		var msg RegisterMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, fmt.Errorf("failed to parse REGISTER message: %w", err)
		}
		return &msg, nil

	case MsgTypeYield:
		var msg YieldMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, fmt.Errorf("failed to parse YIELD message: %w", err)
		}
		return &msg, nil

	case MsgTypeActivate:
		var msg ActivateMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, fmt.Errorf("failed to parse ACTIVATE message: %w", err)
		}
		return &msg, nil

	case MsgTypeError:
		var msg ErrorMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, fmt.Errorf("failed to parse ERROR message: %w", err)
		}
		return &msg, nil

	case MsgTypeAckRegister:
		var msg AckRegisterMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, fmt.Errorf("failed to parse ACK_REGISTER message: %w", err)
		}
		return &msg, nil

	default:
		return nil, fmt.Errorf("unknown message type: %s", typeParser.Type)
	}
}

// SerializeMessage serializes a revolutionary message to JSON
func SerializeMessage(msg interface{}) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}
	return data, nil
}

// ValidateMessage validates a revolutionary message according to collective principles
func ValidateMessage(msg interface{}) error {
	switch m := msg.(type) {
	case *RegisterMessage:
		return validateRegisterMessage(m)
	case *YieldMessage:
		return validateYieldMessage(m)
	case *ActivateMessage:
		return validateActivateMessage(m)
	case *ErrorMessage:
		return validateErrorMessage(m)
	case *AckRegisterMessage:
		return validateAckRegisterMessage(m)
	default:
		return fmt.Errorf("unknown message type for validation")
	}
}

func validateRegisterMessage(msg *RegisterMessage) error {
	if msg.Type != MsgTypeRegister {
		return fmt.Errorf("invalid message type for RegisterMessage: %s", msg.Type)
	}
	if msg.Role == "" {
		return fmt.Errorf("role cannot be empty in REGISTER message")
	}
	if msg.Role == "people" {
		return fmt.Errorf("'people' is a reserved role and cannot be registered as an agent")
	}
	return nil
}

func validateYieldMessage(msg *YieldMessage) error {
	if msg.Type != MsgTypeYield {
		return fmt.Errorf("invalid message type for YieldMessage: %s", msg.Type)
	}
	if msg.ToRole == "" {
		return fmt.Errorf("to_role cannot be empty in YIELD message")
	}
	return nil
}

func validateActivateMessage(msg *ActivateMessage) error {
	if msg.Type != MsgTypeActivate {
		return fmt.Errorf("invalid message type for ActivateMessage: %s", msg.Type)
	}
	if msg.FromRole == "" {
		return fmt.Errorf("from_role cannot be empty in ACTIVATE message")
	}
	return nil
}

func validateErrorMessage(msg *ErrorMessage) error {
	if msg.Type != MsgTypeError {
		return fmt.Errorf("invalid message type for ErrorMessage: %s", msg.Type)
	}
	if msg.Message == "" {
		return fmt.Errorf("message cannot be empty in ERROR message")
	}
	return nil
}

func validateAckRegisterMessage(msg *AckRegisterMessage) error {
	if msg.Type != MsgTypeAckRegister {
		return fmt.Errorf("invalid message type for AckRegisterMessage: %s", msg.Type)
	}
	return nil
}

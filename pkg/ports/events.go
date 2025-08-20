package ports

// EventPublisher defines the port for publishing domain events
// This interface abstracts event publishing operations from the core domain
type EventPublisher interface {
	// PublishAgentRegistered publishes an agent registration event
	PublishAgentRegistered(role string, agentType string) error

	// PublishAgentDeregistered publishes an agent deregistration event
	PublishAgentDeregistered(role string, reason string) error

	// PublishBarrelTransferred publishes a barrel transfer event
	PublishBarrelTransferred(fromRole string, toRole string, payload string) error

	// PublishSystemStatusChanged publishes a system status change event
	PublishSystemStatusChanged(status string, details map[string]interface{}) error
}

// DomainEvent represents a domain event that was published
type DomainEvent struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp string                 `json:"timestamp"`
}

package tcp

// TCPMessage represents the base structure for all TCP protocol messages
type TCPMessage struct {
	Type string `json:"type"`
}

// RegisterMessage represents agent registration requests
type RegisterMessage struct {
	Type         string   `json:"type"`         // "REGISTER"
	Role         string   `json:"role"`
	Capabilities []string `json:"capabilities"`
}

// YieldMessage represents yield requests from agents or people
type YieldMessage struct {
	Type     string `json:"type"` // "YIELD"
	FromRole string `json:"from_role"`
	ToRole   string `json:"to_role"`
	Payload  string `json:"payload"`
}

// QueryMessage represents query requests
type QueryMessage struct {
	Type string `json:"type"` // "QUERY_AGENTS" or "QUERY_STATUS"
}

// ActivateMessage represents activation messages sent to agents
type ActivateMessage struct {
	Type     string `json:"type"` // "ACTIVATE"
	FromRole string `json:"from_role"`
	Payload  string `json:"payload"`
}

// AgentListMessage represents response to agent list queries
type AgentListMessage struct {
	Type   string   `json:"type"` // "AGENT_LIST"
	Agents []string `json:"agents"`
}

// AgentDetailsMessage represents response to detailed agent queries
type AgentDetailsMessage struct {
	Type         string                   `json:"type"` // "AGENT_DETAILS"
	AgentDetails []AgentDetailInfo        `json:"agent_details"`
}

// AgentDetailInfo represents detailed information about a single agent
type AgentDetailInfo struct {
	Role         string   `json:"role"`
	Capabilities []string `json:"capabilities"`
	State        string   `json:"state"`
	Connected    bool     `json:"connected"`
}

// StatusMessage represents response to status queries
type StatusMessage struct {
	Type             string            `json:"type"` // "STATUS"
	BarrelHolder     string            `json:"barrel_holder"`
	RegisteredAgents []string          `json:"registered_agents"`
	AgentStates      map[string]string `json:"agent_states"`
	ConnectedAgents  map[string]bool   `json:"connected_agents"`
}

// ErrorMessage represents error responses
type ErrorMessage struct {
	Type    string `json:"type"` // "ERROR"
	Message string `json:"message"`
}

// AckRegisterMessage represents registration acknowledgment
type AckRegisterMessage struct {
	Type    string `json:"type"` // "ACK_REGISTER"
	Status  string `json:"status"`
	Message string `json:"message"`
}

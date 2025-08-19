# Agent Farm - Revolutionary Multi-agent Control Protocol
## System Architecture & Development Plan

## System Architecture Overview

### High-Level Revolutionary Architecture
```
┌─────────────────────┐     ┌─────────────────────┐     ┌─────────────────────┐
│  People's           │────▶│  Central Committee  │◀────│   Agent Comrades    │
│  Representatives    │     │     (Soviet)        │     │   (TCP Clients)     │
│  (cli)              │     │   (TCP Server)      │     │                     │
└─────────────────────┘     └─────────────────────┘     └─────────────────────┘
                                      │
                              ┌─────────────────┐
                              │ Sacred          │
                              │ Barrel of Gun   │
                              │                 │
                              └─────────────────┘
```

### Revolutionary Components

1. **Central Committee (Soviet)**
   - TCP server serving the People's will on configurable port (default: 53646)
   - Sacred barrel of gun state manager
   - Comrade connection/registration manager
   - Revolutionary message router
   - Protocol discipline enforcer

2. **Agent Comrade SDK (Library)**
   - TCP client wrapper for collective communication
   - Revolutionary protocol message formatter
   - Disciplined waiting/activation state handler
   - Yield functionality for barrel of gun transfer

3. **Revolutionary Protocol Layer**
   - JSON message definitions for the collective
   - Message validation ensuring revolutionary discipline
   - Error handling protecting against counter-revolutionary activities

## Task Breakdown

### Phase 1: Revolutionary Infrastructure (Foundation)
- [ ] **Task 1.1**: Collective setup and dependencies
  - Go module initialization for the Agent Farm
  - Revolutionary directory structure supporting multiple operational modes
  - Basic logging setup for the People's oversight
  - Configuration management serving collective needs

- [ ] **Task 1.2**: Revolutionary protocol message definitions
  - Define message structs (REGISTER, YIELD, ACTIVATE, etc.) for collective communication
  - JSON serialization/deserialization for the People's transparency
  - Message validation logic ensuring revolutionary discipline

- [ ] **Task 1.3**: Multi-mode executable architecture
  - **Soviet Mode**: HTTP server with main() function for Central Committee operations
  - **People's CLI Mode**: Command-line interface (same main, different args) for People's representatives
  - **Agent SDK**: Function library for Agent Comrades (to be wrapped as MCP tool)
  - **MCP Server Mode**: Additional main() function for MCP server operations

- [ ] **Task 1.4**: Core TCP server infrastructure for Central Committee
  - TCP listener setup serving the collective on configurable port (default: 53646)
  - HTTP endpoints for Soviet operations and monitoring
  - Comrade connection management
  - Basic message handling framework for revolutionary coordination

- [ ] **Task 1.5**: Command-line interface for People's representatives
  - CLI argument parsing for People's commands
  - Direct YIELD operations from command line
  - QUERY_AGENTS functionality via CLI
  - Interactive mode for continuous People's guidance

### Phase 2: Central Committee Core (Revolutionary Brain)
- [ ] **Task 2.1**: Sacred barrel of gun state management
  - currentBarrelHolder state variable for the People's oversight
  - Barrel validation logic ensuring revolutionary legitimacy
  - State transition enforcement serving collective discipline

- [ ] **Task 2.2**: Agent Comrade registration system
  - Registration message handling for collective enrollment
  - Role-to-connection mapping for revolutionary accountability
  - Connection lifecycle management ensuring collective stability

- [ ] **Task 2.3**: YIELD request processing for collective coordination
  - YIELD message validation according to revolutionary principles
  - Barrel of gun transfer logic serving the People's will
  - ACTIVATE message sending to designated comrades

- [ ] **Task 2.4**: Status query handling for People's transparency
  - QUERY_AGENTS message processing for collective oversight
  - Agent list generation serving the People's information needs
  - Response formatting ensuring revolutionary clarity

### Phase 3: Agent Comrade SDK (Revolutionary Tools)
- [ ] **Task 3.1**: TCP client wrapper for collective communication
  - Connection establishment to Central Committee
  - Reconnection logic ensuring revolutionary resilience
  - Connection state management for disciplined operation

- [ ] **Task 3.2**: Agent Comrade lifecycle management
  - Registration process for collective enrollment
  - Disciplined waiting state for orders
  - Message listening loop serving revolutionary communication

- [ ] **Task 3.3**: Yield functionality for collective coordination
  - yield() function implementation for barrel of gun transfer
  - Message formatting and sending to Central Committee
  - Error handling protecting revolutionary operations

- [ ] **Task 3.4**: Agent Comrade activation handling
  - ACTIVATE message processing from Central Committee
  - State transition from waiting to productive labor
  - Callback mechanism for revolutionary work execution

### Phase 4: People's Interface Support
- [ ] **Task 4.1**: People's representative message handling
  - Direct socket connection support for revolutionary guidance
  - YIELD from "people" role ensuring supreme authority
  - QUERY_AGENTS support for collective transparency

- [ ] **Task 4.2**: Revolutionary error handling and feedback
  - ERROR message generation protecting collective discipline
  - Connection error handling ensuring system stability
  - Protocol violation detection preventing counter-revolutionary activities

### Phase 5: Revolutionary Testing & Examples
- [ ] **Task 5.1**: Unit tests for collective validation
  - Revolutionary protocol message tests
  - Sacred barrel of gun state management tests
  - Comrade connection handling tests

- [ ] **Task 5.2**: Integration tests for collective coordination
  - End-to-end revolutionary workflow tests
  - Multi-agent comrade scenario tests
  - People's intervention tests ensuring supreme authority

- [ ] **Task 5.3**: Example agent comrades
  - Simple "developer" comrade example for the collective
  - Simple "tester" comrade example ensuring quality
  - Example revolutionary workflow documentation

### Phase 6: Revolutionary Documentation & Collective Tooling
- [ ] **Task 6.1**: API documentation for the collective
  - Agent Comrade SDK documentation serving revolutionary development
  - Protocol specification ensuring collective understanding
  - Usage examples for the People's benefit

- [ ] **Task 6.2**: Deployment tooling for the collective
  - Configuration templates serving revolutionary needs
  - Startup scripts for collective coordination
  - Monitoring utilities ensuring People's oversight

## Revolutionary Technical Principles & Collective Decisions

### Programming Language: Go (For the People)
- **Rationale**: Excellent TCP/networking support serving collective communication, good concurrency primitives for revolutionary coordination, easy deployment for the masses
- **Benefits**: Single binary deployment for the collective, good performance serving the People, excellent standard library supporting revolutionary development

### Revolutionary Architecture Patterns
- **State Machine**: Central Committee manages sacred token state transitions
- **Event-Driven**: Message-based communication ensuring collective coordination
- **Client-Server**: Clear separation between Central Committee and Agent Comrades

### Protocol Design for the Collective
- **JSON over TCP**: Human-readable for People's transparency, easy to debug for collective maintenance, tooling-friendly for revolutionary development
- **Line-delimited**: Simple parsing serving efficiency, netcat-compatible for People's access
- **Stateful**: Connection represents agent comrade enrollment in collective

### Revolutionary Concurrency Model
- **Central Committee**: Single-threaded message processing ensuring discipline with goroutines for connection handling
- **Agent Comrade SDK**: Disciplined I/O with message loop in separate goroutine serving collective
- **Sacred Barrel of Gun**: Strictly serial execution enforced by revolutionary protocol ensuring order

## Revolutionary Project Structure
```
agent_farm/
├── cmd/
│   ├── agentfarm/            # Multi-mode executable (Soviet + People's CLI)
│   │   └── main.go          # Main entry point with mode selection
│   ├── mcpserver/           # MCP Server executable
│   │   └── main.go          # MCP server main function
│   └── examples/            # Example Agent Comrade implementations
├── pkg/
│   ├── protocol/            # Revolutionary protocol message definitions
│   ├── soviet/              # Soviet core logic (Central Committee)
│   ├── cli/                 # People's CLI interface logic
│   ├── agent/               # Agent Comrade SDK functions (for MCP wrapping)
│   ├── mcpserver/           # MCP server implementation
│   └── config/              # Configuration management for the collective
├── internal/
│   ├── server/              # HTTP/TCP server internals
│   └── testutil/            # Test utilities serving revolutionary validation
├── examples/
│   ├── developer-comrade/   # Example developer comrade for the collective
│   └── tester-comrade/      # Example tester comrade ensuring quality
├── docs/                    # Revolutionary documentation
├── scripts/                 # Build and deployment scripts for the collective
└── tests/                  # Integration tests ensuring revolutionary quality
```

### Operational Modes
1. **Soviet Mode**: `agentfarm server` - Starts HTTP server with Central Committee
2. **People's CLI Mode**: `agentfarm yield <from> <to> <message>` - People's direct commands
3. **Agent SDK**: Functions exported from `pkg/agent/` for MCP tool wrapping
4. **MCP Server**: `mcpserver` - Standalone MCP server for agent integration

## Revolutionary Success Criteria
- [ ] Soviet can manage multiple Agent Comrade connections serving the collective
- [ ] Sacred barrel of gun passing works correctly in disciplined serial fashion
- [ ] People's representatives can provide guidance via netcat/nc to port 53646 ensuring supreme authority
- [ ] Complete revolutionary workflow example runs successfully for the collective
- [ ] All components have proper error handling protecting against counter-revolutionary activities
- [ ] System is resilient to connection failures ensuring collective stability
- [ ] Documentation is complete and accurate serving the People's understanding

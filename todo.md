# Agent Farm - Revolutionary Multi-agent Control Protocol
## Hexagonal Architecture Development Plan

## System Architecture Overview - Hexagonal Design

### Hexagonal Architecture Layout
```
                        ┌──────────────### Hexagonal Architecture Patterns
- **Domain-Driven Design**: Core business logic isolated in domain layer
- **Dependency Inversion**: Core depends on abstractions, not implementations  
- **Interface Segregation**: Small, focused interfaces for each concern
- **Single Responsibility**: Each layer has one clear purpose

### Revolutionary Domain Model Improvements
- **AgentComrade State Model**: Refined from 3 states (waiting/working/yielding) to 2 meaningful states (waiting/working)
- **Domain Operations**: Yield is now a proper domain method, not a persistent state
- **Semantic Correctness**: States represent what agents *are*, operations represent what they *do*
- **Test Coverage**: 100% test coverage with comprehensive TDD implementation──────────────────────┐
                        │             ADAPTERS                    │
                        │  ┌─────────────────┐ ┌─────────────────┐│
                        │  │   TCP Server    │ │   CLI Handler   ││
                        │  │   (Primary)     │ │   (Primary)     ││
                        │  └─────────────────┘ └─────────────────┘│
                        └─────────────┬───────────────────────────┘
                                      │
                        ┌─────────────▼───────────────────────────┐
                        │             PORTS                       │
                        │  ┌─────────────────┐ ┌─────────────────┐│
                        │  │  Soviet Service │ │  Agent Service  ││
                        │  │  (Interface)    │ │  (Interface)    ││
                        │  └─────────────────┘ └─────────────────┘│
                        └─────────────┬───────────────────────────┘
                                      │
                        ┌─────────────▼───────────────────────────┐
                        │              CORE                       │
                        │  ┌─────────────────┐ ┌─────────────────┐│
                        │  │ Barrel of Gun   │ │ Agent Registry  ││
                        │  │ Domain Model    │ │ Domain Model    ││
                        │  │                 │ │                 ││
                        │  └─────────────────┘ └─────────────────┘│
                        │  ┌─────────────────┐ ┌─────────────────┐│
                        │  │ Soviet          │ │ Revolutionary   ││
                        │  │ Coordinator     │ │ Protocol        ││
                        │  │ (Domain Logic)  │ │ (Domain Logic)  ││
                        │  └─────────────────┘ └─────────────────┘│
                        └─────────────────────────────────────────┘
```

### Hexagonal Components

**CORE DOMAIN (Center of Hexagon):**
- **Domain Models**: Barrel of Gun, Agent Comrade, Message entities
- **Domain Services**: Soviet Coordinator, Protocol Validator, Yield Manager
- **Business Logic**: All revolutionary rules and coordination logic

**PORTS (Interfaces):**
- **Primary Ports**: SovietService (command interface), AgentService (query interface)
- **Secondary Ports**: AgentRepository (persistence), EventPublisher (notifications)

**ADAPTERS:**
- **Primary Adapters**: TCP Server Handler, CLI Command Handler
- **Secondary Adapters**: In-Memory Repository, Console Logger

## Development Plan - Hexagonal Architecture Approach

**HEXAGONAL PRINCIPLE: Build from the inside out**
1. **CORE FIRST**: Implement domain models and business logic with zero external dependencies
2. **PORTS SECOND**: Define clean interfaces for external communication
3. **ADAPTERS LAST**: Implement TCP and CLI as interchangeable adapters

**IMPORTANT NOTICE: This project MUST follow Test-Driven Development (TDD) principles. Write unit tests FIRST before implementing any code. Follow the Red-Green-Refactor cycle for all development tasks.**

**REVOLUTIONARY DISCIPLINE: Complete ONE subtask at a time. After each subtask completion, present the work for comrade review before proceeding to the next subtask. No rushing through multiple tasks without proper collective oversight.**

### Phase 1: CORE DOMAIN IMPLEMENTATION (Pure Business Logic) 🎯 CURRENT FOCUS
**Goal: Implement the core domain without any external dependencies (no TCP, no JSON, no I/O)**

- [x] **Task 1.1**: Core Domain Models (Pure Go structs) ✅ **COMPLETED**
  - [x] BarrelOfGun entity with ownership tracking
  - [x] AgentComrade entity with role, state, and capabilities (refactored: removed yielding state, added Activate/Yield methods)
  - [x] RevolutionaryMessage entity for internal communication  
  - [x] SovietState entity managing the collective state
  - [x] Unit tests for all domain models (pure Go testing with Testify)

- [ ] **Task 1.2**: Core Business Logic - Soviet Coordinator Service
  - SovietCoordinator struct implementing all barrel management logic
  - Agent registration/deregistration logic (in-memory)
  - Barrel yield validation and transfer logic
  - Agent state transition logic (waiting -> working -> yielding)
  - Unit tests for all business rules (mocked dependencies)

- [ ] **Task 1.3**: Core Business Logic - Protocol Validation
  - ProtocolValidator service for revolutionary discipline enforcement
  - Message validation rules (independent of transport format)
  - Revolutionary rule enforcement (only barrel holder can yield)
  - State consistency validation logic
  - Unit tests for protocol validation rules

- [ ] **Task 1.4**: Core Business Logic - Yield Management
  - YieldManager service handling barrel transfers
  - Target validation (agent exists, can receive barrel)
  - Yield workflow orchestration (validate -> transfer -> notify)
  - Reconnection recovery logic for disconnected barrel holders
  - Unit tests for yield scenarios and edge cases

### Phase 2: PORTS DEFINITION (Clean Interfaces) 
**Goal: Define interfaces that the core needs to interact with the outside world**

- [ ] **Task 2.1**: Primary Ports (Driving the application)
  - SovietService interface (commands: RegisterAgent, ProcessYield, QueryStatus)
  - AgentService interface (queries: GetAgentState, GetBarrelStatus)
  - CommandHandler interface (abstraction for different input sources)

- [ ] **Task 2.2**: Secondary Ports (Driven by the application)  
  - AgentRepository interface (persistence abstraction)
  - MessageSender interface (communication abstraction)
  - EventPublisher interface (notification abstraction)
  - Logger interface (logging abstraction)

### Phase 3: TCP ADAPTER IMPLEMENTATION (Network Transport)
**Goal: Implement TCP as one possible adapter, easily replaceable**

- [ ] **Task 3.1**: TCP Primary Adapter (Server)
  - TCPSovietAdapter implementing CommandHandler interface
  - JSON message serialization/deserialization
  - TCP connection management and message routing
  - Integration with SovietService port

- [ ] **Task 3.2**: TCP Secondary Adapter (Client Communication)
  - TCPMessageSender implementing MessageSender interface
  - Connection lifecycle management
  - Error handling and reconnection logic
  - Message delivery confirmation

### Phase 4: CLI ADAPTER IMPLEMENTATION (Human Interface)
**Goal: Implement CLI as another adapter for People's representatives**

- [ ] **Task 4.1**: CLI Primary Adapter
  - CLICommandHandler implementing CommandHandler interface
  - Interactive command parsing and validation
  - Human-readable output formatting
  - Integration with SovietService port

- [ ] **Task 4.2**: CLI Secondary Adapter
  - ConsoleSender implementing MessageSender interface
  - Formatted output for human consumption
  - Status display and monitoring capabilities

### Phase 5: INFRASTRUCTURE ADAPTERS (Supporting Services)
**Goal: Implement infrastructure concerns as pluggable adapters**

- [ ] **Task 5.1**: Repository Adapters
  - InMemoryAgentRepository implementing AgentRepository interface
  - (Future: FileRepository, DatabaseRepository)

- [ ] **Task 5.2**: Logging Adapters  
  - ConsoleLogger implementing Logger interface
  - (Future: FileLogger, StructuredLogger)

### Phase 6: INTEGRATION & ASSEMBLY (Dependency Injection)
**Goal: Wire everything together with clean dependency injection**

- [ ] **Task 6.1**: Application Assembly
  - Dependency injection container/factory
  - Adapter configuration and wiring
  - Application startup coordination

- [ ] **Task 6.2**: Integration Testing
  - End-to-end tests with real adapters
  - TCP + Core integration tests
  - CLI + Core integration tests
  - Cross-adapter compatibility tests

## Directory Structure (Hexagonal Layout)

```
agent_farm/
├── cmd/                          # Application entry points
│   ├── soviet/                   # Central Committee executable
│   └── agent/                    # Agent Comrade executable
├── pkg/
│   ├── core/                     # DOMAIN CORE (no external dependencies)
│   │   ├── domain/               # Domain models and entities
│   │   │   ├── barrel.go         # BarrelOfGun entity
│   │   │   ├── agent.go          # AgentComrade entity  
│   │   │   ├── message.go        # RevolutionaryMessage entity
│   │   │   └── soviet.go         # SovietState entity
│   │   ├── services/             # Domain services (business logic)
│   │   │   ├── coordinator.go    # SovietCoordinator service
│   │   │   ├── validator.go      # ProtocolValidator service
│   │   │   └── yield_manager.go  # YieldManager service
│   │   └── errors/               # Domain-specific errors
│   ├── ports/                    # PORTS (interfaces)
│   │   ├── primary/              # Primary ports (driving)
│   │   │   ├── soviet_service.go # SovietService interface
│   │   │   └── agent_service.go  # AgentService interface  
│   │   └── secondary/            # Secondary ports (driven)
│   │       ├── repository.go     # AgentRepository interface
│   │       ├── sender.go         # MessageSender interface
│   │       └── logger.go         # Logger interface
│   └── adapters/                 # ADAPTERS (implementations)
│       ├── primary/              # Primary adapters
│       │   ├── tcp/              # TCP server adapter
│       │   └── cli/              # CLI command adapter
│       └── secondary/            # Secondary adapters
│           ├── memory/           # In-memory repository
│           ├── tcp/              # TCP client sender
│           └── console/          # Console logger
├── internal/                     # Internal utilities and helpers
└── test/                         # Integration tests
```

## Development Workflow (Red-Green-Refactor)

For each task:
1. **RED**: Write failing unit tests first
2. **GREEN**: Write minimal code to make tests pass  
3. **REFACTOR**: Clean up code while keeping tests green
4. **REVIEW**: Present work for collective approval
5. **NEXT**: Move to next task only after approval

## Benefits of This Hexagonal Approach

1. **Testability**: Core domain logic can be tested in isolation
2. **Flexibility**: TCP can be replaced with HTTP, gRPC, or any other transport
3. **Independence**: Business logic doesn't depend on frameworks or infrastructure
4. **Maintainability**: Clear separation of concerns and responsibilities
5. **Revolutionary Discipline**: Each layer has a single, well-defined purpose

## Revolutionary Technical Principles & Collective Decisions

### Programming Language: Go (For the People)
- **Rationale**: Excellent TCP/networking support, good concurrency primitives, easy deployment
- **Benefits**: Single binary deployment, good performance, excellent standard library

### Testing Framework: Testify (For Revolutionary Quality)
- **Test Assertions**: Use `github.com/stretchr/testify/assert` for all test assertions
- **Test Suites**: Use `github.com/stretchr/testify/suite` for complex test scenarios
- **Mocking**: Use `github.com/stretchr/testify/mock` for dependency mocking in unit tests
- **Benefits**: Readable test code, comprehensive assertion library, excellent mocking capabilities

### Hexagonal Architecture Patterns
- **Domain-Driven Design**: Core business logic isolated in domain layer
- **Dependency Inversion**: Core depends on abstractions, not implementations  
- **Interface Segregation**: Small, focused interfaces for each concern
- **Single Responsibility**: Each layer has one clear purpose

**NEXT STEP**: Complete Task 1.1 by implementing SovietState entity, then proceed to Task 1.2 - Core Business Logic (Soviet Coordinator Service). We continue building the revolutionary system from the heart outward, following strict hexagonal principles!

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

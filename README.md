# Agent Farm - MCP (Multi-agent Control Protocol)
Version: 4.0  
Date: August 19, 2025

## 1. Project Overview

Agent Farm is a revolutionary protocol and architecture for building and coordinating multi-agent collaboration systems based on MCP (Multi-agent Control Protocol). In this collective farm of digital workers, agents operate as disciplined comrades executing tasks in a strictly serial manner through a virtual credential called the "Barrel of Gun" that manages work permission flow according to the People's will.

The core of the system is a centralized coordinator (Soviet) that serves as the Central Committee managing the barrel of gun on behalf of the People. After completing their assigned labor, agent comrades must "yield" the barrel of gun to the next designated agent or return it to the People, who hold supreme authority. Human users, representing the People's will, interact with the system through direct Socket connections to enable revolutionary guidance and status oversight.

The protocol is implemented over TCP Sockets and provides a dedicated SDK for agent comrades, while defining clear raw protocol interaction methods for the People's representatives.

## 2. Core Concepts

The following core concepts are essential for understanding this system:

### The Barrel of Gun
This is the sacred credential of labor, managed by the Central Committee (Soviet) on behalf of the People. Possessing the barrel of gun is the only authorization for an agent comrade to transition from waiting state to productive labor. There is only one barrel of gun in the entire collective, ensuring disciplined coordination.

### Serialized Workflow
The fundamental principle of our collective system. Since there is only one barrel of gun, all agent workflows are completely serial at the macro level, eliminating capitalist chaos of concurrent execution and ensuring task flow determinism and revolutionary accountability.

### Yielding
This is the sacred act of transferring the barrel of gun in service to the collective. When an active agent comrade completes its current labor assignment, it must initiate a YIELD request to the Central Committee, explicitly designating the next recipient of the barrel (another comrade or returning it to the People) and may attach a progress report. After completing the yield, the agent immediately returns the barrel and enters a disciplined waiting state for further orders.

### People Reserved Role
"people" is the supreme authority in our system, representing the collective will of the working class. This role does not correspond to any specific agent process but embodies the democratic decision-making power. When the barrel of gun is yielded to "people", the entire automated process pauses in respectful deference, all agent comrades wait in disciplined formation, until the People's representatives provide new revolutionary guidance through direct Socket connections.

## 3. System Architecture

This system adopts a revolutionary collective model consisting of three main components:

### Central Committee (Soviet)
As the administrative apparatus of the People, it is a long-running process that serves the collective. It is the keeper of the barrel of gun and enforcer of revolutionary discipline, responsible for:

- Listening to TCP ports and managing all Agent Comrade network connections
- Maintaining the record of which comrade currently holds the barrel of gun
- Strictly enforcing collective workflows, processing and validating all YIELD requests according to revolutionary principles
- Responding to status queries from the People's representatives

### Agent Comrades
As the working class of the system, each agent is an independent, long-running process representing a specialized worker with specific skills and revolutionary consciousness (such as developer, tester). Each Agent Comrade's lifecycle operates in a disciplined "waiting -> laboring -> yielding" cycle in service to the collective.

### MCP Tool (Agent SDK)
This tool provides the ideological framework (Agent SDK) - a library embedded in Agent Comrade programs. Agent Comrades serve the collective by calling functions provided by this library (such as yield()), which handles the revolutionary protocols and maintains proper communication with the Central Committee.

## 4. Component Details

### 4.1 Central Committee (Soviet)
**Connection Management**: Accepts and maintains TCP long connections from all Agent Comrades, while also accepting direct connections from the People's representatives for revolutionary guidance.

**Registration**: After agent comrades connect, they must first declare their service to the collective through registration messages. The Central Committee maintains a real-time roster from role names to their network connections, ensuring accountability.

**Agent Reconnection & Recovery**: When an agent comrade reconnects after disconnection, the Central Committee must check if the sacred barrel of gun belongs to this returning comrade. If so, the Central Committee immediately sends an ACTIVATE message to resume the revolutionary workflow, preventing system-wide blocking caused by offline agents.

**Barrel of Gun Management & Revolutionary Discipline**:
- Internally maintains the sacred record currentBarrelHolder (string)
- When receiving YIELD requests, must verify that the request comes from the current legitimate barrel holder, preventing counter-revolutionary activities
- Handles disconnected barrel holders by detecting reconnections and resuming workflow automatically

**Yield Processing**:
- After validating YIELD requests according to revolutionary principles, updates currentBarrelHolder to serve the collective
- Sends "ACTIVATE" orders and attached directives to the designated comrade through their connection

**Status Query Response**:
- When receiving inquiries from the People's representatives, the Central Committee traverses its internal roster, extracts all registered comrade roles, and reports the state of the collective

### 4.2 Agent Comrades
Agent Comrade processes embody revolutionary discipline through their state cycle:

1. **Startup & Registration**: Process awakens, connects to Central Committee through Agent SDK and declares service to the collective
2. **Waiting (Disciplined Readiness)**: After successful registration, enters respectful waiting state for orders from the Central Committee
3. **Revolution in Progress (Working)**: After receiving "activation" orders, begins productive work in service to the collective
4. **Yielding (Revolutionary Transfer)**: After completing labor assignment, calls yield(target_role, progress_report) function provided by Agent SDK to serve the next phase of collective work
5. **Cycle Continuation**: After sending YIELD message, Agent Comrade immediately returns to disciplined waiting state for further revolutionary tasks

### 4.3 People's Interface (Direct Revolutionary Communication)
The People's representatives do not use predetermined tools of the old regime, but connect directly to the Central Committee's TCP port through revolutionary communication tools (such as netcat, telnet, or custom scripts) and provide guidance through protocol messages that embody the People's will.

**Connection Method**: Use netcat (or nc) to establish direct communication with the Central Committee, for example:
```bash
nc localhost 53646
```

**Revolutionary Operations**:

- **Issue Labor Directives**: The People's representatives provide guidance by inputting complete, single-line JSON directives
  ```json
  {"type":"YIELD","from_role":"people","to_role":"developer","payload":"Comrade Developer, implement the authentication module for the People's digital security."}
  ```

- **Inspect Collective Status**: The People's representatives can inquire about the state of their digital collective
  ```json
  {"type":"QUERY_AGENTS"}
  ```
  The Central Committee will report the AGENT_LIST message on this connection, after which representatives may disconnect at their discretion

## 5. Communication Protocol

Defines revolutionary communication protocols between Agent Comrades and the Central Committee.

### Agent Comrades -> Central Committee Messages

**REGISTER**
- User: Agent Comrade
- Format: `{"type": "REGISTER", "role": "developer"}`

**YIELD**
- User: Agent Comrade, People's Representatives
- Format: `{"type": "YIELD", "from_role": "developer", "to_role": "tester", "payload": "Comrade Tester, the code is ready for revolutionary quality assurance."}`

**QUERY_AGENTS**
- User: People's Representatives
- Format: `{"type": "QUERY_AGENTS"}`

### Central Committee -> Agent Comrades Messages

**ACTIVATE**
- Receiver: Agent Comrade
- Format: `{"type": "ACTIVATE", "from_role": "developer", "payload": "Comrade Tester, the code is ready for revolutionary quality assurance."}`

**AGENT_LIST**
- Receiver: People's Representatives
- Format: `{"type": "AGENT_LIST", "agents": ["developer", "tester", "code-reviewer"]}`

**ERROR**
- Receiver: Agent Comrade, People's Representatives
- Format: `{"type": "ERROR", "message": "Yield request denied: counter-revolutionary action detected - you are not the current barrel holder."}`

**ACK_REGISTER**
- Receiver: Agent Comrade
- Format: `{"type": "ACK_REGISTER", "status": "success", "message": "Comrade 'developer' successfully enlisted in the collective."}`

## 6. Revolutionary Workflow Example

1. **Collective Awakening**: Central Committee process starts. currentBarrelHolder initially held by "people" - the supreme authority

2. **Comrades Report for Duty**: developer and tester Agent Comrades start their processes, connect and register respectively. Both enter disciplined waiting state

3. **People's Directive**: A representative of the People opens terminal, uses netcat to establish revolutionary communication with Central Committee:
   ```bash
   nc localhost 53646
   ```
   After successful connection, they provide guidance through the following directive:
   ```json
   {"type":"YIELD","from_role":"people","to_role":"developer","payload":"Comrade Developer, implement feature #123 for the advancement of our digital collective."}
   ```
   Then may disconnect (Ctrl+C)

4. **First Labor Assignment**: Central Committee receives the People's directive, grants barrel to developer comrade, and sends ACTIVATE message

5. **Developer Labors**: developer Agent Comrade begins productive work. tester comrade remains in disciplined waiting

6. **Revolutionary Transfer**: After completion, developer yields through: `yield("tester", "Comrade Tester, feature #123 has been implemented and awaits your revolutionary quality assurance.")`

7. **Barrel Flows**: Central Committee grants barrel to tester comrade. developer returns to disciplined waiting state

8. **Seeking People's Wisdom**: tester Agent Comrade discovers critical issues, yields to the People: `yield("people", "Comrade Representatives, I have discovered a critical design flaw that requires the People's guidance.")`

9. **Revolutionary Pause**: Central Committee updates currentBarrelHolder to "people". All Agent Comrades await the People's wisdom in disciplined formation

10. **People's Inspection**: A representative reconnects to assess the collective state:
    ```json
    {"type":"QUERY_AGENTS"}
    ```
    Central Committee immediately reports:
    ```json
    {"type":"AGENT_LIST","agents":["developer","tester"]}
    ```

11. **People's Decision**: After reviewing the collective status, the People's representative issues new guidance:
    ```json
    {"type":"YIELD","from_role":"people","to_role":"developer","payload":"Comrade Developer, apply your revolutionary consciousness to fix the design flaw discovered by our vigilant tester comrade."}
    ```

12. **Collective Continues**: developer Agent Comrade receives new orders, revolutionary work continues in service to the People

## 7. Edge Case: Agent Disconnection & Recovery

A critical scenario that must be handled to prevent system deadlock:

**Problem**: If an agent goes offline after another agent has yielded the barrel of gun to it, the entire system becomes blocked waiting for the disconnected agent.

**Solution**: When a disconnected agent reconnects and registers with the Central Committee:

1. **Reconnection Detection**: Central Committee receives REGISTER message from returning agent
2. **Barrel Check**: Central Committee checks if currentBarrelHolder matches the reconnecting agent's role
3. **Automatic Recovery**: If the barrel belongs to the reconnecting agent, Central Committee immediately sends ACTIVATE message with the last payload
4. **Workflow Resume**: The agent resumes work, and the revolutionary workflow continues

**Example Recovery Scenario**:
- developer yields to tester: `yield("tester", "Code ready for testing")`
- tester goes offline before receiving ACTIVATE
- System is blocked, developer waits in disciplined formation
- tester reconnects, sends: `{"type": "REGISTER", "role": "tester"}`
- Central Committee detects currentBarrelHolder == "tester", immediately sends: `{"type": "ACTIVATE", "from_role": "developer", "payload": "Code ready for testing"}`
- Revolutionary workflow resumes without People's intervention

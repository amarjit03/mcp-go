# Architecture & Design

## High-Level Architecture

```
┌─────────────────────────────────────────────────────┐
│                  AI Models / Clients                 │
│         (Claude, GPT, Custom Applications)           │
└────────────────────┬────────────────────────────────┘
                     │ JSON-RPC 2.0
                     │ (TCP Port 9090)
                     ▼
┌─────────────────────────────────────────────────────┐
│              MCP Server (Go)                         │
├─────────────────────────────────────────────────────┤
│  pkg/mcp/                                           │
│  ├── types.go       - Protocol types                │
│  └── protocol.go    - JSON-RPC handler              │
├─────────────────────────────────────────────────────┤
│  internal/server/                                   │
│  └── server.go      - Connection & request routing  │
├─────────────────────────────────────────────────────┤
│  internal/tools/                                    │
│  └── manager.go     - Tool implementations          │
├─────────────────────────────────────────────────────┤
│  internal/system/                                   │
│  ├── executor.go    - Command & file execution      │
│  └── monitor.go     - System monitoring             │
└─────────────────────────────────────────────────────┘
                     │
         ┌───────────┴────────────┬──────────────┐
         │                        │              │
         ▼                        ▼              ▼
    ┌─────────┐            ┌──────────┐    ┌─────────┐
    │ Shell   │            │ Files    │    │ System  │
    │Commands │            │ & Logs   │    │ Stats   │
    └─────────┘            └──────────┘    └─────────┘
```

## Component Breakdown

### 1. Protocol Layer (`pkg/mcp/`)

**Responsibility**: Handle MCP protocol communication

**Files**:
- `types.go` - JSON-RPC message types, tool definitions
- `protocol.go` - Request/response handling, routing

**Key Types**:
- `Request`, `Response`, `Error` - JSON-RPC messages
- `Tool`, `InputSchema` - Tool definitions
- `Protocol` - Request handler registry

### 2. Server Layer (`internal/server/`)

**Responsibility**: Connection management and request dispatching

**Files**:
- `server.go` - TCP listener, connection handler, protocol handlers

**Key Types**:
- `Server` - Main server instance
- `Config` - Server configuration

**Handlers**:
- `handleInitialize` - Protocol initialization
- `handleListTools` - Tool enumeration
- `handleCallTool` - Tool execution

### 3. Tools Layer (`internal/tools/`)

**Responsibility**: Tool implementations and management

**Files**:
- `manager.go` - Tool registry and implementations

**Available Tools**:
1. **File Operations**
   - `read_file` - Read file contents
   - `write_file` - Write to files
   - `list_directory` - List directory contents
   - `read_logs` - Read log files

2. **System Commands**
   - `execute_command` - Run shell commands

3. **System Health**
   - `get_cpu_usage` - CPU information
   - `get_memory_usage` - Memory information
   - `check_port` - Port availability
   - `get_process_info` - Process details
   - `health_check` - Full system check

### 4. System Layer (`internal/system/`)

**Responsibility**: Safe system operations with security controls

**Files**:
- `executor.go` - Command and file execution with safety checks
- `monitor.go` - System monitoring via gopsutil

**Components**:

#### CommandExecutor
- Executes shell commands safely
- Applies timeouts (default: 30s)
- Enforces command allow/deny lists
- Limits output size (default: 10MB)

#### FileReader
- Reads files safely
- Enforces path restrictions
- Limits file size (default: 50MB)
- Blocks denied paths

#### SystemMonitor
- Gets CPU usage and counts
- Gets memory usage and percentage
- Checks port availability
- Retrieves process information
- Performs comprehensive health checks

## Data Flow

### Initialize Request Flow

```
Client sends "initialize"
    │
    ▼
Server.handleConnection()
    │
    ▼
Protocol.HandleRequest()
    │
    ▼
handleInitialize()
    │
    ├─ Create InitializeResponse
    │  │
    │  ├─ Server info
    │  ├─ Capabilities
    │  └─ Tool list
    │
    ▼
Response sent to client
```

### Tool Execution Flow

```
Client sends "tools/call"
    │
    ▼
Server.handleConnection()
    │
    ▼
Protocol.HandleRequest()
    │
    ▼
handleCallTool()
    │
    ▼
ToolManager.ExecuteTool()
    │
    ├─ Validate tool name
    │
    ├─ Execute tool function
    │  │
    │  ├─ Validate arguments
    │  │
    │  ├─ Apply security checks
    │  │
    │  ├─ Call system components
    │  │  (Executor, FileReader, Monitor)
    │  │
    │  └─ Return result
    │
    ▼
Response sent to client
```

## Security Architecture

### Multi-Layer Protection

```
┌──────────────────────────────────────┐
│   Client Request (untrusted)         │
└──────────────────┬───────────────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │ 1. Request Validation │
        │   - Schema check      │
        │   - Type check        │
        └──────────┬───────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │ 2. Tool Lookup       │
        │   - Tool exists?      │
        └──────────┬───────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │ 3. Argument Check    │
        │   - Type validation  │
        │   - Range check      │
        └──────────┬───────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │ 4. Security Gate     │
        │ - File: Path restrict│
        │ - Cmd: Allow/deny    │
        │ - Exec: Timeout      │
        └──────────┬───────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │ 5. Execute           │
        │   - Isolated context │
        │   - Output limited   │
        └──────────┬───────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │ 6. Response          │
        │   - Sanitized        │
        │   - Structured       │
        └──────────┬───────────┘
                   │
                   ▼
┌──────────────────────────────────────┐
│   Client Response (safe)             │
└──────────────────────────────────────┘
```

## Configuration Hierarchy

```
1. Default Config (hardcoded)
   ├─ Port: 9090
   ├─ Host: 127.0.0.1
   ├─ Timeout: 30s
   └─ Max sizes: 50MB/10MB

2. Config File (config/default.yaml)
   └─ Overrides defaults

3. Command-Line Flags
   └─ Overrides file config
```

## Extension Points

### Adding New Tools

1. **Implement function** in `ToolManager`:
   ```go
   func (tm *ToolManager) myTool(args map[string]interface{}) (interface{}, error) {
       // Implementation
       return result, nil
   }
   ```

2. **Register in `registerBuiltinTools()`**:
   ```go
   tm.RegisterTool("my_tool", tm.myTool)
   ```

3. **Add to `GetTools()`** documentation

### Custom Security Policies

Extend `ExecutorConfig` or `FileReaderConfig`:

```go
config := &system.ExecutorConfig{
    AllowedCommands: []string{"ls", "grep", "ps"},
    DeniedCommands: []string{"rm", "reboot"},
    Timeout: 60 * time.Second,
}
executor := system.NewCommandExecutor(config)
```

## Performance Considerations

### Connection Handling
- Concurrent goroutines per connection
- Graceful shutdown of idle connections
- Configurable read timeouts

### System Calls
- Efficient process listing via gopsutil
- Network checks with TCP dial timeout
- Memory-efficient file reading (chunked for large files)

### Caching Opportunities
- CPU/memory stats (cache for 1-5 seconds)
- Process info (live queries only)
- Port states (cache briefly)

## Testing Strategy

### Unit Tests
- Tool manager functions
- Protocol message handling
- Executor safety checks

### Integration Tests
- Full request/response cycle
- Tool execution end-to-end
- Security enforcement

### Performance Tests
- Large file reading
- Command output limits
- Concurrent connections

## Deployment Considerations

### Minimal Setup
```bash
docker run -p 9090:9090 mcp-go:latest
```

### Restricted Access
```bash
./mcp-server -host 127.0.0.1  # Localhost only
```

### Custom Config
```bash
./mcp-server -config /etc/mcp/config.yaml
```

### Monitoring
- Log all tool executions
- Track execution times
- Monitor resource usage
- Alert on errors

# Development Guide

## Getting Started

### Prerequisites
- Go 1.21 or later
- Unix-like system (Linux, macOS) or WSL on Windows
- Make (optional but recommended)

### Setup

```bash
# Clone repository
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go

# Install dependencies
go mod download

# Build server
go build -o bin/mcp-server ./cmd/server

# Or using make
make build
```

## Project Structure

```
mcp-go/
├── cmd/server/           # Entry point
│   └── main.go          # Server startup
├── internal/
│   ├── server/          # MCP server implementation
│   │   └── server.go    # Connection & protocol handling
│   ├── tools/           # Tool implementations
│   │   └── manager.go   # Tool registry & execution
│   └── system/          # System operations
│       ├── executor.go  # Safe command/file execution
│       └── monitor.go   # System monitoring
├── pkg/mcp/             # MCP protocol
│   ├── types.go         # Protocol types
│   └── protocol.go      # JSON-RPC handler
├── config/              # Configuration templates
├── examples/            # Example clients
├── docs/                # Documentation
├── go.mod               # Module definition
├── Makefile             # Build commands
└── README.md            # User documentation
```

## Key Files to Know

### Protocol Implementation
- **`pkg/mcp/types.go`** - Message types, tool definitions
- **`pkg/mcp/protocol.go`** - Request routing and response encoding

### Server Core
- **`internal/server/server.go`** - Main server logic
  - `Start()` - Listener and accept loop
  - `handleConnection()` - Per-connection handler
  - Protocol handlers (initialize, listTools, callTool)

### Tools
- **`internal/tools/manager.go`** - All tool implementations
  - File tools: read_file, write_file, list_directory, read_logs
  - Command tools: execute_command
  - System tools: get_cpu_usage, get_memory_usage, check_port, get_process_info, health_check

### System Operations
- **`internal/system/executor.go`**
  - `CommandExecutor` - Safe command execution with timeouts, output limiting
  - `FileReader` - Safe file access with path restrictions

- **`internal/system/monitor.go`**
  - `SystemMonitor` - CPU, memory, process, and port monitoring

## Development Workflow

### Running the Server

```bash
# Debug mode with default config
make run

# With custom port
./bin/mcp-server -port 8080

# With config file
./bin/mcp-server -config config/default.yaml

# With all options
./bin/mcp-server -host 0.0.0.0 -port 8080 -config config/default.yaml
```

### Code Formatting

```bash
# Format all code
make fmt

# Or manually
go fmt ./...
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run specific test
go test -v ./internal/tools -run TestToolManager
```

### Running Examples

```bash
# Python example
python3 examples/python_client.py

# JavaScript example
node examples/javascript_client.js

# Make targets
make run-py-example
make run-js-example
```

## Adding New Tools

### Step 1: Implement Function

Add to `internal/tools/manager.go`:

```go
func (tm *ToolManager) myNewTool(args map[string]interface{}) (interface{}, error) {
    // Extract arguments
    param1, ok := args["param1"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid param1")
    }
    
    // Implement tool logic
    result := "success"
    
    // Return structured result
    return map[string]interface{}{
        "result": result,
        "param1": param1,
    }, nil
}
```

### Step 2: Register Tool

In `registerBuiltinTools()`:

```go
tm.RegisterTool("my_new_tool", tm.myNewTool)
```

### Step 3: Add to Tool List

In `GetTools()`:

```go
{
    Name:        "my_new_tool",
    Description: "Description of what the tool does",
    InputSchema: mcp.InputSchema{
        Type: "object",
        Properties: map[string]interface{}{
            "param1": map[string]interface{}{
                "type":        "string",
                "description": "Parameter description",
            },
        },
        Required: []string{"param1"},
    },
}
```

## Extending Security

### Custom File Restrictions

In `cmd/server/main.go` or `internal/server/server.go`:

```go
fileReaderConfig := &system.FileReaderConfig{
    MaxFileSize:  100 * 1024 * 1024, // 100MB
    AllowedPaths: []string{
        "/home/user/projects",
        "/var/log",
    },
    DeniedPaths: []string{
        "/etc/shadow",
        "/root",
    },
}
fileReader := system.NewFileReader(fileReaderConfig)
```

### Custom Command Restrictions

```go
cmdExecConfig := &system.ExecutorConfig{
    Timeout:         60 * time.Second,
    MaxOutputSize:   20 * 1024 * 1024, // 20MB
    AllowedCommands: []string{
        "ls", "cat", "grep", "ps", "netstat",
    },
    DeniedCommands: []string{
        "rm", "reboot", "shutdown",
    },
}
cmdExecutor := system.NewCommandExecutor(cmdExecConfig)
```

## Debugging

### Enable Logging

The server logs to stderr with timestamps. Check logs for:
- Connection events
- Tool execution
- Errors and exceptions

### Debug Connection

Use `nc` or telnet to test:

```bash
# Connect to server
nc localhost 9090

# Send JSON request
{"jsonrpc":"2.0","method":"tools/list","id":1}

# Press Ctrl+D to close connection
```

### Use Client Examples

Python client with debugging:

```python
import json

response = client.execute_command("echo test")
print(json.dumps(response, indent=2))
```

## Performance Profiling

### CPU Profiling

```bash
go build -o bin/mcp-server ./cmd/server

# Run with profiling
pprof -http=:6060 bin/mcp-server
```

### Memory Profiling

```bash
# Add to main.go for production profiles
import _ "net/http/pprof"

// Then profile with
go tool pprof http://localhost:6060/debug/pprof/heap
```

## Testing

### Write Unit Tests

Example test in `internal/tools/manager_test.go`:

```go
package tools

import (
    "testing"
)

func TestReadFile(t *testing.T) {
    tm := NewToolManager(...)
    
    result, err := tm.readFile(map[string]interface{}{
        "path": "test.txt",
    })
    
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    
    if result == nil {
        t.Fatal("Expected result, got nil")
    }
}
```

### Run Tests with Coverage

```bash
make coverage
# Opens coverage.html in browser
```

## Building for Production

### Static Build

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -o bin/mcp-server \
  -ldflags="-w -s" \
  ./cmd/server
```

### Build for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/mcp-server-linux ./cmd/server

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/mcp-server-macos ./cmd/server

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/mcp-server.exe ./cmd/server
```

### Docker Build

```bash
docker build -t mcp-go:latest .
docker run -p 9090:9090 mcp-go:latest
```

## Common Issues

### Port Already in Use

```bash
# Find process using port
lsof -i :9090

# Kill process
kill -9 <PID>

# Or use different port
./bin/mcp-server -port 9091
```

### Permission Denied

```bash
# Make binary executable
chmod +x bin/mcp-server

# Run with appropriate permissions
./bin/mcp-server
```

### Module Not Found

```bash
# Ensure go.mod is present
go mod init github.com/amarjit-singh/mcp-go

# Download modules
go mod download

# Tidy modules
go mod tidy
```

## Contributing

1. Fork repository
2. Create feature branch: `git checkout -b feature/name`
3. Make changes and format: `make fmt`
4. Add tests
5. Verify tests pass: `make test`
6. Commit and push
7. Create pull request

## Resources

- [Go Documentation](https://golang.org/doc/)
- [JSON-RPC 2.0 Spec](https://www.jsonrpc.org/specification)
- [MCP Protocol](https://modelcontextprotocol.io/)
- [gopsutil Library](https://github.com/shirou/gopsutil)

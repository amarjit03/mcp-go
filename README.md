# MCP Dev Assistant - Go Implementation

A lightweight **Model Context Protocol (MCP)** server in Go that acts as a bridge between AI models and real-world system operations. This server enables AI to safely interact with a developer's environment through structured tools, reducing context switching and making development faster and smarter.

## Features

✨ **Safe System Access**
- File reading/writing with path restrictions
- Command execution with allowlisting/denylisting
- Process monitoring and port checking
- System health checks (CPU, memory)
- Log file retrieval

🔐 **Security First**
- Configurable file access restrictions
- Command execution safety guards
- Timeout protection for long-running commands
- Output size limiting
- Structured error handling

📡 **JSON-RPC 2.0 Protocol**
- Standard MCP protocol implementation
- Easy integration with AI models and clients
- Request/response handling with proper error codes

🎯 **Real-World Use Cases**
- "Is my backend running?" → Port check + process info
- "Show recent errors" → Log file reading
- "Check system health" → CPU, memory, and port status
- "Run this command" → Safe command execution
- "What's in this file?" → File reading with restrictions

## Architecture

```
mcp-go/
├── cmd/server/              # Main entry point
├── internal/
│   ├── server/             # MCP server implementation
│   ├── tools/              # Tool implementations
│   └── system/             # System utilities (executor, monitor)
├── pkg/mcp/                # MCP protocol types and handlers
├── config/                 # Configuration templates
├── examples/               # Example client implementations
└── README.md
```

## Quick Start

### Prerequisites
- Go 1.21 or later

### Installation

```bash
# Clone repository
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go

# Install dependencies
go mod download

# Build server
go build -o bin/mcp-server ./cmd/server
```

### Running the Server

**Default configuration (localhost:9090):**
```bash
./bin/mcp-server
```

**With custom port:**
```bash
./bin/mcp-server -port 8080 -host 0.0.0.0
```

**With configuration file:**
```bash
./bin/mcp-server -config config/default.yaml
```

## How to Use It (Step-by-Step Guide)

### 🖥️ macOS/Linux Setup

#### Step 1: Install Ollama (for the AI Chatbot)

**macOS:**
```bash
brew install ollama
```

**Linux:**
```bash
curl -fsSL https://ollama.ai/install.sh | sh
```

#### Step 2: Download a Text Generation Model

```bash
# Recommended (fast, lightweight, good quality)
ollama pull tinyllama

# Or try these alternatives:
ollama pull mistral       # Faster, better quality
ollama pull neural-chat   # Optimized for conversations
ollama pull llama2        # More capable
```

#### Step 3: Clone and Build MCP Server

```bash
# Clone the repository
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go

# Build the server
go build -o bin/mcp-server ./cmd/server

# Build the chatbot
go build -o bin/chatbot ./cmd/chatbot
```

#### Step 4: Start Services (3 Terminal Windows)

**Terminal 1 - Start Ollama:**
```bash
ollama serve
# Output: Listening on 127.0.0.1:11434
```

**Terminal 2 - Start MCP Server:**
```bash
cd /path/to/mcp-go
./bin/mcp-server -port 9090
# Output: [MCP Server] Listening on localhost:9090
```

**Terminal 3 - Start Chatbot:**
```bash
cd /path/to/mcp-go
./bin/chatbot -model tinyllama
# Chatbot ready for input!
```

#### Step 5: Try Example Queries

```
You: What is my CPU usage?
🤔 Thinking... 
Bot: Your CPU is currently using 23% of available processing power...

You: Is port 8080 open?
🤔 Thinking... 
Bot: Port 8080 is currently closed on your system...

You: Check system health
🤔 Thinking... 
Bot: Your system is running well with good resource availability...

You: exit
👋 Goodbye!
```

---

### 🪟 Windows Setup

#### Step 1: Install Ollama for Windows

1. Visit [ollama.ai](https://ollama.ai)
2. Download the Windows installer
3. Run the installer and follow the prompts
4. Open Command Prompt or PowerShell

#### Step 2: Download a Text Generation Model

```powershell
# Download model (choose one)
ollama pull tinyllama

# Or alternatives:
ollama pull mistral
ollama pull neural-chat
```

#### Step 3: Install Go (if not already installed)

1. Visit [golang.org](https://golang.org/dl)
2. Download and run the Windows installer
3. Verify installation:
```powershell
go version
```

#### Step 4: Clone and Build MCP Server

```powershell
# Open Command Prompt or PowerShell

# Clone repository
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go

# Build the server
go build -o bin\mcp-server.exe .\cmd\server

# Build the chatbot
go build -o bin\chatbot.exe .\cmd\chatbot
```

#### Step 5: Start Services (3 Command Prompts)

**Command Prompt 1 - Start Ollama:**
```powershell
ollama serve
# Output: Listening on 127.0.0.1:11434
```

**Command Prompt 2 - Start MCP Server:**
```powershell
cd C:\path\to\mcp-go
.\bin\mcp-server.exe -port 9090
# Output: [MCP Server] Listening on localhost:9090
```

**Command Prompt 3 - Start Chatbot:**
```powershell
cd C:\path\to\mcp-go
.\bin\chatbot.exe -model tinyllama
# Chatbot ready!
```

#### Step 6: Try Example Queries

Same as macOS/Linux examples above.

---

### 📚 Alternative: Using MCP Server with Python Client

If you don't want to use the chatbot, you can communicate with the MCP server directly:

```python
import socket
import json

# Connect to MCP server
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('localhost', 9090))

# Send initialize request
request = {
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {
        "protocolVersion": "2024-11-05",
        "clientInfo": {"name": "python-client", "version": "1.0"}
    },
    "id": 1
}

sock.send(json.dumps(request).encode())
response = sock.recv(4096).decode()
print(response)

sock.close()
```

See `examples/python_client.py` for a complete example.

---

### 🛠️ Troubleshooting

| Issue | Solution |
|-------|----------|
| "Ollama not available" | Make sure `ollama serve` is running in Terminal 1 |
| "Cannot connect to MCP server" | Make sure `./bin/mcp-server` is running in Terminal 2 |
| "Model not found" | Run `ollama pull tinyllama` first |
| "Port 9090 already in use" | Use different port: `./bin/mcp-server -port 8090` |
| "Permission denied" on Linux | Run `chmod +x bin/mcp-server bin/chatbot` |

---

### 📖 Configuration

Edit `config/default.yaml` to customize server behavior:

```yaml
port: 9090
host: 127.0.0.1
name: mcp-dev-assistant
version: 1.0.0

# Execution settings
commandTimeout: 30s
fileReadMaxSize: 52428800  # 50MB
commandMaxOutput: 10485760 # 10MB

# Security: Restrict to specific paths
allowedPaths:
  - /home/user/projects
  - /var/log

# Security: Allow only specific commands
allowedCommands:
  - ls
  - cat
  - grep
  - ps

# Always denied commands (for safety)
deniedCommands:
  - rm
  - rmdir
  - reboot
  - shutdown
```

## MCP Protocol

The server implements the **Model Context Protocol** using JSON-RPC 2.0.

### Message Format

All communication uses JSON-RPC 2.0:

```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "clientInfo": {
      "name": "my-ai-client",
      "version": "1.0.0"
    }
  },
  "id": 1
}
```

### Available Methods

#### 1. **initialize** - Initialize connection
```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "clientInfo": {
      "name": "client-name",
      "version": "1.0.0"
    }
  },
  "id": 1
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": true
    },
    "serverInfo": {
      "name": "mcp-dev-assistant",
      "version": "1.0.0"
    },
    "tools": [...]
  },
  "id": 1
}
```

#### 2. **tools/list** - List available tools
```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 2
}
```

#### 3. **tools/call** - Execute a tool

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "execute_command",
    "arguments": {
      "command": "ps aux | grep python"
    }
  },
  "id": 3
}
```

## Available Tools

### File Operations

#### **read_file** - Read file contents
```json
{
  "name": "read_file",
  "arguments": {
    "path": "/path/to/file.txt"
  }
}
```

#### **write_file** - Write to file
```json
{
  "name": "write_file",
  "arguments": {
    "path": "/path/to/file.txt",
    "content": "file content",
    "append": false
  }
}
```

#### **list_directory** - List directory contents
```json
{
  "name": "list_directory",
  "arguments": {
    "path": "/path/to/dir"
  }
}
```

#### **read_logs** - Read last N lines of log file
```json
{
  "name": "read_logs",
  "arguments": {
    "path": "/var/log/app.log",
    "lines": 50
  }
}
```

### System Commands

#### **execute_command** - Run shell command
```json
{
  "name": "execute_command",
  "arguments": {
    "command": "ls -la /home/user"
  }
}
```

### System Health

#### **get_cpu_usage** - Get CPU information
```json
{
  "name": "get_cpu_usage",
  "arguments": {}
}
```

**Response:**
```json
{
  "percent": 45.2,
  "count": 4,
  "countLogical": 8,
  "timestamp": "2024-03-31T10:30:00Z"
}
```

#### **get_memory_usage** - Get memory information
```json
{
  "name": "get_memory_usage",
  "arguments": {}
}
```

**Response:**
```json
{
  "total": 8589934592,
  "used": 4294967296,
  "free": 4294967296,
  "available": 5000000000,
  "usedPercent": 50.0,
  "timestamp": "2024-03-31T10:30:00Z"
}
```

#### **check_port** - Check if port is open
```json
{
  "name": "check_port",
  "arguments": {
    "port": 8080,
    "protocol": "tcp"
  }
}
```

**Response:**
```json
{
  "port": 8080,
  "protocol": "tcp",
  "state": "open",
  "timestamp": "2024-03-31T10:30:00Z"
}
```

#### **get_process_info** - Get process information
```json
{
  "name": "get_process_info",
  "arguments": {
    "pid_or_name": "1234"
  }
}
```

#### **health_check** - Perform system health check
```json
{
  "name": "health_check",
  "arguments": {
    "ports": [3000, 5432, 8080]
  }
}
```

**Response:**
```json
{
  "cpu": {...},
  "memory": {...},
  "ports": [...],
  "status": "healthy",
  "timestamp": "2024-03-31T10:30:00Z"
}
```

## Example: Python Client

```python
import socket
import json

def send_request(method, params=None):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(("127.0.0.1", 9090))
    
    request = {
        "jsonrpc": "2.0",
        "method": method,
        "params": params or {},
        "id": 1
    }
    
    sock.sendall(json.dumps(request).encode() + b'\n')
    
    response = b''
    while True:
        chunk = sock.recv(4096)
        if not chunk:
            break
        response += chunk
    
    sock.close()
    return json.loads(response.decode())

# Initialize
response = send_request("initialize", {
    "protocolVersion": "2024-11-05",
    "clientInfo": {
        "name": "python-client",
        "version": "1.0"
    }
})
print("Initialized:", response)

# Check port
response = send_request("tools/call", {
    "name": "check_port",
    "arguments": {"port": 8080}
})
print("Port status:", response)

# Get CPU usage
response = send_request("tools/call", {
    "name": "get_cpu_usage",
    "arguments": {}
})
print("CPU usage:", response)
```

## Example: JavaScript Client

```javascript
const net = require('net');

function sendRequest(method, params = {}) {
  return new Promise((resolve, reject) => {
    const socket = net.createConnection({ port: 9090, host: '127.0.0.1' });
    
    const request = {
      jsonrpc: '2.0',
      method: method,
      params: params,
      id: 1
    };
    
    socket.write(JSON.stringify(request) + '\n');
    
    let data = '';
    socket.on('data', (chunk) => {
      data += chunk.toString();
    });
    
    socket.on('end', () => {
      try {
        resolve(JSON.parse(data));
      } catch (e) {
        reject(e);
      }
    });
    
    socket.on('error', reject);
  });
}

// Example usage
(async () => {
  const response = await sendRequest('tools/call', {
    name: 'execute_command',
    arguments: { command: 'whoami' }
  });
  console.log('Command output:', response);
})();
```

## Security Considerations

### 1. **Path Restrictions**
- Use `allowedPaths` to restrict file access to specific directories
- Default: All paths allowed (except explicitly denied)
- Recommended: Whitelist only necessary paths

### 2. **Command Safety**
- Built-in denial of dangerous commands (rm, reboot, etc.)
- Use `allowedCommands` to implement strict command allowlisting
- Use `deniedCommands` to add additional restrictions

### 3. **Timeouts**
- Commands have a default 30-second timeout
- Configurable via `commandTimeout`
- Prevents hanging or resource exhaustion

### 4. **Output Limiting**
- Default 50MB limit for file reads
- Default 10MB limit for command output
- Prevents memory overflow from large outputs

### 5. **Access Control**
- Run server with minimal privileges
- Use localhost (127.0.0.1) by default
- Restrict network access via firewall
- Consider authentication if exposing over network

## Development & Extension

### Adding a New Tool

1. Implement a function in `internal/tools/manager.go`:
```go
func (tm *ToolManager) myNewTool(args map[string]interface{}) (interface{}, error) {
  // Implementation
  return result, nil
}
```

2. Register in `registerBuiltinTools()`:
```go
tm.RegisterTool("my_new_tool", tm.myNewTool)
```

3. Add to `GetTools()` documentation

### Custom Configuration

Create a custom config file and pass via `-config`:
```bash
./bin/mcp-server -config /path/to/custom.yaml
```

## Performance Tips

- **Use `allowedPaths`** to reduce filesystem access checks
- **Set reasonable `commandTimeout`** based on your workload
- **Limit `fileReadMaxSize`** to prevent memory issues
- **Monitor port checks** - TCP connect is slower than UDP
- **Consider caching** for frequently-requested information

## Troubleshooting

### Server won't start
- Check if port is already in use: `lsof -i :9090`
- Verify firewall settings
- Check config file syntax with `yamllint`

### Tool execution fails
- Check permission errors in logs
- Verify command is not in denied list
- Confirm path is in allowed list (if configured)
- Check timeout settings

### Slow performance
- Reduce file read size limit
- Optimize command execution
- Check system resources (CPU, memory)

## Contributing

Contributions are welcome! Areas for enhancement:
- Additional tools (git, docker, kubernetes)
- Authentication and authorization
- Persistent logging
- Metrics and monitoring
- WebSocket support
- gRPC integration

## License

MIT License - See LICENSE file

## Resources

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [Go Standard Library](https://pkg.go.dev/std)

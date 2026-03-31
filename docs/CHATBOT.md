# MCP System Chatbot

A lightweight AI-powered chatbot that bridges local language models (via Ollama) with the MCP server to enable natural language queries about system status.

## Overview

The chatbot acts as an intelligent interface between you and system tools. Simply ask questions in natural language, and the chatbot will:

1. **Understand** your intent (CPU usage, port checking, health status, etc.)
2. **Query** the appropriate MCP server tools
3. **Synthesize** real system data into a human-readable response using a lightweight LLM

## Features

- **Ultra-lightweight**: Uses TinyLlama by default (~1.1GB RAM)
- **Privacy-first**: All models run locally (no cloud API calls)
- **Zero-latency**: Direct connection to MCP server
- **Tool-aware**: Automatically selects relevant tools based on query
- **Conversation history**: Maintains context during session

## Prerequisites

### 1. Ollama (LLM Runtime)

Install from [ollama.ai](https://ollama.ai)

```bash
# macOS
brew install ollama

# Linux - download from ollama.ai
curl -fsSL https://ollama.ai/install.sh | sh
```

### 2. TinyLlama Model (Ultra-lightweight)

Download the default model:

```bash
ollama pull tinyllama
```

Alternative lightweight models:
```bash
ollama pull llama2            # ~4GB, more accurate
ollama pull neural-chat       # ~5GB, optimized for chat
```

### 3. MCP Server

Build the MCP server (if not already done):

```bash
cd /path/to/mcp-go
go build -o bin/mcp-server ./cmd/server
```

## Running the Chatbot

### 1. Start Ollama (Terminal 1)

```bash
ollama serve
# Output: Listening on 127.0.0.1:11434
```

### 2. Start MCP Server (Terminal 2)

```bash
./bin/mcp-server -port 9090 -host localhost
# Output: [MCP Server] Listening on localhost:9090
```

### 3. Start Chatbot (Terminal 3)

```bash
./bin/chatbot
# Or with custom addresses:
./bin/chatbot -ollama http://127.0.0.1:11434 -mcp localhost:9090
```

### Example Session

```
🤖 MCP System Chatbot - Powered by Ollama + Lightweight LLM
============================================================

Initializing chatbot...
  - Ollama URL: http://127.0.0.1:11434
  - MCP Server: localhost:9090

[ChatBot] Using model: tinyllama
✅ Chatbot ready!

💡 Example queries:
  - What is my CPU usage?
  - Is my backend running on port 8080?
  - Check system health
  - Show recent logs
  - What processes are running?

You: What is my CPU usage?

🤔 Thinking... 

Bot: Your CPU is currently running at 23% utilization across 8 cores. This 
indicates moderate system load with plenty of available processing capacity.

You: Is port 8080 open?

🤔 Thinking... 

Bot: Port 8080 is currently CLOSED on your system. No service is listening 
on that port at the moment.

You: exit

👋 Goodbye!
```

## Supported Query Types

The chatbot recognizes keywords to determine which tools to invoke:

| Keyword | Tools Invoked | Example |
|---------|-------------|---------|
| cpu, processor | `get_cpu_usage` | "What is my CPU usage?" |
| memory, ram | `get_memory_usage` | "How much RAM is in use?" |
| port, running, backend | `check_port`, `health_check` | "Is port 8080 open?" |
| process | `get_process_info` | "Is nginx running?" |
| health, status, system | `health_check` | "Check system health" |
| log, error | `read_logs` | "Show recent errors" |
| file | `read_file` | "Read /etc/hostname" |

Default (if no keywords match): `health_check` (full system snapshot)

## How It Works

### Architecture

```
User Input
    ↓
Keyword Analysis (determineTools)
    ↓
MCP Server Queries (executeMCPTools)
    ↓
System Data Collection
    ↓
LLM Prompt Building
    ↓
Ollama Generation
    ↓
Response Display
```

### MCP Tool Integration

The chatbot communicates with the MCP server via TCP/JSON-RPC:

1. **Parse query** for keywords
2. **Build arguments** for relevant tools
3. **Send JSON-RPC request** to MCP server (port 9090)
4. **Collect responses** from multiple tools
5. **Format data** for LLM consumption
6. **Send to Ollama** with system prompt
7. **Return synthesized response** to user

### Example: "Is my backend running on port 8080?"

```
Query: "Is my backend running on port 8080?"
↓
Keywords Found: "port", "backend", "running"
↓
Tools Selected: check_port, health_check
↓
MCP Calls:
  - check_port(port=8080, protocol="tcp")
  - health_check(ports=[22,80,443,3000,5432,8080,9090])
↓
System Data:
  [check_port]
  Port 8080: CLOSED
  
  [health_check]
  CPU: 15.2%
  Memory: 45.3%
  Ports: 22(OPEN), 80(OPEN), 443(OPEN), 8080(CLOSED)
↓
LLM Prompt:
  System: You are a helpful Linux system assistant...
  User Query: Is my backend running on port 8080?
  System Data: [collected above]
  Your Response (2-3 sentences):
↓
LLM Response (TinyLlama):
  "Port 8080 is currently closed on your system, which means no backend 
   service is listening there. You may need to check if your backend 
   application is running or has been properly configured to listen on 
   that port."
```

## Configuration

### Command-Line Flags

```bash
./bin/chatbot [flags]

Flags:
  -ollama string     Ollama API URL (default "http://127.0.0.1:11434")
  -mcp string        MCP server address host:port (default "localhost:9090")
```

### Environment Variables (Optional)

Can be extended to support:
```bash
export OLLAMA_URL=http://custom.ollama:11434
export MCP_ADDR=remotehost:9090
./bin/chatbot
```

## Troubleshooting

### "Ollama not available"

```bash
# Check if Ollama is running
curl http://127.0.0.1:11434/api/tags

# If not, start Ollama
ollama serve

# If still failing, check port
lsof -i :11434
```

### "Cannot connect to MCP server"

```bash
# Check if MCP server is running
lsof -i :9090

# If not, start it
./bin/mcp-server -port 9090

# If connection refused, check firewall
# For custom MCP port:
./bin/chatbot -mcp localhost:9091
```

### "No models available"

```bash
# Pull a model
ollama pull tinyllama

# List available models
ollama list
```

### LLM Response is Slow

- TinyLlama is optimized for speed but still needs ~2-5 seconds per response
- For faster responses on CPU-only systems, reduce model complexity
- Consider Llama2 if you have 8+ GB RAM for better accuracy

### Custom Models

To use a different model:

1. Pull the model: `ollama pull llama2`
2. The chatbot automatically detects and uses available models
3. First model in list is used by default

To force a specific model, edit [internal/chatbot/chatbot.go](../internal/chatbot/chatbot.go#L45):

```go
// Change this line:
model := models[0]

// To:
model := "llama2"  // Your preferred model
```

## Performance Characteristics

| Model | Size | Memory | Speed | Accuracy |
|-------|------|--------|-------|----------|
| TinyLlama | 1.1GB | 1-2GB | 2-3s/response | Good |
| Llama2 | 4GB | 6-8GB | 5-10s/response | Excellent |
| Neural-Chat | 5GB | 8-10GB | 3-5s/response | Very Good |

## Advanced Usage

### Conversation History

The chatbot maintains conversation history internally. Each message is tracked:

```go
cb.conversationHistory  // Slice of Message{Role, Content}
```

This could be extended to:
- Persist conversations to disk
- Provide multi-turn context for follow-up questions
- Generate conversation logs

### Custom Tool Arguments

Modify [internal/chatbot/chatbot.go#L130](../internal/chatbot/chatbot.go#L130) `buildToolArguments()` to customize default tool parameters:

```go
case "health_check":
    args["ports"] = []int{22, 80, 443, 3000, 5432, 8080, 9090}
    // Add custom ports here
```

### System Prompt Tuning

Edit [internal/chatbot/chatbot.go#L76](../internal/chatbot/chatbot.go#L76) `buildSystemPrompt()` to adjust LLM behavior:

```go
func (cb *Chatbot) buildSystemPrompt() string {
    return `You are a helpful Linux system assistant...`
    // Customize tone, style, or instructions
}
```

## Architecture Details

### File Structure

```
internal/chatbot/
├── chatbot.go          # Core chatbot logic (410 lines)
│   ├── Chatbot struct
│   ├── NewChatbot()    # Factory function
│   ├── Chat()          # Main conversation loop
│   ├── determineTools()  # Keyword → tool mapping
│   ├── executeMCPTools() # Call MCP server
│   └── formatResult()  # Parse MCP responses

internal/llm/
├── ollama.go           # Ollama HTTP client (175 lines)
│   ├── OllamaClient struct
│   ├── Generate()      # Send prompt, get response
│   ├── IsAvailable()   # Check Ollama health
│   └── GetModels()     # List available models

cmd/chatbot/
├── main.go             # CLI entry point (80 lines)
│   ├── Flag parsing
│   ├── Error handling
│   ├── REPL loop
│   └── Exit handling
```

### Data Flow Diagram

```
┌──────────┐
│   User   │
│  Input   │
└────┬─────┘
     │
     ↓
┌─────────────────────────────┐
│  Chatbot.Chat()             │
│  - Parse keywords           │
│  - Select tools             │
│  - Call MCP server          │
│  - Collect data             │
└────┬────────────────────────┘
     │
     ├─→ TCP/JSON-RPC → MCP Server (port 9090)
     │                      ↓
     │            ┌──────────────────┐
     │            │ System Tools:    │
     │            │ - CPU Usage      │
     │            │ - Memory Usage   │
     │            │ - Port Status    │
     │            │ - Process Info   │
     │            │ - Health Check   │
     │            └────────┬─────────┘
     │                     ↓
     │         Returns JSON-RPC response
     │                     ↑
     │                     │
     ├──────────────────────┘
     │
     ↓
┌────────────────────────────┐
│ Build LLM Prompt:          │
│ - System message           │
│ - User query               │
│ - System data              │
└────┬───────────────────────┘
     │
     ↓
┌──────────────────────────────┐
│ HTTP → Ollama (port 11434)   │
│ POST /api/generate           │
└────┬───────────────────────┬─┘
     │                       │
     ↓                       ↓
  Input                   Model Inference
  "What is my           (TinyLlama generates
   CPU usage?"          natural response)
                             ↓
┌───────────────────────────────────────┐
│ Response:                             │
│ "Your CPU is at 23%..."               │
└───────────────────────────────────────┘
     ↑
     │
┌────┴──────────────────────────────────┐
│  Format & Display to User              │
│  Update conversation history           │
│  Prompt for next input                 │
└────────────────────────────────────────┘
```

## Performance Tips

1. **Pre-load model**: Models are cached in memory after first use
2. **Keep queries focused**: Shorter queries are processed faster
3. **Batch queries**: Multiple related questions in one query
4. **Monitor resources**: Watch CPU/RAM during LLM inference

## Future Enhancements

- [ ] Context-aware follow-up questions
- [ ] Multi-language support
- [ ] Custom knowledge base injection
- [ ] Tool output caching
- [ ] Batch query execution
- [ ] Web UI interface
- [ ] Export conversation logs
- [ ] Real-time streaming responses

## Debugging

Enable detailed logging by modifying [cmd/chatbot/main.go](../cmd/chatbot/main.go):

```go
fmt.Printf("[DEBUG] Determined tools: %v\n", tools)
fmt.Printf("[DEBUG] Tool results: %s\n", toolResults)
```

## License

Same as parent MCP project

## Related Documentation

- [MCP Server Architecture](./ARCHITECTURE.md)
- [Protocol Specification](./PROTOCOL.md)
- [System Tools Reference](./TOOLS.md) (if created)

# MCP Go - Complete Setup Guide

This guide covers everything you need to know to set up and use MCP Go.

## Table of Contents

1. [What is MCP Go?](#what-is-mcp-go)
2. [System Requirements](#system-requirements)
3. [Getting Started](#getting-started)
4. [Using the CLI](#using-the-cli)
5. [Manual Setup](#manual-setup)
6. [Next Steps](#next-steps)

## What is MCP Go?

MCP Go is a **Model Context Protocol (MCP)** server that connects AI models to your system. It provides:

- **Real System Access**: Get actual CPU usage, running processes, file contents
- **AI Chatbot**: Ask questions in natural language, get intelligent answers
- **Safe Execution**: Run commands safely with permission checks
- **JSON-RPC Interface**: Easy integration with any AI application

### Architecture

```
┌─────────────────┐
│   AI Model      │
│  (e.g., Claude) │
└────────┬────────┘
         │
    JSON-RPC 2.0
         │
    ┌────▼────────────────────┐
    │   MCP Go Server         │
    │  (Port 9090)            │
    ├─────────────────────────┤
    │  • Tool Registry        │
    │  • Security Framework   │
    │  • Command Executor     │
    │  • System Monitor       │
    └────┬───────────────────┬┘
         │                   │
    ┌────▼───────┐    ┌──────▼─────┐
    │ System Ops │    │ Groq API   │
    │ (Files,    │    │ (LLM)      │
    │  Processes)│    │            │
    └────────────┘    └────────────┘
```

## System Requirements

- **OS**: macOS, Linux, or Windows
- **Go**: 1.21 or later
- **Network**: Internet access for Groq API
- **Groq Account**: Free account at https://console.groq.com

### Check Go Version

```bash
go version
# Output: go version go1.22 linux/amd64
```

## Getting Started

### Option 1: CLI Setup (Recommended - 5 minutes)

The CLI makes setup incredibly easy with an interactive wizard.

#### Step 1: Clone & Build

```bash
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go

# Build the CLI tool
make build-cli
```

#### Step 2: Interactive Setup

```bash
./bin/mcp-cli init
```

Follow the prompts:
```
╔════════════════════════════════════════════╗
║  🚀 MCP Go - Setup Wizard                  ║
║  AI-Powered System Assistant               ║
╚════════════════════════════════════════════╝

🔍 Checking prerequisites...

📌 Setting up Groq API Key
==========================

Get your API key from: https://console.groq.com/keys

Enter your Groq API key (gsk_...): [paste your key]

🔨 Building binaries...

✅ Setup Complete!
```

#### Step 3: Start Using

```bash
# Start both services
./bin/mcp-cli start

# Or in separate terminals:
# Terminal 1:
./bin/mcp-cli server

# Terminal 2:
./bin/mcp-cli chat

# Check status
./bin/mcp-cli status
```

### Option 2: Manual Setup (10 minutes)

For more control over the setup process.

#### Step 1: Clone Repository

```bash
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go
```

#### Step 2: Get Groq API Key

1. Visit https://console.groq.com/keys
2. Create an account or sign in
3. Click "Create API Key"
4. Copy the key (starts with `gsk_`)

#### Step 3: Create .env File

```bash
cat > .env << EOF
GROQ_API_KEY=gsk_your_actual_key_here
EOF
```

Replace `gsk_your_actual_key_here` with your actual key.

#### Step 4: Build Components

```bash
# Build everything at once
make build-all

# Or build individually:
go build -o bin/mcp-server ./cmd/server
go build -o bin/chatbot ./cmd/chatbot
go build -o bin/mcp-cli ./cmd/cli
```

#### Step 5: Start Services

**In Terminal 1:**
```bash
./bin/mcp-server
# Output: MCP Server listening on :9090
```

**In Terminal 2:**
```bash
export GROQ_API_KEY=$(grep GROQ_API_KEY .env | cut -d'=' -f2)
./bin/chatbot
# Output: Welcome to MCP Go Chatbot!
```

## Using the CLI

Once set up, use the CLI to manage services:

```bash
# View help
./bin/mcp-cli help

# Start everything
./bin/mcp-cli start

# Check what's running
./bin/mcp-cli status

# See configuration
./bin/mcp-cli config

# Stop services
./bin/mcp-cli stop
```

### CLI Commands Reference

| Command | Purpose | Example |
|---------|---------|---------|
| `init` | Setup wizard | `./bin/mcp-cli init` |
| `start` | Start all services | `./bin/mcp-cli start` |
| `server` | Start MCP server | `./bin/mcp-cli server` |
| `chat` | Start chatbot | `./bin/mcp-cli chat` |
| `stop` | Stop services | `./bin/mcp-cli stop` |
| `status` | Check status | `./bin/mcp-cli status` |
| `config` | View config | `./bin/mcp-cli config` |
| `test-api` | Test API | `./bin/mcp-cli test-api` |
| `help` | Get help | `./bin/mcp-cli help` |

## Manual Setup

### Building from Source

```bash
# Download dependencies
go mod download
go mod tidy

# Build server
go build -o bin/mcp-server ./cmd/server

# Build chatbot
go build -o bin/chatbot ./cmd/chatbot

# Build CLI
go build -o bin/mcp-cli ./cmd/cli
```

### Running Without CLI

If you prefer to run services directly:

```bash
# Terminal 1: Start server
./bin/mcp-server

# Terminal 2: Start chatbot
export GROQ_API_KEY=$(grep GROQ_API_KEY .env | cut -d'=' -f2)
./bin/chatbot
```

### Configuration Options

#### Environment Variables

```bash
# Required
export GROQ_API_KEY=gsk_your_key_here

# Optional
export MCP_PORT=9090        # Default: 9090
export MCP_HOST=localhost   # Default: localhost
```

#### Configuration File

Create `config/default.yaml`:

```yaml
server:
  port: 9090
  host: localhost
  
llm:
  provider: groq
  model: llama-3.3-70b-versatile
  temperature: 0.3
  max_tokens: 500
  
security:
  allowed_paths:
    - /var/log
    - /tmp
  allowed_commands:
    - ps
    - df
    - top
  blocked_commands:
    - rm -rf /
    - sudo
```

## File Locations

After setup, files are organized as:

```
mcp-go/
├── .env                     # Configuration (API key)
├── bin/
│   ├── mcp-cli             # CLI tool
│   ├── mcp-server          # MCP server binary
│   ├── chatbot             # Chatbot binary
│   ├── .server.pid         # Server process ID
│   └── .chatbot.pid        # Chatbot process ID
├── logs/                   # Log files (created if needed)
└── config/                 # Configuration files
```

## Verification

### Check Everything is Working

```bash
# 1. Check status
./bin/mcp-cli status

# Expected output:
#   🟢 MCP Server: RUNNING
#   🟢 Chatbot: RUNNING
#   ✓ GROQ_API_KEY: Set

# 2. Test API key
./bin/mcp-cli test-api

# Expected output:
#   ✅ API key format is valid

# 3. View configuration
./bin/mcp-cli config

# Expected output shows your settings
```

### Test With Chatbot

Once running, try these in the chatbot:

```
You: What's my CPU usage?
Assistant: Your CPU usage is currently at 45.2%

You: List running processes
Assistant: Here are the running processes...

You: Check system health
Assistant: Your system is running well with...

You: help
Assistant: Available commands: info, processes, logs...
```

## Troubleshooting

### "GROQ_API_KEY not set"

```bash
# Re-run setup
./bin/mcp-cli init

# Or manually set it
export GROQ_API_KEY=gsk_your_key_here
```

### "Connection refused" (Port 9090)

The server isn't running. Start it:
```bash
./bin/mcp-cli server

# Or manually
./bin/mcp-server
```

### "Permission denied"

Make binaries executable:
```bash
chmod +x bin/mcp-cli bin/mcp-server bin/chatbot
```

### "Module not found" when building

Download dependencies:
```bash
go mod download
go mod tidy
go build -o bin/mcp-cli ./cmd/cli
```

## Next Steps

### 1. Learn About MCP Protocol
See [MCP Architecture Guide](./docs/MCP_ARCHITECTURE.md)

### 2. Explore Available Tools
See [System Tools Reference](./docs/SYSTEM_TOOLS.md)

### 3. Understand Security
See [Security Framework](./docs/SECURITY.md)

### 4. Integrate With Your Application
See [Integration Guide](./docs/INTEGRATION.md)

### 5. Read Full CLI Documentation
See [CLI Complete Guide](./docs/CLI.md)

## Common Workflows

### Workflow 1: Quick Diagnostics

```bash
# Start services
./bin/mcp-cli start

# In chatbot:
You: What's the system status?
You: Check network connections
You: Show memory usage
```

### Workflow 2: Development Testing

```bash
# Terminal 1: Server
./bin/mcp-cli server

# Terminal 2: Chatbot
./bin/mcp-cli chat

# Terminal 3: Run tests
make test
```

### Workflow 3: Running in Production

```bash
# Install globally
sudo cp bin/mcp-cli /usr/local/bin/mcp-go

# Start services
mcp-go start

# Monitor
mcp-go status
```

## Performance Metrics

- **CLI Startup**: < 100ms
- **Server Startup**: < 500ms
- **Chatbot Startup**: ~1s
- **API Response**: 1-3s (Groq API latency)
- **Binary Sizes**:
  - CLI: 2.5 MB
  - Server: 4.8 MB
  - Chatbot: 7.3 MB

## Getting Help

1. **Check CLI help**: `./bin/mcp-cli help`
2. **Read documentation**: See `docs/` directory
3. **View logs**: Check CLI output for errors
4. **GitHub Issues**: Open an issue on GitHub

## Security Considerations

- ✅ API key stored locally in `.env` (git-ignored)
- ✅ No data sent except to Groq API
- ✅ Commands executed locally with restrictions
- ✅ File access limited to specified paths
- ✅ Timeout protection on all operations

## Next Commands

After successful setup:

```bash
# View available tools
./bin/mcp-cli config

# Check documentation
less docs/CLI.md

# Read chatbot guide
less docs/CHATBOT.md

# Explore system tools
less docs/SYSTEM_TOOLS.md
```

---

**Ready to go?** Run `./bin/mcp-cli init` now! 🚀

For questions or issues, see [Troubleshooting](#troubleshooting) or open a GitHub issue.

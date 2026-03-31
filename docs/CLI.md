# MCP Go CLI - User Guide

The MCP Go CLI (`mcp-cli`) is a professional command-line utility that makes it easy to set up, manage, and use the MCP Go system.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Configuration](#configuration)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go

# Build the CLI
make build-cli

# The binary will be available at: ./bin/mcp-cli
```

### Add to PATH (Optional)

To use `mcp-cli` from anywhere:

```bash
# Copy to a directory in PATH (macOS/Linux)
sudo cp bin/mcp-cli /usr/local/bin/mcp-go

# Or create a symbolic link
sudo ln -s $(pwd)/bin/mcp-cli /usr/local/bin/mcp-go

# Verify installation
mcp-go help
```

### Build All Components

To build the MCP server, chatbot, and CLI together:

```bash
make build-all
```

## Quick Start

### 1. Initial Setup

Run the interactive setup wizard:

```bash
./bin/mcp-cli init
```

This will:
- Check your system prerequisites
- Create necessary directories
- Prompt for your Groq API key
- Build the server and chatbot binaries
- Save configuration to `.env` file

### 2. Start Services

Start both MCP server and chatbot:

```bash
./bin/mcp-cli start
```

### 3. Check Status

Verify services are running:

```bash
./bin/mcp-cli status
```

### 4. Stop Services

Stop all running services:

```bash
./bin/mcp-cli stop
```

## Commands

### `mcp-go init`

**Interactive Setup Wizard**

Runs a step-by-step configuration process. Perfect for first-time setup.

```bash
./bin/mcp-cli init
```

**What it does:**
- ✓ Checks system prerequisites (Go, environment)
- ✓ Creates necessary directories (bin, logs, config)
- ✓ Prompts for Groq API key
- ✓ Builds server and chatbot binaries
- ✓ Saves configuration to `.env`

**Output:**
```
╔════════════════════════════════════════════╗
║  🚀 MCP Go - Setup Wizard                  ║
║  AI-Powered System Assistant               ║
╚════════════════════════════════════════════╝

🔍 Checking prerequisites...

📌 Setting up Groq API Key
==========================

Get your API key from: https://console.groq.com/keys

Enter your Groq API key (gsk_...): [user input]

🔨 Building binaries...
   Building MCP server...
   ✓ Server built
   Building chatbot...
   ✓ Chatbot built

╔════════════════════════════════════════════╗
║  ✅ Setup Complete!                        ║
╚════════════════════════════════════════════╝
```

---

### `mcp-go start`

**Start All Services**

Launches both the MCP server and chatbot in the background, with the chatbot in interactive mode.

```bash
./bin/mcp-cli start
```

**What it does:**
- Checks if binaries exist (builds if needed)
- Loads configuration from `.env`
- Starts MCP server on port 9090
- Starts chatbot (interactive)
- Saves process PIDs for management

**Output:**
```
🚀 Starting MCP Go services...

📌 Starting MCP server...
   ✓ Server started (PID: 12345)
💬 Starting chatbot...
   ✓ Chatbot started (PID: 12346)

✅ All services started! Type 'help' in chatbot for commands.
```

---

### `mcp-go server`

**Start MCP Server Only**

Launches just the MCP server without the chatbot. Useful when you want to keep the server running in a separate terminal.

```bash
./bin/mcp-cli server
```

**What it does:**
- Checks/builds server binary if needed
- Loads GROQ_API_KEY from `.env`
- Starts MCP server on port 9090
- Runs in foreground (press Ctrl+C to stop)

**Output:**
```
📌 Starting MCP Server...

MCP Server listening on :9090
```

---

### `mcp-go chat`

**Start Chatbot Only**

Launches just the chatbot. Requires MCP server to be running on localhost:9090.

```bash
./bin/mcp-cli chat
```

**What it does:**
- Checks if MCP server is running
- Checks/builds chatbot binary if needed
- Loads GROQ_API_KEY from `.env`
- Starts interactive chatbot

**Output:**
```
💬 Starting Chatbot...

Welcome to MCP Go Chatbot!
Type 'help' for available commands

You: [waiting for input]
```

---

### `mcp-go stop`

**Stop All Services**

Gracefully shuts down both MCP server and chatbot.

```bash
./bin/mcp-cli stop
```

**What it does:**
- Reads saved process PIDs
- Sends SIGTERM (signal 15) to processes
- Cleans up PID files
- Falls back to `pkill` if PID files not found

**Output:**
```
🛑 Stopping services...

   ✓ Stopped MCP server
   ✓ Stopped chatbot

✅ All services stopped
```

---

### `mcp-go status`

**Show Service Status**

Displays the current status of both services and configuration.

```bash
./bin/mcp-cli status
```

**Output Examples:**

When services are running:
```
╔════════════════════════════════════════════╗
║  📊 MCP Go Status                          ║
╚════════════════════════════════════════════╝

  🟢 MCP Server: RUNNING (port 9090)
  🟢 Chatbot: RUNNING

  Configuration:
    ✓ GROQ_API_KEY: Set
```

When services are stopped:
```
╔════════════════════════════════════════════╗
║  📊 MCP Go Status                          ║
╚════════════════════════════════════════════╝

  🔴 MCP Server: STOPPED
  🔴 Chatbot: STOPPED

  Configuration:
    ✓ GROQ_API_KEY: Set
```

---

### `mcp-go config`

**View Configuration**

Displays the current configuration settings.

```bash
./bin/mcp-cli config
```

**Output:**
```
📋 Current Configuration
========================

  GROQ_API_KEY: gsk_SGaI...PYb8
  MCP Port: 9090
  Chatbot URL: localhost:9090
```

The configuration is stored in the `.env` file in your project root. You can edit it directly with any text editor.

---

### `mcp-go test-api`

**Test API Key**

Validates the Groq API key format and connectivity.

```bash
./bin/mcp-cli test-api
```

**Output (Success):**
```
🧪 Testing Groq API Key...

Sending test request to Groq API...
✅ API key format is valid
```

**Output (Failure):**
```
🧪 Testing Groq API Key...

Sending test request to Groq API...
❌ API test failed: GROQ_API_KEY not set. Run 'mcp-go init' first
```

---

### `mcp-go help`

**Show Help Information**

Displays comprehensive help about all commands.

```bash
./bin/mcp-cli help
```

---

## Configuration

### Environment File (.env)

The CLI stores configuration in a `.env` file in your project root:

```bash
GROQ_API_KEY=gsk_your_api_key_here
```

### Getting a Groq API Key

1. Visit https://console.groq.com/keys
2. Sign in or create an account
3. Create a new API key
4. Copy the key (starts with `gsk_`)
5. Use `mcp-go init` or edit `.env` directly

### Manual Configuration

Edit the `.env` file directly:

```bash
# Edit .env
nano .env

# Add or update:
GROQ_API_KEY=gsk_SGaI...PYb8
```

Then run:
```bash
./bin/mcp-cli status
```

to verify configuration is loaded.

## Examples

### Complete First-Time Setup

```bash
# 1. Clone repository
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go

# 2. Run setup wizard
./bin/mcp-cli init
# Follow prompts, enter your Groq API key

# 3. Start services
./bin/mcp-cli start
# Services now running!
```

### Running in Separate Terminals

```bash
# Terminal 1: Run MCP server
./bin/mcp-cli server

# Terminal 2: Run chatbot
./bin/mcp-cli chat
```

### Checking Service Health

```bash
# Check status
./bin/mcp-cli status

# Test API connectivity
./bin/mcp-cli test-api

# View configuration
./bin/mcp-cli config
```

### Building Everything

```bash
# Build all components
make build-all

# Or build individually
make build           # MCP server
make build-cli       # CLI utility
go build -o bin/chatbot ./cmd/chatbot  # Chatbot
```

### Adding to System PATH

```bash
# On macOS/Linux:
sudo cp bin/mcp-cli /usr/local/bin/mcp-go

# Then use from anywhere:
mcp-go status
mcp-go start
```

## Troubleshooting

### Problem: "GROQ_API_KEY not set"

**Solution:**
```bash
# Run setup wizard
./bin/mcp-cli init

# Or edit .env manually
echo "GROQ_API_KEY=gsk_your_key_here" > .env
```

### Problem: "MCP server is not running"

**Solution:**
```bash
# Start server in background first
./bin/mcp-cli server &

# Wait a moment, then start chatbot
./bin/mcp-cli chat
```

### Problem: "pgrep: command not found"

**Context:** Status checking might fail on some systems.

**Solution:**
- The CLI will use PID files if pgrep is unavailable
- Ensure server was started with CLI (saves PIDs)

### Problem: "Failed to build binaries"

**Solution:**
```bash
# Install Go dependencies
go mod tidy
go mod download

# Try building again
make build-all
```

### Problem: "Permission denied" when running CLI

**Solution:**
```bash
# Make binary executable
chmod +x bin/mcp-cli

# Run
./bin/mcp-cli status
```

## Architecture

The CLI is structured as follows:

```
mcp-go/
├── cmd/cli/
│   └── main.go              # CLI entry point and command handlers
├── internal/cli/
│   └── init.go              # Setup wizard and configuration
└── bin/
    ├── mcp-cli              # CLI binary
    ├── mcp-server           # MCP server binary
    └── chatbot              # Chatbot binary
```

## Performance

- **CLI startup time:** < 100ms
- **Server startup time:** < 500ms
- **Chatbot startup time:** < 1s
- **Binary size:**
  - CLI: ~2.5 MB
  - Server: ~4.8 MB
  - Chatbot: ~7.3 MB

## Security Notes

1. **API Key Storage:**
   - `.env` file is in `.gitignore`
   - Never commit API keys to version control
   - `.env` file has restricted permissions (0600)

2. **Process Management:**
   - Services run with user privileges
   - PIDs stored in `bin/` directory for cleanup
   - Graceful shutdown with SIGTERM

3. **Network:**
   - Server listens on `localhost:9090` by default
   - Not exposed to external network
   - Only accessible from local machine

## Next Steps

- Learn about [MCP Server Architecture](./MCP_ARCHITECTURE.md)
- Explore [Chatbot Features](./CHATBOT.md)
- Check [System Tools](./SYSTEM_TOOLS.md)
- Review [Security Framework](./SECURITY.md)

## Support

For issues, questions, or feature requests:

1. Check this documentation
2. Review [Troubleshooting](#troubleshooting) section
3. Open an issue on GitHub
4. Check existing issues for similar problems

## Version

- **CLI Version:** 1.0.0
- **Go Version:** 1.21+
- **Groq API:** Latest

---

**Last Updated:** March 31, 2024
**Maintainer:** Amarjit Singh

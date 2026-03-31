# MCP Go CLI - Quick Start (2 minutes)

Get MCP Go up and running in just 2 minutes!

## Prerequisites

- Go 1.21 or higher
- A Groq API key (get one free at https://console.groq.com/keys)

## Step 1: Clone & Build (30 seconds)

```bash
git clone https://github.com/amarjit-singh/mcp-go.git
cd mcp-go

# Build the CLI
make build-cli

# Or build everything
make build-all
```

## Step 2: Initial Setup (1 minute)

```bash
./bin/mcp-cli init
```

Follow the prompts:
- Press Enter to check prerequisites ✓
- Enter your Groq API key (starts with `gsk_`)
- Wait for binaries to build

**That's it!** Configuration saved to `.env`

## Step 3: Start Services (30 seconds)

```bash
./bin/mcp-cli start
```

You should see:
```
🚀 Starting MCP Go services...
📌 Starting MCP server...
   ✓ Server started (PID: XXXX)
💬 Starting chatbot...
   ✓ Chatbot started (PID: XXXX)
✅ All services started!
```

Now type commands in the chatbot! Examples:

```
You: What's my CPU usage?
Assistant: Your CPU usage is at 45.2%

You: List all running processes
Assistant: Here are the currently running processes:
1. mcp-server (PID: 12345)
2. chatbot (PID: 12346)
...

You: Check system health
Assistant: Your system is running well with 60% memory usage...

You: help
```

## Common Commands

```bash
# Check status
./bin/mcp-cli status

# View config
./bin/mcp-cli config

# Stop services
./bin/mcp-cli stop

# Get help
./bin/mcp-cli help
```

## Separate Terminal Mode

Want to run in different terminals?

```bash
# Terminal 1
./bin/mcp-cli server

# Terminal 2
./bin/mcp-cli chat
```

## Issues?

### "GROQ_API_KEY not set"
```bash
./bin/mcp-cli init
```

### "Permission denied"
```bash
chmod +x bin/mcp-cli
```

### Need more help?
```bash
./bin/mcp-cli help
```

---

That's all! You're ready to use MCP Go. 🚀

For detailed documentation, see [CLI.md](./CLI.md)

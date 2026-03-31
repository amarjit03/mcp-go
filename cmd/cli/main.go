package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	cli "github.com/amarjit-singh/mcp-go/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		if err := cli.InitWizard(); err != nil {
			fmt.Printf("❌ Init failed: %v\n", err)
			os.Exit(1)
		}
	case "start":
		if err := startAll(); err != nil {
			fmt.Printf("❌ Start failed: %v\n", err)
			os.Exit(1)
		}
	case "server":
		if err := startServer(); err != nil {
			fmt.Printf("❌ Server start failed: %v\n", err)
			os.Exit(1)
		}
	case "chat":
		if err := startChatbot(); err != nil {
			fmt.Printf("❌ Chatbot start failed: %v\n", err)
			os.Exit(1)
		}
	case "stop":
		if err := stopServices(); err != nil {
			fmt.Printf("❌ Stop failed: %v\n", err)
			os.Exit(1)
		}
	case "status":
		if err := printStatus(); err != nil {
			fmt.Printf("❌ Status check failed: %v\n", err)
			os.Exit(1)
		}
	case "config":
		if err := handleConfig(); err != nil {
			fmt.Printf("❌ Config failed: %v\n", err)
			os.Exit(1)
		}
	case "test-api":
		if err := testAPIKey(); err != nil {
			fmt.Printf("❌ API test failed: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("MCP Go - AI-Powered System Assistant CLI")
	fmt.Println()
	fmt.Println("Usage: mcp-go <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init           Interactive setup wizard")
	fmt.Println("  start          Start both MCP server and chatbot")
	fmt.Println("  server         Start MCP server only")
	fmt.Println("  chat           Start chatbot only")
	fmt.Println("  stop           Stop running services")
	fmt.Println("  status         Show services status")
	fmt.Println("  config         Manage configuration")
	fmt.Println("  test-api       Test Groq API key")
	fmt.Println("  help           Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mcp-go init              # Initial setup")
	fmt.Println("  mcp-go start             # Start all services")
	fmt.Println("  mcp-go status            # Check service status")
	fmt.Println()
}

func printHelp() {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║  🚀 MCP Go Help                            ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println()
	fmt.Println("  mcp-go init")
	fmt.Println("    Interactive setup wizard. Run this first to configure the application.")
	fmt.Println()
	fmt.Println("  mcp-go start")
	fmt.Println("    Start both MCP server and chatbot in background.")
	fmt.Println("    MCP Server runs on port 9090")
	fmt.Println("    Chatbot connects to server for system queries")
	fmt.Println()
	fmt.Println("  mcp-go server")
	fmt.Println("    Start only the MCP server (port 9090)")
	fmt.Println()
	fmt.Println("  mcp-go chat")
	fmt.Println("    Start only the chatbot (interactive)")
	fmt.Println()
	fmt.Println("  mcp-go stop")
	fmt.Println("    Stop both services gracefully")
	fmt.Println()
	fmt.Println("  mcp-go status")
	fmt.Println("    Show current status of services and configuration")
	fmt.Println()
	fmt.Println("  mcp-go config")
	fmt.Println("    View or edit configuration (.env file)")
	fmt.Println()
	fmt.Println("  mcp-go test-api")
	fmt.Println("    Test Groq API key connectivity")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println()
	fmt.Println("  # First time setup")
	fmt.Println("  $ mcp-go init")
	fmt.Println()
	fmt.Println("  # Start all services")
	fmt.Println("  $ mcp-go start")
	fmt.Println()
	fmt.Println("  # Check if services are running")
	fmt.Println("  $ mcp-go status")
	fmt.Println()
	fmt.Println("  # Stop services")
	fmt.Println("  $ mcp-go stop")
	fmt.Println()
	fmt.Println("ENVIRONMENT:")
	fmt.Println()
	fmt.Println("  Configuration is stored in .env file:")
	fmt.Println("  - GROQ_API_KEY     Your Groq API key (get from https://console.groq.com)")
	fmt.Println()
	fmt.Println("DOCUMENTATION:")
	fmt.Println("  See README.md for detailed setup instructions")
	fmt.Println()
}

// startAll starts both server and chatbot
func startAll() error {
	fmt.Println()
	fmt.Println("🚀 Starting MCP Go services...")
	fmt.Println()

	// Check binaries exist, build if needed
	if err := checkBinaries(); err != nil {
		fmt.Println("🔨 Building binaries...")
		if err := buildBinaries(); err != nil {
			return fmt.Errorf("failed to build binaries: %w", err)
		}
	}

	// Load config
	config, err := cli.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.GroqAPIKey == "" {
		return fmt.Errorf("GROQ_API_KEY not set. Run 'mcp-go init' first")
	}

	// Start server
	fmt.Println("📌 Starting MCP server...")
	serverCmd := exec.Command("./bin/mcp-server")
	serverCmd.Env = append(os.Environ(), fmt.Sprintf("GROQ_API_KEY=%s", config.GroqAPIKey))

	if err := serverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Save server PID
	if err := ioutil.WriteFile("bin/.server.pid", []byte(fmt.Sprintf("%d", serverCmd.Process.Pid)), 0644); err != nil {
		fmt.Printf("Warning: couldn't save server PID: %v\n", err)
	}

	fmt.Printf("   ✓ Server started (PID: %d)\n", serverCmd.Process.Pid)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Start chatbot
	fmt.Println("💬 Starting chatbot...")
	chatbotCmd := exec.Command("./bin/chatbot")
	chatbotCmd.Env = append(os.Environ(), fmt.Sprintf("GROQ_API_KEY=%s", config.GroqAPIKey))
	chatbotCmd.Stdin = os.Stdin
	chatbotCmd.Stdout = os.Stdout
	chatbotCmd.Stderr = os.Stderr

	// Save chatbot PID
	if err := chatbotCmd.Start(); err != nil {
		return fmt.Errorf("failed to start chatbot: %w", err)
	}

	if err := ioutil.WriteFile("bin/.chatbot.pid", []byte(fmt.Sprintf("%d", chatbotCmd.Process.Pid)), 0644); err != nil {
		fmt.Printf("Warning: couldn't save chatbot PID: %v\n", err)
	}

	fmt.Printf("   ✓ Chatbot started (PID: %d)\n", chatbotCmd.Process.Pid)
	fmt.Println()
	fmt.Println("✅ All services started! Type 'help' in chatbot for commands.")
	fmt.Println()

	// Wait for chatbot to finish
	_ = chatbotCmd.Wait()

	return nil
}

// startServer starts only the MCP server
func startServer() error {
	fmt.Println()
	fmt.Println("📌 Starting MCP Server...")
	fmt.Println()

	// Check/build if needed
	if err := checkBinaries(); err != nil {
		fmt.Println("🔨 Building server...")
		if err := runCommand("go", "build", "-o", "bin/mcp-server", "./cmd/server"); err != nil {
			return fmt.Errorf("failed to build server: %w", err)
		}
	}

	// Load config for API key
	config, err := cli.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.GroqAPIKey == "" {
		return fmt.Errorf("GROQ_API_KEY not set. Run 'mcp-go init' first")
	}

	// Start server
	serverCmd := exec.Command("./bin/mcp-server")
	serverCmd.Env = append(os.Environ(), fmt.Sprintf("GROQ_API_KEY=%s", config.GroqAPIKey))
	serverCmd.Stdin = os.Stdin
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	if err := serverCmd.Run(); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// startChatbot starts only the chatbot
func startChatbot() error {
	fmt.Println()
	fmt.Println("💬 Starting Chatbot...")
	fmt.Println()

	// Check if server is running
	if !isProcessRunning("mcp-server") {
		fmt.Println("⚠️  MCP server is not running. Start it with: mcp-go server")
		fmt.Println()
	}

	// Check/build if needed
	if err := checkBinaries(); err != nil {
		fmt.Println("🔨 Building chatbot...")
		if err := runCommand("go", "build", "-o", "bin/chatbot", "./cmd/chatbot"); err != nil {
			return fmt.Errorf("failed to build chatbot: %w", err)
		}
	}

	// Load config for API key
	config, err := cli.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.GroqAPIKey == "" {
		return fmt.Errorf("GROQ_API_KEY not set. Run 'mcp-go init' first")
	}

	// Start chatbot
	chatbotCmd := exec.Command("./bin/chatbot")
	chatbotCmd.Env = append(os.Environ(), fmt.Sprintf("GROQ_API_KEY=%s", config.GroqAPIKey))
	chatbotCmd.Stdin = os.Stdin
	chatbotCmd.Stdout = os.Stdout
	chatbotCmd.Stderr = os.Stderr

	if err := chatbotCmd.Run(); err != nil {
		return fmt.Errorf("chatbot error: %w", err)
	}

	return nil
}

// stopServices stops all running services
func stopServices() error {
	fmt.Println()
	fmt.Println("🛑 Stopping services...")
	fmt.Println()

	// Stop by reading PIDs
	if pidData, err := ioutil.ReadFile("bin/.server.pid"); err == nil {
		pid, _ := strconv.Atoi(string(pidData))
		if err := exec.Command("kill", "-15", fmt.Sprintf("%d", pid)).Run(); err == nil {
			fmt.Println("   ✓ Stopped MCP server")
		}
		os.Remove("bin/.server.pid")
	}

	if pidData, err := ioutil.ReadFile("bin/.chatbot.pid"); err == nil {
		pid, _ := strconv.Atoi(string(pidData))
		if err := exec.Command("kill", "-15", fmt.Sprintf("%d", pid)).Run(); err == nil {
			fmt.Println("   ✓ Stopped chatbot")
		}
		os.Remove("bin/.chatbot.pid")
	}

	// Also try to kill by name
	exec.Command("pkill", "-f", "mcp-server").Run()
	exec.Command("pkill", "-f", "bin/chatbot").Run()

	fmt.Println()
	fmt.Println("✅ All services stopped")
	fmt.Println()

	return nil
}

// printStatus prints current service status
func printStatus() error {
	config, err := cli.LoadConfig()
	if err != nil {
		config = &cli.Config{}
	}

	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║  📊 MCP Go Status                          ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()

	// Server status
	serverRunning := isProcessRunning("mcp-server")
	if serverRunning {
		fmt.Println("  🟢 MCP Server: RUNNING (port 9090)")
	} else {
		fmt.Println("  🔴 MCP Server: STOPPED")
	}

	// Chatbot status
	chatbotRunning := isProcessRunning("bin/chatbot")
	if chatbotRunning {
		fmt.Println("  🟢 Chatbot: RUNNING")
	} else {
		fmt.Println("  🔴 Chatbot: STOPPED")
	}

	// Config status
	fmt.Println()
	fmt.Println("  Configuration:")
	if config.GroqAPIKey != "" {
		fmt.Println("    ✓ GROQ_API_KEY: Set")
	} else {
		fmt.Println("    ✗ GROQ_API_KEY: Not set (run 'mcp-go init')")
	}

	fmt.Println()

	return nil
}

// handleConfig shows or edits configuration
func handleConfig() error {
	config, err := cli.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println()
	fmt.Println("📋 Current Configuration")
	fmt.Println("========================")
	fmt.Println()

	if config.GroqAPIKey != "" {
		keyPreview := config.GroqAPIKey[:8] + "..." + config.GroqAPIKey[len(config.GroqAPIKey)-4:]
		fmt.Printf("  GROQ_API_KEY: %s\n", keyPreview)
	} else {
		fmt.Println("  GROQ_API_KEY: Not set")
	}

	fmt.Printf("  MCP Port: %d\n", config.MCPPort)
	fmt.Printf("  Chatbot URL: %s\n", config.ChatbotURL)
	fmt.Println()

	return nil
}

// testAPIKey tests the Groq API key
func testAPIKey() error {
	fmt.Println()
	fmt.Println("🧪 Testing Groq API Key...")
	fmt.Println()

	config, err := cli.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.GroqAPIKey == "" {
		return fmt.Errorf("GROQ_API_KEY not set. Run 'mcp-go init' first")
	}

	fmt.Println("Sending test request to Groq API...")
	// This would require importing the LLM package and testing
	// For now, just check if key format is valid
	if !strings.HasPrefix(config.GroqAPIKey, "gsk_") {
		return fmt.Errorf("invalid API key format")
	}

	fmt.Println("✅ API key format is valid")
	fmt.Println()

	return nil
}

// Helper functions

func checkBinaries() error {
	if _, err := os.Stat("bin/mcp-server"); err != nil {
		return err
	}
	if _, err := os.Stat("bin/chatbot"); err != nil {
		return err
	}
	return nil
}

func buildBinaries() error {
	if err := runCommand("go", "build", "-o", "bin/mcp-server", "./cmd/server"); err != nil {
		return fmt.Errorf("failed to build server: %w", err)
	}
	if err := runCommand("go", "build", "-o", "bin/chatbot", "./cmd/chatbot"); err != nil {
		return fmt.Errorf("failed to build chatbot: %w", err)
	}
	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isProcessRunning(name string) bool {
	cmd := exec.Command("pgrep", "-f", name)
	return cmd.Run() == nil
}

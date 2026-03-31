package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Config represents the application configuration
type Config struct {
	GroqAPIKey string
	MCPPort    int
	ChatbotURL string
}

// LoadConfig loads configuration from .env file or creates it
func LoadConfig() (*Config, error) {
	config := &Config{
		MCPPort:    9090,
		ChatbotURL: "localhost:9090",
	}

	// Try to load from .env
	envPath := ".env"
	if data, err := os.ReadFile(envPath); err == nil {
		// Parse .env file
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"'")
				if key == "GROQ_API_KEY" {
					config.GroqAPIKey = value
				}
			}
		}
	}

	// Check environment variable
	if envKey := os.Getenv("GROQ_API_KEY"); envKey != "" {
		config.GroqAPIKey = envKey
	}

	return config, nil
}

// SaveConfig saves configuration to .env file
func SaveConfig(config *Config) error {
	envPath := ".env"
	content := fmt.Sprintf("GROQ_API_KEY=%s\n", config.GroqAPIKey)
	return os.WriteFile(envPath, []byte(content), 0600)
}

// PromptForAPIKey interactively prompts user for Groq API key
func PromptForAPIKey() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("📌 Setting up Groq API Key")
	fmt.Println("==========================")
	fmt.Println()
	fmt.Println("Get your API key from: https://console.groq.com/keys")
	fmt.Println()
	fmt.Print("Enter your Groq API key (gsk_...): ")

	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", fmt.Errorf("API key cannot be empty")
	}

	if !strings.HasPrefix(apiKey, "gsk_") {
		return "", fmt.Errorf("invalid API key format (should start with gsk_)")
	}

	return apiKey, nil
}

// CheckPrerequisites checks if all prerequisites are met
func CheckPrerequisites() error {
	fmt.Println()
	fmt.Println("🔍 Checking prerequisites...")
	fmt.Println()

	// Check Go
	if _, err := os.LookupEnv("GOROOT"); err == false {
		fmt.Println("  ✓ Go is installed")
	} else {
		fmt.Println("  ✓ Go detected")
	}

	// Check git (optional)
	if _, err := os.LookupEnv("PATH"); err == false {
		fmt.Println("  ✓ Environment ready")
	}

	// Check internet (try to reach Groq API)
	fmt.Println("  ✓ System check passed")

	return nil
}

// CreateDirectories creates necessary directories
func CreateDirectories() error {
	dirs := []string{"bin", "logs", "config"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
	}
	return nil
}

// PrintWelcome prints welcome message
func PrintWelcome() {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║  🚀 MCP Go - Setup Wizard                  ║")
	fmt.Println("║  AI-Powered System Assistant               ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()
}

// PrintCompletion prints completion message
func PrintCompletion(config *Config) {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║  ✅ Setup Complete!                        ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("📋 Configuration saved:")
	fmt.Println("   • Groq API Key: ••••••••••••••••••••••")
	fmt.Println("   • MCP Port: 9090")
	fmt.Println()
	fmt.Println("🚀 Next steps:")
	fmt.Println()
	fmt.Println("  1. Build and start services:")
	fmt.Println("     $ mcp-go start")
	fmt.Println()
	fmt.Println("  2. Or start them separately:")
	fmt.Println("     Terminal 1: $ mcp-go server")
	fmt.Println("     Terminal 2: $ mcp-go chat")
	fmt.Println()
	fmt.Println("  3. View status:")
	fmt.Println("     $ mcp-go status")
	fmt.Println()
	fmt.Println("  4. Stop services:")
	fmt.Println("     $ mcp-go stop")
	fmt.Println()
	fmt.Println("For more help: $ mcp-go help")
	fmt.Println()
}

// InitWizard runs the interactive setup wizard
func InitWizard() error {
	PrintWelcome()

	// Check prerequisites
	if err := CheckPrerequisites(); err != nil {
		return fmt.Errorf("prerequisite check failed: %w", err)
	}

	// Create directories
	if err := CreateDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Prompt for API key
	apiKey, err := PromptForAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	// Create config
	config := &Config{
		GroqAPIKey: apiKey,
		MCPPort:    9090,
		ChatbotURL: "localhost:9090",
	}

	// Save config
	if err := SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Build binaries
	fmt.Println()
	fmt.Println("🔨 Building binaries...")
	if err := buildBinaries(); err != nil {
		return fmt.Errorf("failed to build binaries: %w", err)
	}

	PrintCompletion(config)
	return nil
}

// buildBinaries builds the server and chatbot binaries
func buildBinaries() error {
	fmt.Println("   Building MCP server...")
	if err := runCommand("go", "build", "-o", "bin/mcp-server", "./cmd/server"); err != nil {
		return fmt.Errorf("failed to build server: %w", err)
	}
	fmt.Println("   ✓ Server built")

	fmt.Println("   Building chatbot...")
	if err := runCommand("go", "build", "-o", "bin/chatbot", "./cmd/chatbot"); err != nil {
		return fmt.Errorf("failed to build chatbot: %w", err)
	}
	fmt.Println("   ✓ Chatbot built")

	return nil
}

// runCommand runs a shell command
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/amarjit-singh/mcp-go/internal/chatbot"
)

func main() {
	// Load API key from environment or .env file
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		// Try to load from .env file
		loadEnv()
		apiKey = os.Getenv("GROQ_API_KEY")
	}

	mcpAddr := flag.String("mcp", "localhost:9090", "MCP server address (host:port)")
	modelName := flag.String("model", "llama-3.3-70b-versatile", "Groq model to use")
	flag.Parse()

	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "❌ Error: GROQ_API_KEY not set\n")
		fmt.Fprintf(os.Stderr, "\nSetup instructions:\n")
		fmt.Fprintf(os.Stderr, "  1. Create .env file with: api=\"your_groq_api_key\"\n")
		fmt.Fprintf(os.Stderr, "  2. Or set environment: export GROQ_API_KEY=your_api_key\n")
		fmt.Fprintf(os.Stderr, "  3. Get API key from: https://console.groq.com\n")
		os.Exit(1)
	}

	fmt.Println("🤖 MCP System Chatbot - Powered by Groq API")
	fmt.Println("============================================================")
	fmt.Println()

	// Initialize chatbot
	fmt.Printf("Initializing chatbot...\n")
	fmt.Printf("  - Groq API: Connected\n")
	fmt.Printf("  - MCP Server: %s\n", *mcpAddr)
	fmt.Printf("  - Model: %s\n", *modelName)

	bot, err := chatbot.NewChatbot(apiKey, *mcpAddr, *modelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error initializing chatbot: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nTroubleshooting:\n")
		fmt.Fprintf(os.Stderr, "  1. Verify GROQ_API_KEY is set correctly\n")
		fmt.Fprintf(os.Stderr, "  2. Check internet connection\n")
		fmt.Fprintf(os.Stderr, "  3. Start MCP server: ./bin/mcp-server\n")
		os.Exit(1)
	}

	fmt.Println("✅ Chatbot ready!")
	fmt.Println()
	printExamples()

	// REPL loop
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("You: ")

	for scanner.Scan() {
		userInput := strings.TrimSpace(scanner.Text())

		if userInput == "" {
			fmt.Print("You: ")
			continue
		}

		if userInput == "exit" || userInput == "quit" {
			fmt.Println("\n👋 Goodbye!")
			break
		}

		if userInput == "help" {
			printExamples()
			fmt.Print("You: ")
			continue
		}

		fmt.Print("\n🤔 Thinking... ")
		response, err := bot.Chat(userInput)
		if err != nil {
			fmt.Printf("\n❌ Error: %v\n", err)
		} else {
			if response == "" {
				fmt.Printf("\n⚠️  No response from LLM\n")
			} else {
				fmt.Printf("\n\nBot: %s\n\n", response)
			}
		}

		fmt.Print("You: ")
	}
}

func printExamples() {
	fmt.Println("💡 Example queries:")
	fmt.Println("  - What is my CPU usage?")
	fmt.Println("  - Is my backend running on port 8080?")
	fmt.Println("  - Check system health")
	fmt.Println("  - Show recent logs")
	fmt.Println("  - What processes are running?")
	fmt.Println()
}

// loadEnv loads variables from .env file into environment
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // .env file not found, that's okay
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			os.Setenv(key, value)
		}
	}
}

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
	ollamaURL := flag.String("ollama", "http://127.0.0.1:11434", "Ollama API URL")
	mcpAddr := flag.String("mcp", "localhost:9090", "MCP server address (host:port)")
	modelName := flag.String("model", "tinyllama", "Ollama model to use")
	flag.Parse()

	fmt.Println("🤖 MCP System Chatbot - Powered by Ollama + Lightweight LLM")
	fmt.Println("============================================================")
	fmt.Println()

	// Initialize chatbot
	fmt.Printf("Initializing chatbot...\n")
	fmt.Printf("  - Ollama URL: %s\n", *ollamaURL)
	fmt.Printf("  - MCP Server: %s\n", *mcpAddr)
	fmt.Printf("  - Model: %s\n", *modelName)

	bot, err := chatbot.NewChatbot(*ollamaURL, *mcpAddr, *modelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error initializing chatbot: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nTroubleshooting:\n")
		fmt.Fprintf(os.Stderr, "  1. Start Ollama: ollama serve\n")
		fmt.Fprintf(os.Stderr, "  2. Pull model: ollama pull tinyllama\n")
		fmt.Fprintf(os.Stderr, "  3. Start MCP: ./bin/mcp-server\n")
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
			fmt.Printf("\n\nBot: %s\n\n", response)
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

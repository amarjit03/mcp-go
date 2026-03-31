package chatbot

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/amarjit-singh/mcp-go/internal/llm"
	"github.com/amarjit-singh/mcp-go/pkg/mcp"
)

// Chatbot represents the system chatbot
type Chatbot struct {
	llmClient           *llm.OllamaClient
	mcpHost             string
	mcpPort             string
	model               string
	conversationHistory []Message
}

// Message represents a conversation message
type Message struct {
	Role    string
	Content string
}

// NewChatbot creates a new chatbot instance
func NewChatbot(ollamaURL, mcpAddr, modelName string) (*Chatbot, error) {
	llmClient := llm.NewOllamaClient(ollamaURL)

	// Check if Ollama is available
	if !llmClient.IsAvailable() {
		return nil, fmt.Errorf("Ollama not available at %s", ollamaURL)
	}

	// Get available models
	models, err := llmClient.GetModels()
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("no models available. Install a text generation model with: ollama pull mistral")
	}

	// Use specified model or first available
	model := modelName
	if model == "" {
		model = models[0]
	}

	// Warn if using embedding model (doesn't support generation)
	if isEmbeddingModel(model) {
		fmt.Printf("⚠️  WARNING: '%s' is an embedding model and cannot generate text!\n", model)
		fmt.Printf("Install a text generation model instead:\n")
		fmt.Printf("  ollama pull mistral       (7B, fast, recommended)\n")
		fmt.Printf("  ollama pull neural-chat   (7B, chat optimized)\n")
		fmt.Printf("  ollama pull llama2        (7B, very capable)\n")
		fmt.Printf("  ollama pull tinyllama     (1.1B, ultra-lightweight)\n")
		return nil, fmt.Errorf("model '%s' does not support text generation", model)
	}

	parts := strings.Split(mcpAddr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid MCP address format, expected host:port")
	}

	return &Chatbot{
		llmClient: llmClient,
		mcpHost:   parts[0],
		mcpPort:   parts[1],
		model:     model,
	}, nil
}

// Chat sends a message and returns a response
func (cb *Chatbot) Chat(userMessage string) (string, error) {
	// Add user message to history
	cb.conversationHistory = append(cb.conversationHistory, Message{
		Role:    "user",
		Content: userMessage,
	})

	// Determine and execute appropriate MCP tools
	tools := cb.determineTools(userMessage)
	toolResults := cb.executeMCPTools(tools)

	// Create system prompt
	systemPrompt := cb.buildSystemPrompt()

	// Build prompt for LLM
	fullPrompt := fmt.Sprintf(`%s

User Query: %s

System Data:
%s

Your Response (be concise, 2-3 sentences):`, systemPrompt, userMessage, toolResults)

	// Get LLM response
	llmResponse, err := cb.llmClient.Generate(fullPrompt, cb.model)
	if err != nil {
		return "", fmt.Errorf("LLM error: %w", err)
	}

	// Clean up response
	finalAnswer := strings.TrimSpace(llmResponse)

	// Add assistant response to history
	cb.conversationHistory = append(cb.conversationHistory, Message{
		Role:    "assistant",
		Content: finalAnswer,
	})

	return finalAnswer, nil
}

// buildSystemPrompt creates the system prompt
func (cb *Chatbot) buildSystemPrompt() string {
	return `You are a helpful Linux system assistant. You have real-time access to system metrics and can explain system status clearly. Be brief and technical.`
}

// determineTools determines which MCP tools to use based on user query
func (cb *Chatbot) determineTools(userMessage string) []string {
	query := strings.ToLower(userMessage)
	toolMap := map[string][]string{
		"cpu":       {"get_cpu_usage"},
		"processor": {"get_cpu_usage"},
		"memory":    {"get_memory_usage"},
		"ram":       {"get_memory_usage"},
		"port":      {"check_port", "health_check"},
		"running":   {"get_process_info", "health_check"},
		"process":   {"get_process_info"},
		"backend":   {"check_port", "health_check"},
		"health":    {"health_check"},
		"status":    {"health_check"},
		"log":       {"read_logs"},
		"error":     {"read_logs"},
		"file":      {"read_file"},
		"system":    {"health_check"},
	}

	var tools []string
	seen := make(map[string]bool)

	for keyword, keyTools := range toolMap {
		if strings.Contains(query, keyword) {
			for _, tool := range keyTools {
				if !seen[tool] {
					tools = append(tools, tool)
					seen[tool] = true
				}
			}
		}
	}

	// Default to health check if no specific tool found
	if len(tools) == 0 {
		tools = []string{"health_check"}
	}

	return tools
}

// executeMCPTools executes MCP tools and returns results
func (cb *Chatbot) executeMCPTools(tools []string) string {
	var results strings.Builder

	for _, toolName := range tools {
		result, err := cb.callMCPTool(toolName)
		if err != nil {
			results.WriteString(fmt.Sprintf("[%s] Error: %v\n", toolName, err))
			continue
		}
		results.WriteString(fmt.Sprintf("[%s]\n%s\n", toolName, result))
	}

	return results.String()
}

// callMCPTool calls an MCP tool via TCP
func (cb *Chatbot) callMCPTool(toolName string) (string, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", cb.mcpHost, cb.mcpPort))
	if err != nil {
		return "", fmt.Errorf("cannot connect to MCP server at %s:%s: %w", cb.mcpHost, cb.mcpPort, err)
	}
	defer conn.Close()

	args := cb.buildToolArguments(toolName)

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": args,
		},
		"id": 1,
	}

	reqBody, _ := json.Marshal(request)
	conn.Write(reqBody)

	var resp mcp.Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("tool error: %s", resp.Error.Message)
	}

	return cb.formatResult(resp.Result), nil
}

// buildToolArguments builds appropriate arguments for each tool
func (cb *Chatbot) buildToolArguments(toolName string) map[string]interface{} {
	args := make(map[string]interface{})

	switch toolName {
	case "health_check":
		args["ports"] = []int{22, 80, 443, 3000, 5432, 8080, 9090}
	case "check_port":
		args["port"] = 8080
		args["protocol"] = "tcp"
	case "read_logs":
		args["path"] = "/var/log/syslog"
		args["lines"] = 20
	case "read_file":
		args["path"] = "/etc/hostname"
	case "list_directory":
		args["path"] = "/home"
	}

	return args
}

// formatResult formats tool result for display
func (cb *Chatbot) formatResult(result interface{}) string {
	if callResp, ok := result.(mcp.CallToolResponse); ok {
		if len(callResp.Content) > 0 {
			data, _ := json.MarshalIndent(callResp.Content[0].Data, "", "  ")
			return string(data)
		}
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data)
}

// isEmbeddingModel checks if a model is an embedding model (not a text generation model)
func isEmbeddingModel(modelName string) bool {
	embeddingModels := map[string]bool{
		"mxbai-embed-large":      true,
		"nomic-embed-text":       true,
		"all-minilm":             true,
		"all-minilm:v2":          true,
		"snowflake-arctic-embed": true,
	}
	return embeddingModels[strings.Split(modelName, ":")[0]]
}

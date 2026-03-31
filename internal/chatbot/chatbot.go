package chatbot

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/amarjit-singh/mcp-go/internal/llm"
	"github.com/amarjit-singh/mcp-go/pkg/mcp"
)

// Chatbot represents the system chatbot
type Chatbot struct {
	llmClient           llm.LLMClient
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
func NewChatbot(apiKey, mcpAddr, modelName string) (*Chatbot, error) {
	llmClient := llm.NewGroqClient(apiKey)

	// Check if Groq API is available
	if !llmClient.IsAvailable() {
		return nil, fmt.Errorf("Groq API not available - verify GROQ_API_KEY is valid")
	}

	// Get available models
	models, err := llmClient.GetModels()
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("no models available from Groq API")
	}

	// Use specified model or first available
	model := modelName
	if model == "" {
		model = models[0]
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

	// Build prompt for LLM with clear data separation
	fullPrompt := fmt.Sprintf(`%s

QUESTION: %s

SYSTEM DATA:
%s

INSTRUCTIONS:
- Answer ONLY based on the system data provided
- Do NOT repeat the system data or raw JSON
- Provide a natural language response
- Be specific with numbers and values
- Keep response to 2-3 sentences maximum
- Do NOT show JSON or raw data in your response

ANSWER:`, systemPrompt, userMessage, toolResults)

	// Get LLM response
	llmResponse, err := cb.llmClient.Generate(fullPrompt, cb.model)
	if err != nil {
		return "", fmt.Errorf("LLM error: %w", err)
	}

	// DEBUG: Check what we got back
	if llmResponse == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] LLM returned empty response\n")
	}

	// Clean up response
	finalAnswer := strings.TrimSpace(llmResponse)

	// Remove any residual formatting that got included
	finalAnswer = strings.TrimPrefix(finalAnswer, "ANSWER:")
	finalAnswer = strings.TrimPrefix(finalAnswer, "ANSWER")
	finalAnswer = strings.TrimSpace(finalAnswer)

	// If response is empty or looks like it's still raw data, handle it
	if finalAnswer == "" || strings.HasPrefix(finalAnswer, "[") {
		// LLM failed to synthesize, provide a fallback
		finalAnswer = "Unable to synthesize response from system data."
	}

	// Add assistant response to history
	cb.conversationHistory = append(cb.conversationHistory, Message{
		Role:    "assistant",
		Content: finalAnswer,
	})

	return finalAnswer, nil
}

// buildSystemPrompt creates the system prompt
func (cb *Chatbot) buildSystemPrompt() string {
	return `You are a precise Linux system information assistant.

RULES:
1. Use ONLY the system data provided in the SYSTEM DATA section
2. Be concise - maximum 2-3 sentences
3. Provide specific numbers, percentages, and values
4. Do not speculate, hallucinate, or provide generic statements
5. Answer the question directly based on actual system data
6. If data is unavailable, say "Data not available" instead of guessing
7. Format: [Observation] [Reason] [Action if needed]

EXAMPLES:
Q: What is my CPU usage?
Data: CPU: 15.3%
A: Your CPU usage is at 15.3%, which is low and indicates your system has plenty of processing capacity available.

Q: Is port 8080 open?
Data: Port 8080: CLOSED
A: Port 8080 is currently closed - no service is listening on that port.`
}

// determineTools determines which MCP tools to use based on user query
func (cb *Chatbot) determineTools(userMessage string) []string {
	query := strings.ToLower(userMessage)
	toolMap := map[string][]string{
		"cpu":         {"get_cpu_usage"},
		"processor":   {"get_cpu_usage"},
		"usage":       {"get_cpu_usage", "get_memory_usage"},
		"memory":      {"get_memory_usage"},
		"ram":         {"get_memory_usage"},
		"disk":        {"health_check"},
		"port":        {"check_port"},
		"running":     {"get_process_info"},
		"process":     {"get_process_info"},
		"service":     {"get_process_info"},
		"backend":     {"check_port"},
		"application": {"check_port", "get_process_info"},
		"listen":      {"check_port"},
		"open":        {"check_port"},
		"health":      {"health_check"},
		"status":      {"health_check"},
		"system":      {"health_check"},
		"log":         {"read_logs"},
		"logs":        {"read_logs"},
		"error":       {"read_logs"},
		"file":        {"read_file"},
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
	addr := net.JoinHostPort(cb.mcpHost, cb.mcpPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("cannot connect to MCP server at %s: %w", addr, err)
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

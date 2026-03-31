package mcp

import "encoding/json"

// Version defines the MCP protocol version
const Version = "2024-11-05"

// RequestMethod represents an MCP request method
type RequestMethod string

const (
	MethodInitialize    RequestMethod = "initialize"
	MethodListTools     RequestMethod = "tools/list"
	MethodCallTool      RequestMethod = "tools/call"
	MethodListResources RequestMethod = "resources/list"
	MethodReadResource  RequestMethod = "resources/read"
)

// Request represents an MCP JSON-RPC request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  RequestMethod   `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

// Response represents an MCP JSON-RPC response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// Error represents an MCP error response
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// InitializeRequest is the parameters for initialize method
type InitializeRequest struct {
	ProtocolVersion string      `json:"protocolVersion"`
	Capabilities    interface{} `json:"capabilities,omitempty"`
	ClientInfo      ClientInfo  `json:"clientInfo"`
}

// InitializeResponse is the response for initialize method
type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    map[string]bool    `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
	Tools           []Tool             `json:"tools,omitempty"`
}

// ClientInfo describes the client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// ServerInfo describes the server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Tool describes an available tool
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema describes the parameters for a tool
type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

// CallToolRequest is the parameters for tools/call method
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// CallToolResponse is the response for tools/call method
type CallToolResponse struct {
	Content []ContentBlock `json:"content"`
}

// ContentBlock represents a content block in a response
type ContentBlock struct {
	Type string      `json:"type"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// ErrorCodes for MCP protocol
const (
	ErrorCodeParse      = -32700
	ErrorCodeInvalidReq = -32600
	ErrorCodeNotFound   = -32601
	ErrorCodeInvalidParams = -32602
	ErrorCodeInternal   = -32603
	ErrorCodeServerErr  = -32000
)

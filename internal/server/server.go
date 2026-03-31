package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amarjit-singh/mcp-go/internal/system"
	"github.com/amarjit-singh/mcp-go/internal/tools"
	"github.com/amarjit-singh/mcp-go/pkg/mcp"
)

// Server represents the MCP server
type Server struct {
	config      *Config
	protocol    *mcp.Protocol
	toolManager *tools.ToolManager
	listener    net.Listener
	logger      *log.Logger
}

// Config holds server configuration
type Config struct {
	Port             int
	Host             string
	Name             string
	Version          string
	CommandTimeout   time.Duration
	FileReadMaxSize  int64
	CommandMaxOutput int64
	AllowedPaths     []string
	AllowedCommands  []string
	DeniedCommands   []string
}

// DefaultConfig returns default server configuration
func DefaultConfig() *Config {
	return &Config{
		Port:             9090,
		Host:             "127.0.0.1",
		Name:             "mcp-dev-assistant",
		Version:          "1.0.0",
		CommandTimeout:   30 * time.Second,
		FileReadMaxSize:  50 * 1024 * 1024,
		CommandMaxOutput: 10 * 1024 * 1024,
	}
}

// NewServer creates a new MCP server
func NewServer(config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	logger := log.New(os.Stderr, "[MCP] ", log.LstdFlags)

	// Create system components
	fileReaderConfig := &system.FileReaderConfig{
		MaxFileSize:  config.FileReadMaxSize,
		AllowedPaths: config.AllowedPaths,
	}
	fileReader := system.NewFileReader(fileReaderConfig)

	cmdExecConfig := &system.ExecutorConfig{
		Timeout:         config.CommandTimeout,
		MaxOutputSize:   config.CommandMaxOutput,
		AllowedCommands: config.AllowedCommands,
		DeniedCommands:  config.DeniedCommands,
	}
	cmdExecutor := system.NewCommandExecutor(cmdExecConfig)

	systemMonitor := system.NewSystemMonitor()

	// Create tool manager
	toolManager := tools.NewToolManager(fileReader, cmdExecutor, systemMonitor)

	// Create protocol handler
	protocol := mcp.NewProtocol()

	server := &Server{
		config:      config,
		protocol:    protocol,
		toolManager: toolManager,
		logger:      logger,
	}

	// Register protocol handlers
	protocol.RegisterHandler(mcp.MethodInitialize, server.handleInitialize)
	protocol.RegisterHandler(mcp.MethodListTools, server.handleListTools)
	protocol.RegisterHandler(mcp.MethodCallTool, server.handleCallTool)

	return server
}

// Start starts the MCP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.listener = listener
	s.logger.Printf("Server listening on %s", addr)

	// Handle graceful shutdown
	go s.handleShutdown()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if listener was closed
			if _, ok := err.(net.Error); ok {
				return nil
			}
			s.logger.Printf("Accept error: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var rawReq json.RawMessage
		if err := decoder.Decode(&rawReq); err != nil {
			if err.Error() == "EOF" {
				return
			}
			s.logger.Printf("Decode error: %v", err)
			break
		}

		// Process request
		respData, err := s.protocol.HandleRequest(rawReq)
		if err != nil {
			s.logger.Printf("Handle request error: %v", err)
			continue
		}

		// Parse response for logging
		var resp mcp.Response
		if err := json.Unmarshal(respData, &resp); err == nil {
			if resp.Error != nil {
				s.logger.Printf("Response error [%d]: %s", resp.Error.Code, resp.Error.Message)
			}
		}

		// Send response
		if err := encoder.Encode(json.RawMessage(respData)); err != nil {
			s.logger.Printf("Encode error: %v", err)
			break
		}

		// Reset deadline
		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(params json.RawMessage) (interface{}, error) {
	var initReq mcp.InitializeRequest
	if err := json.Unmarshal(params, &initReq); err != nil {
		return nil, fmt.Errorf("invalid initialize params: %w", err)
	}

	s.logger.Printf("Client connected: %s v%s", initReq.ClientInfo.Name, initReq.ClientInfo.Version)

	return mcp.InitializeResponse{
		ProtocolVersion: mcp.Version,
		Capabilities: map[string]bool{
			"tools": true,
		},
		ServerInfo: mcp.ServerInfo{
			Name:    s.config.Name,
			Version: s.config.Version,
		},
		Tools: s.toolManager.GetTools(),
	}, nil
}

// handleListTools handles the tools/list request
func (s *Server) handleListTools(params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{
		"tools": s.toolManager.GetTools(),
	}, nil
}

// handleCallTool handles the tools/call request
func (s *Server) handleCallTool(params json.RawMessage) (interface{}, error) {
	var callReq mcp.CallToolRequest
	if err := json.Unmarshal(params, &callReq); err != nil {
		return nil, fmt.Errorf("invalid tool call params: %w", err)
	}

	s.logger.Printf("Calling tool: %s", callReq.Name)

	result, err := s.toolManager.ExecuteTool(callReq.Name, callReq.Arguments)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	return mcp.CallToolResponse{
		Content: []mcp.ContentBlock{
			{
				Type: "text",
				Data: result,
			},
		},
	}, nil
}

// handleShutdown handles graceful shutdown
func (s *Server) handleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	s.logger.Println("Shutdown signal received")

	if s.listener != nil {
		s.listener.Close()
	}

	os.Exit(0)
}

// Stop stops the server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amarjit-singh/mcp-go/internal/system"
	"github.com/amarjit-singh/mcp-go/pkg/mcp"
)

// ToolManager manages available tools and their execution
type ToolManager struct {
	tools         map[string]ToolFunc
	fileReader    *system.FileReader
	cmdExecutor   *system.CommandExecutor
	systemMonitor *system.SystemMonitor
}

// ToolFunc is a function that implements a tool
type ToolFunc func(args map[string]interface{}) (interface{}, error)

// NewToolManager creates a new tool manager
func NewToolManager(fileReader *system.FileReader, cmdExecutor *system.CommandExecutor, monitor *system.SystemMonitor) *ToolManager {
	tm := &ToolManager{
		tools:         make(map[string]ToolFunc),
		fileReader:    fileReader,
		cmdExecutor:   cmdExecutor,
		systemMonitor: monitor,
	}

	tm.registerBuiltinTools()
	return tm
}

// registerBuiltinTools registers all built-in tools
func (tm *ToolManager) registerBuiltinTools() {
	tm.RegisterTool("read_file", tm.readFile)
	tm.RegisterTool("write_file", tm.writeFile)
	tm.RegisterTool("list_directory", tm.listDirectory)
	tm.RegisterTool("execute_command", tm.executeCommand)
	tm.RegisterTool("get_cpu_usage", tm.getCPUUsage)
	tm.RegisterTool("get_memory_usage", tm.getMemoryUsage)
	tm.RegisterTool("check_port", tm.checkPort)
	tm.RegisterTool("get_process_info", tm.getProcessInfo)
	tm.RegisterTool("health_check", tm.healthCheck)
	tm.RegisterTool("read_logs", tm.readLogs)
}

// RegisterTool registers a new tool
func (tm *ToolManager) RegisterTool(name string, fn ToolFunc) {
	tm.tools[name] = fn
}

// ExecuteTool executes a tool by name
func (tm *ToolManager) ExecuteTool(name string, args map[string]interface{}) (interface{}, error) {
	tool, ok := tm.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool(args)
}

// GetTools returns a list of all available tools
func (tm *ToolManager) GetTools() []mcp.Tool {
	return []mcp.Tool{
		{
			Name:        "read_file",
			Description: "Read the contents of a file",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The path to the file to read",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "write_file",
			Description: "Write contents to a file",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The path to the file to write",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "The content to write",
					},
					"append": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to append to the file (default: false)",
					},
				},
				Required: []string{"path", "content"},
			},
		},
		{
			Name:        "list_directory",
			Description: "List the contents of a directory",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The path to the directory",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "execute_command",
			Description: "Execute a shell command",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The shell command to execute",
					},
				},
				Required: []string{"command"},
			},
		},
		{
			Name:        "get_cpu_usage",
			Description: "Get current CPU usage information",
			InputSchema: mcp.InputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
				Required:   []string{},
			},
		},
		{
			Name:        "get_memory_usage",
			Description: "Get current memory usage information",
			InputSchema: mcp.InputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
				Required:   []string{},
			},
		},
		{
			Name:        "check_port",
			Description: "Check if a network port is open",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"port": map[string]interface{}{
						"type":        "integer",
						"description": "The port number to check",
					},
					"protocol": map[string]interface{}{
						"type":        "string",
						"description": "Protocol to use: tcp or udp (default: tcp)",
					},
				},
				Required: []string{"port"},
			},
		},
		{
			Name:        "get_process_info",
			Description: "Get information about a running process",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"pid_or_name": map[string]interface{}{
						"type":        "string",
						"description": "Process ID (number) or process name (string)",
					},
				},
				Required: []string{"pid_or_name"},
			},
		},
		{
			Name:        "health_check",
			Description: "Perform a system health check",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"ports": map[string]interface{}{
						"type":        "array",
						"description": "Optional list of ports to check",
						"items": map[string]interface{}{
							"type": "integer",
						},
					},
				},
				Required: []string{},
			},
		},
		{
			Name:        "read_logs",
			Description: "Read log files (last N lines)",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the log file",
					},
					"lines": map[string]interface{}{
						"type":        "integer",
						"description": "Number of lines to read from the end (default: 50)",
					},
				},
				Required: []string{"path"},
			},
		},
	}
}

// readFile implements read_file tool
func (tm *ToolManager) readFile(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid path argument")
	}

	content, err := tm.fileReader.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"path":    path,
		"content": content,
	}, nil
}

// writeFile implements write_file tool
func (tm *ToolManager) writeFile(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid path argument")
	}

	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content argument")
	}

	append_, ok := args["append"].(bool)
	if !ok {
		append_ = false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	var flag int
	if append_ {
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	f, err := os.OpenFile(absPath, flag, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return map[string]interface{}{
		"path":         path,
		"bytesWritten": len(content),
		"success":      true,
	}, nil
}

// listDirectory implements list_directory tool
func (tm *ToolManager) listDirectory(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid path argument")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []map[string]interface{}
	for _, entry := range entries {
		info, _ := entry.Info()
		files = append(files, map[string]interface{}{
			"name":    entry.Name(),
			"isDir":   entry.IsDir(),
			"size":    info.Size(),
			"mode":    info.Mode().String(),
			"modTime": info.ModTime(),
		})
	}

	return map[string]interface{}{
		"path":  path,
		"files": files,
	}, nil
}

// executeCommand implements execute_command tool
func (tm *ToolManager) executeCommand(args map[string]interface{}) (interface{}, error) {
	command, ok := args["command"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid command argument")
	}

	output, err := tm.cmdExecutor.ExecuteCommand(command)

	result := map[string]interface{}{
		"command": command,
		"output":  output,
	}

	if err != nil {
		result["error"] = err.Error()
		result["success"] = false
	} else {
		result["success"] = true
	}

	return result, nil
}

// getCPUUsage implements get_cpu_usage tool
func (tm *ToolManager) getCPUUsage(args map[string]interface{}) (interface{}, error) {
	cpuInfo, err := tm.systemMonitor.GetCPUUsage()
	if err != nil {
		return nil, err
	}
	return cpuInfo, nil
}

// getMemoryUsage implements get_memory_usage tool
func (tm *ToolManager) getMemoryUsage(args map[string]interface{}) (interface{}, error) {
	memInfo, err := tm.systemMonitor.GetMemoryUsage()
	if err != nil {
		return nil, err
	}
	return memInfo, nil
}

// checkPort implements check_port tool
func (tm *ToolManager) checkPort(args map[string]interface{}) (interface{}, error) {
	port, ok := args["port"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid port argument")
	}

	protocol, _ := args["protocol"].(string)

	portStatus := tm.systemMonitor.CheckPort(int(port), protocol)
	return portStatus, nil
}

// getProcessInfo implements get_process_info tool
func (tm *ToolManager) getProcessInfo(args map[string]interface{}) (interface{}, error) {
	pidOrName, ok := args["pid_or_name"]
	if !ok {
		return nil, fmt.Errorf("invalid pid_or_name argument")
	}

	procInfo, err := tm.systemMonitor.GetProcessInfo(pidOrName)
	if err != nil {
		return nil, err
	}
	return procInfo, nil
}

// healthCheck implements health_check tool
func (tm *ToolManager) healthCheck(args map[string]interface{}) (interface{}, error) {
	var ports []int
	if portsArg, ok := args["ports"]; ok {
		if portsList, ok := portsArg.([]interface{}); ok {
			for _, p := range portsList {
				if portNum, ok := p.(float64); ok {
					ports = append(ports, int(portNum))
				}
			}
		}
	}

	hc, err := tm.systemMonitor.PerformHealthCheck(ports)
	if err != nil {
		return nil, err
	}
	return hc, nil
}

// readLogs implements read_logs tool
func (tm *ToolManager) readLogs(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid path argument")
	}

	lines := 50
	if linesArg, ok := args["lines"].(float64); ok {
		lines = int(linesArg)
	}

	content, err := tm.fileReader.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Split into lines and get last N
	fileLines := strings.Split(content, "\n")
	start := len(fileLines) - lines
	if start < 0 {
		start = 0
	}

	selectedLines := fileLines[start:]

	return map[string]interface{}{
		"path":       path,
		"linesRead":  len(selectedLines),
		"totalLines": len(fileLines),
		"lines":      selectedLines,
	}, nil
}

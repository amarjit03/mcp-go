package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CommandExecutor handles safe command execution
type CommandExecutor struct {
	config *ExecutorConfig
}

// ExecutorConfig configures command execution behavior
type ExecutorConfig struct {
	// Timeout for command execution
	Timeout time.Duration
	// AllowedCommands is a list of commands that can be executed (empty = allow all)
	AllowedCommands []string
	// DeniedCommands is a list of commands that cannot be executed
	DeniedCommands []string
	// MaxOutputSize is the maximum size of command output in bytes
	MaxOutputSize int64
}

// DefaultExecutorConfig returns a safe default configuration
func DefaultExecutorConfig() *ExecutorConfig {
	return &ExecutorConfig{
		Timeout:       30 * time.Second,
		MaxOutputSize: 10 * 1024 * 1024, // 10MB
		DeniedCommands: []string{
			"rm", "rmdir", "mkfs", "dd", "format", "fdisk",
			"reboot", "shutdown", "halt", "poweroff",
		},
	}
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(config *ExecutorConfig) *CommandExecutor {
	if config == nil {
		config = DefaultExecutorConfig()
	}
	return &CommandExecutor{config: config}
}

// ExecuteCommand runs a shell command with safety checks
func (ce *CommandExecutor) ExecuteCommand(cmdLine string) (string, error) {
	// Parse command
	parts := strings.Fields(cmdLine)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := parts[0]

	// Check denied commands
	if ce.isDenied(cmd) {
		return "", fmt.Errorf("command '%s' is not allowed", cmd)
	}

	// Check allowed commands if specified
	if len(ce.config.AllowedCommands) > 0 && !ce.isAllowed(cmd) {
		return "", fmt.Errorf("command '%s' is not in allowed list", cmd)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ce.config.Timeout)
	defer cancel()

	// Execute command
	execCmd := exec.CommandContext(ctx, "sh", "-c", cmdLine)
	output, err := execCmd.CombinedOutput()

	// Check output size
	if int64(len(output)) > ce.config.MaxOutputSize {
		output = output[:ce.config.MaxOutputSize]
	}

	if err != nil {
		// Don't treat context timeout as a fatal error
		if ctx.Err() == context.DeadlineExceeded {
			return string(output), fmt.Errorf("command timed out after %v", ce.config.Timeout)
		}
		return string(output), err
	}

	return string(output), nil
}

// isDenied checks if a command is in the denied list
func (ce *CommandExecutor) isDenied(cmd string) bool {
	baseName := filepath.Base(cmd)
	for _, denied := range ce.config.DeniedCommands {
		if baseName == denied {
			return true
		}
	}
	return false
}

// isAllowed checks if a command is in the allowed list
func (ce *CommandExecutor) isAllowed(cmd string) bool {
	baseName := filepath.Base(cmd)
	for _, allowed := range ce.config.AllowedCommands {
		if baseName == allowed {
			return true
		}
	}
	return false
}

// FileReader handles safe file reading
type FileReader struct {
	config *FileReaderConfig
}

// FileReaderConfig configures file reading behavior
type FileReaderConfig struct {
	// AllowedPaths is a list of allowed base paths (empty = allow all)
	AllowedPaths []string
	// DeniedPaths is a list of denied paths
	DeniedPaths []string
	// MaxFileSize is the maximum file size to read in bytes
	MaxFileSize int64
}

// DefaultFileReaderConfig returns a safe default configuration
func DefaultFileReaderConfig() *FileReaderConfig {
	return &FileReaderConfig{
		MaxFileSize: 50 * 1024 * 1024, // 50MB
		DeniedPaths: []string{
			"/etc/shadow",
			"/etc/gshadow",
			"/root/.ssh",
		},
	}
}

// NewFileReader creates a new file reader
func NewFileReader(config *FileReaderConfig) *FileReader {
	if config == nil {
		config = DefaultFileReaderConfig()
	}
	return &FileReader{config: config}
}

// ReadFile reads a file with safety checks
func (fr *FileReader) ReadFile(path string) (string, error) {
	// Normalize path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Check denied paths
	if fr.isDenied(absPath) {
		return "", fmt.Errorf("access denied to path: %s", path)
	}

	// Check allowed paths if specified
	if len(fr.config.AllowedPaths) > 0 && !fr.isAllowed(absPath) {
		return "", fmt.Errorf("path not in allowed directories: %s", path)
	}

	// Check if file exists
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("file not found: %w", err)
	}

	// Check if it's a file
	if fileInfo.IsDir() {
		return "", fmt.Errorf("path is a directory, not a file: %s", path)
	}

	// Check file size
	if fileInfo.Size() > fr.config.MaxFileSize {
		return "", fmt.Errorf("file too large: %d bytes (max %d)", fileInfo.Size(), fr.config.MaxFileSize)
	}

	// Read file
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// isDenied checks if a path is in the denied list
func (fr *FileReader) isDenied(absPath string) bool {
	for _, denied := range fr.config.DeniedPaths {
		if strings.HasPrefix(absPath, denied) {
			return true
		}
	}
	return false
}

// isAllowed checks if a path is in an allowed directory
func (fr *FileReader) isAllowed(absPath string) bool {
	for _, allowed := range fr.config.AllowedPaths {
		if strings.HasPrefix(absPath, allowed) {
			return true
		}
	}
	return false
}

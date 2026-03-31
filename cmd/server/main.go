package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amarjit-singh/mcp-go/internal/server"
	"gopkg.in/yaml.v3"
)

// ServerConfig holds YAML configuration
type ServerConfig struct {
	Port             int      `yaml:"port"`
	Host             string   `yaml:"host"`
	Name             string   `yaml:"name"`
	Version          string   `yaml:"version"`
	CommandTimeout   string   `yaml:"commandTimeout"`
	FileReadMaxSize  int64    `yaml:"fileReadMaxSize"`
	CommandMaxOutput int64    `yaml:"commandMaxOutput"`
	AllowedPaths     []string `yaml:"allowedPaths"`
	AllowedCommands  []string `yaml:"allowedCommands"`
	DeniedCommands   []string `yaml:"deniedCommands"`
}

func main() {
	configFile := flag.String("config", "", "Path to configuration file")
	port := flag.Int("port", 9090, "Server port")
	host := flag.String("host", "127.0.0.1", "Server host")
	flag.Parse()

	config := server.DefaultConfig()

	// Load from config file if provided
	if *configFile != "" {
		if err := loadConfigFile(*configFile, config); err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}

	// Override with command-line flags
	if *port != 9090 {
		config.Port = *port
	}
	if *host != "127.0.0.1" {
		config.Host = *host
	}

	// Create and start server
	srv := server.NewServer(config)

	fmt.Printf("Starting MCP Dev Assistant Server\n")
	fmt.Printf("Name: %s\n", config.Name)
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Host: %s\n", config.Host)
	fmt.Printf("Port: %d\n", config.Port)
	fmt.Println()

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// loadConfigFile loads configuration from a YAML file
func loadConfigFile(filePath string, serverConfig *server.Config) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg ServerConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply config values
	if cfg.Port > 0 {
		serverConfig.Port = cfg.Port
	}
	if cfg.Host != "" {
		serverConfig.Host = cfg.Host
	}
	if cfg.Name != "" {
		serverConfig.Name = cfg.Name
	}
	if cfg.Version != "" {
		serverConfig.Version = cfg.Version
	}
	if cfg.CommandTimeout != "" {
		timeout, err := time.ParseDuration(cfg.CommandTimeout)
		if err == nil {
			serverConfig.CommandTimeout = timeout
		}
	}
	if cfg.FileReadMaxSize > 0 {
		serverConfig.FileReadMaxSize = cfg.FileReadMaxSize
	}
	if cfg.CommandMaxOutput > 0 {
		serverConfig.CommandMaxOutput = cfg.CommandMaxOutput
	}
	if len(cfg.AllowedPaths) > 0 {
		serverConfig.AllowedPaths = cfg.AllowedPaths
	}
	if len(cfg.AllowedCommands) > 0 {
		serverConfig.AllowedCommands = cfg.AllowedCommands
	}
	if len(cfg.DeniedCommands) > 0 {
		serverConfig.DeniedCommands = cfg.DeniedCommands
	}

	return nil
}

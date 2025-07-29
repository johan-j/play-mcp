package main

import (
	"flag"
	"os"

	"github.com/johan-j/play-mcp/internal/plugins"
	"github.com/johan-j/play-mcp/internal/plugins/financial"
	"github.com/johan-j/play-mcp/internal/plugins/housing"
	"github.com/johan-j/play-mcp/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	// Command line flags
	var (
		host  = flag.String("host", "localhost", "Server host")
		port  = flag.Int("port", 8080, "Server port")
		debug = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	// Setup logger
	logger := logrus.New()
	if *debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Create plugin registry
	registry := plugins.NewRegistry()

	// Register plugins
	if err := registry.Register(financial.NewPlugin()); err != nil {
		logger.Fatalf("Failed to register financial plugin: %v", err)
	}
	logger.Info("Registered financial plugin")

	if err := registry.Register(housing.NewPlugin()); err != nil {
		logger.Fatalf("Failed to register housing plugin: %v", err)
	}
	logger.Info("Registered housing plugin")

	// Log registered tools
	tools := registry.GetAllTools()
	logger.Infof("Registered %d tools:", len(tools))
	for _, tool := range tools {
		logger.Infof("  - %s: %s", tool.Name, tool.Description)
	}

	// Create and start server
	mcpServer := server.NewMCPServer(registry, logger)

	config := server.Config{
		Host: *host,
		Port: *port,
	}

	logger.Info("Starting MCP server for GitHub Copilot...")
	if err := mcpServer.Start(config); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}

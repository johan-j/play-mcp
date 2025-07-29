package plugins

import (
	"context"
	"fmt"

	"github.com/johan-j/play-mcp/pkg/mcp"
)

// Registry manages all registered plugins
type Registry struct {
	plugins map[string]mcp.Plugin
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]mcp.Plugin),
	}
}

// Register registers a plugin with the registry
func (r *Registry) Register(plugin mcp.Plugin) error {
	name := plugin.Name()
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}
	r.plugins[name] = plugin
	return nil
}

// GetPlugin retrieves a plugin by name
func (r *Registry) GetPlugin(name string) (mcp.Plugin, bool) {
	plugin, exists := r.plugins[name]
	return plugin, exists
}

// GetAllPlugins returns all registered plugins
func (r *Registry) GetAllPlugins() map[string]mcp.Plugin {
	return r.plugins
}

// GetAllTools returns all tools from all registered plugins
func (r *Registry) GetAllTools() []mcp.Tool {
	var tools []mcp.Tool
	for _, plugin := range r.plugins {
		tools = append(tools, plugin.GetTools()...)
	}
	return tools
}

// GetAllResources returns all resources from all registered plugins
func (r *Registry) GetAllResources() []mcp.Resource {
	var resources []mcp.Resource
	for _, plugin := range r.plugins {
		resources = append(resources, plugin.GetResources()...)
	}
	return resources
}

// HandleToolCall routes a tool call to the appropriate plugin
func (r *Registry) HandleToolCall(ctx context.Context, request mcp.ToolCallRequest) (*mcp.ToolCallResponse, error) {
	// Find which plugin handles this tool
	for _, plugin := range r.plugins {
		for _, tool := range plugin.GetTools() {
			if tool.Name == request.Name {
				return plugin.HandleToolCall(ctx, request)
			}
		}
	}
	return nil, fmt.Errorf("tool %s not found", request.Name)
}

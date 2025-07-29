package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/johan-j/play-mcp/internal/plugins"
	"github.com/johan-j/play-mcp/pkg/mcp"
	"github.com/sirupsen/logrus"
)

// MCPServer represents the main MCP server
type MCPServer struct {
	registry *plugins.Registry
	logger   *logrus.Logger
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]*Client
	mutex    sync.RWMutex
}

// Client represents a connected MCP client
type Client struct {
	conn        *websocket.Conn
	initialized bool
	clientInfo  mcp.ClientInfo
}

// Config represents server configuration
type Config struct {
	Host string
	Port int
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(registry *plugins.Registry, logger *logrus.Logger) *MCPServer {
	return &MCPServer{
		registry: registry,
		logger:   logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		clients: make(map[*websocket.Conn]*Client),
	}
}

// Start starts the MCP server
func (s *MCPServer) Start(config Config) error {
	r := chi.NewRouter()

	// Add Chi middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Setup routes
	r.HandleFunc("/mcp", s.handleWebSocket)
	r.Post("/", s.handleHTTPMCP) // Add HTTP MCP handler for root path
	r.Get("/health", s.handleHealth)
	r.Get("/tools", s.handleToolsList)
	r.Get("/resources", s.handleResourcesList)

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	s.logger.Infof("Starting MCP server on %s", address)

	return http.ListenAndServe(address, r)
}

// handleWebSocket handles WebSocket connections for MCP protocol
func (s *MCPServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	client := &Client{
		conn:        conn,
		initialized: false,
	}

	s.mutex.Lock()
	s.clients[conn] = client
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.clients, conn)
		s.mutex.Unlock()
	}()

	s.logger.Info("New WebSocket connection established")

	for {
		var message mcp.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.Errorf("WebSocket error: %v", err)
			}
			break
		}

		s.logger.Debugf("Received message: %s", message.Method)

		response := s.handleMessage(client, &message)
		if response != nil {
			if err := conn.WriteJSON(response); err != nil {
				s.logger.Errorf("Failed to send response: %v", err)
				break
			}
		}
	}
}

// handleMessage processes incoming MCP messages
func (s *MCPServer) handleMessage(client *Client, message *mcp.Message) *mcp.Message {
	switch message.Method {
	case "initialize":
		return s.handleInitialize(client, message)
	case "tools/list":
		return s.handleToolsList2(client, message)
	case "tools/call":
		return s.handleToolCall(client, message)
	case "resources/list":
		return s.handleResourcesList2(client, message)
	default:
		return &mcp.Message{
			JSONRPC: "2.0",
			ID:      message.ID,
			Error: &mcp.Error{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", message.Method),
			},
		}
	}
}

// handleInitialize handles the initialize request
func (s *MCPServer) handleInitialize(client *Client, message *mcp.Message) *mcp.Message {
	var request mcp.InitializeRequest
	if err := json.Unmarshal(message.Params, &request); err != nil {
		return &mcp.Message{
			JSONRPC: "2.0",
			ID:      message.ID,
			Error: &mcp.Error{
				Code:    -32602,
				Message: "Invalid initialize parameters",
				Data:    err.Error(),
			},
		}
	}

	client.initialized = true
	client.clientInfo = request.ClientInfo

	response := mcp.InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: map[string]interface{}{
			"tools":     map[string]interface{}{},
			"resources": map[string]interface{}{},
		},
		ServerInfo: mcp.ServerInfo{
			Name:    "play-mcp-server",
			Version: "1.0.0",
		},
	}

	responseData, _ := json.Marshal(response)

	return &mcp.Message{
		JSONRPC: "2.0",
		ID:      message.ID,
		Result:  responseData,
	}
}

// handleToolsList2 handles tools/list requests via WebSocket
func (s *MCPServer) handleToolsList2(client *Client, message *mcp.Message) *mcp.Message {
	if !client.initialized {
		return &mcp.Message{
			JSONRPC: "2.0",
			ID:      message.ID,
			Error: &mcp.Error{
				Code:    -32002,
				Message: "Client not initialized",
			},
		}
	}

	tools := s.registry.GetAllTools()
	responseData, _ := json.Marshal(map[string]interface{}{
		"tools": tools,
	})

	return &mcp.Message{
		JSONRPC: "2.0",
		ID:      message.ID,
		Result:  responseData,
	}
}

// handleToolCall handles tools/call requests
func (s *MCPServer) handleToolCall(client *Client, message *mcp.Message) *mcp.Message {
	if !client.initialized {
		return &mcp.Message{
			JSONRPC: "2.0",
			ID:      message.ID,
			Error: &mcp.Error{
				Code:    -32002,
				Message: "Client not initialized",
			},
		}
	}

	var request mcp.ToolCallRequest
	if err := json.Unmarshal(message.Params, &request); err != nil {
		return &mcp.Message{
			JSONRPC: "2.0",
			ID:      message.ID,
			Error: &mcp.Error{
				Code:    -32602,
				Message: "Invalid tool call parameters",
				Data:    err.Error(),
			},
		}
	}

	ctx := context.Background()
	response, err := s.registry.HandleToolCall(ctx, request)
	if err != nil {
		return &mcp.Message{
			JSONRPC: "2.0",
			ID:      message.ID,
			Error: &mcp.Error{
				Code:    -32603,
				Message: "Tool call failed",
				Data:    err.Error(),
			},
		}
	}

	responseData, _ := json.Marshal(response)

	return &mcp.Message{
		JSONRPC: "2.0",
		ID:      message.ID,
		Result:  responseData,
	}
}

// handleResourcesList2 handles resources/list requests via WebSocket
func (s *MCPServer) handleResourcesList2(client *Client, message *mcp.Message) *mcp.Message {
	if !client.initialized {
		return &mcp.Message{
			JSONRPC: "2.0",
			ID:      message.ID,
			Error: &mcp.Error{
				Code:    -32002,
				Message: "Client not initialized",
			},
		}
	}

	resources := s.registry.GetAllResources()
	responseData, _ := json.Marshal(map[string]interface{}{
		"resources": resources,
	})

	return &mcp.Message{
		JSONRPC: "2.0",
		ID:      message.ID,
		Result:  responseData,
	}
}

// HTTP handlers for REST API access

// handleHealth handles health check requests
func (s *MCPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"plugins": len(s.registry.GetAllPlugins()),
		"tools":   len(s.registry.GetAllTools()),
	})
}

// handleToolsList handles HTTP requests for tools list
func (s *MCPServer) handleToolsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tools := s.registry.GetAllTools()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tools": tools,
	})
}

// handleResourcesList handles HTTP requests for resources list
func (s *MCPServer) handleResourcesList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resources := s.registry.GetAllResources()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"resources": resources,
	})
}

// handleHTTPMCP handles HTTP MCP JSON-RPC requests
func (s *MCPServer) handleHTTPMCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request mcp.JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.logger.Errorf("Failed to decode JSON-RPC request: %v", err)
		response := mcp.JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      nil,
			Error: &mcp.JSONRPCError{
				Code:    -32700,
				Message: "Parse error",
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	switch request.Method {
	case "initialize":
		s.handleHTTPInitialize(w, request)
	case "tools/list":
		s.handleHTTPToolsList(w, request)
	case "tools/call":
		s.handleHTTPToolCall(w, request)
	case "resources/list":
		s.handleHTTPResourcesList(w, request)
	default:
		response := mcp.JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &mcp.JSONRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
		json.NewEncoder(w).Encode(response)
	}
}

func (s *MCPServer) handleHTTPInitialize(w http.ResponseWriter, request mcp.JSONRPCRequest) {
	capabilities := mcp.ServerCapabilities{
		Resources: &mcp.ResourcesCapability{},
		Tools:     &mcp.ToolsCapability{},
	}

	serverInfo := mcp.ServerInfo{
		Name:    "play-mcp-server",
		Version: "1.0.0",
	}

	result := mcp.InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities:    capabilities,
		ServerInfo:      serverInfo,
	}

	response := mcp.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *MCPServer) handleHTTPToolsList(w http.ResponseWriter, request mcp.JSONRPCRequest) {
	tools := s.registry.GetAllTools()
	result := mcp.ToolsListResult{
		Tools: tools,
	}

	response := mcp.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *MCPServer) handleHTTPResourcesList(w http.ResponseWriter, request mcp.JSONRPCRequest) {
	resources := s.registry.GetAllResources()
	result := mcp.ResourcesListResult{
		Resources: resources,
	}

	response := mcp.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *MCPServer) handleHTTPToolCall(w http.ResponseWriter, request mcp.JSONRPCRequest) {
	var params mcp.ToolCallParams
	if err := json.Unmarshal(request.Params, &params); err != nil {
		response := mcp.JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &mcp.JSONRPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Execute the tool
	toolRequest := mcp.ToolCallRequest{
		Name:      params.Name,
		Arguments: params.Arguments,
	}
	result, err := s.registry.HandleToolCall(context.Background(), toolRequest)
	if err != nil {
		response := mcp.JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &mcp.JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := mcp.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}

	json.NewEncoder(w).Encode(response)
}

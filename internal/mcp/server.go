package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// ToolHandler is a function that executes a tool
type ToolHandler func(ctx context.Context, arguments map[string]interface{}) (interface{}, error)

// ResourceHandler is a function that reads a resource
type ResourceHandler func(ctx context.Context, uri string) (string, string, error) // content, mimeType, error

// Server implements the MCP protocol server
type Server struct {
	serverInfo Implementation

	tools     []Tool
	toolHandlers map[string]ToolHandler

	resources []Resource
	resourceHandlers map[string]ResourceHandler

	prompts []Prompt

	input  io.Reader
	output io.Writer
	logger *log.Logger
}

// NewServer creates a new MCP server
func NewServer(name, version, protocolVersion string) *Server {
	return &Server{
		serverInfo: Implementation{
			Name:    name,
			Version: version,
		},
		toolHandlers:     make(map[string]ToolHandler),
		resourceHandlers: make(map[string]ResourceHandler),
		input:            os.Stdin,
		output:           os.Stdout,
		logger:           log.New(os.Stderr, "[MCP] ", log.LstdFlags),
	}
}

// RegisterTool registers a tool with its handler
func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
	s.tools = append(s.tools, tool)
	s.toolHandlers[tool.Name] = handler
}

// RegisterResource registers a resource with its handler
func (s *Server) RegisterResource(resource Resource, handler ResourceHandler) {
	s.resources = append(s.resources, resource)
	s.resourceHandlers[resource.URI] = handler
}

// RegisterPrompt registers a prompt
func (s *Server) RegisterPrompt(prompt Prompt) {
	s.prompts = append(s.prompts, prompt)
}

// Run starts the MCP server (stdio transport)
func (s *Server) Run(ctx context.Context) error {
	s.logger.Println("MCP Server starting...")

	scanner := bufio.NewScanner(s.input)

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("Context cancelled, shutting down")
			return ctx.Err()
		default:
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("scanner error: %w", err)
				}
				// EOF reached
				if s.isDockerMode() {
					// In Docker, block waiting for shutdown signal
					s.logger.Println("Stdin EOF, waiting for shutdown signal...")
					<-ctx.Done()
					return nil
				}
				// In local mode, exit normally
				return nil
			}

			line := scanner.Bytes()

			// Parse JSON-RPC request
			var req JSONRPCRequest
			if err := json.Unmarshal(line, &req); err != nil {
				s.sendError(nil, -32700, "Parse error", err.Error())
				continue
			}

			// Handle request
			if err := s.handleRequest(ctx, &req); err != nil {
				s.logger.Printf("Error handling request: %v", err)
			}
		}
	}
}

// isDockerMode checks if running inside Docker container
func (s *Server) isDockerMode() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

func (s *Server) handleRequest(ctx context.Context, req *JSONRPCRequest) error {
	s.logger.Printf("Received: %s (id=%v)", req.Method, req.ID)

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "initialized":
		// Notification, no response needed
		s.logger.Println("Client initialized")
		return nil
	case "tools/list":
		return s.handleListTools(req)
	case "tools/call":
		return s.handleCallTool(ctx, req)
	case "resources/list":
		return s.handleListResources(req)
	case "resources/read":
		return s.handleReadResource(ctx, req)
	case "prompts/list":
		return s.handleListPrompts(req)
	case "prompts/get":
		return s.handleGetPrompt(req)
	case "ping":
		return s.sendResult(req.ID, map[string]interface{}{})
	default:
		return s.sendError(req.ID, -32601, "Method not found", req.Method)
	}
}

func (s *Server) handleInitialize(req *JSONRPCRequest) error {
	result := InitializeResult{
		ProtocolVersion: "2025-11-25",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
			Prompts: &PromptsCapability{
				ListChanged: false,
			},
			Logging: &LoggingCapability{},
		},
		ServerInfo: s.serverInfo,
	}

	return s.sendResult(req.ID, result)
}

func (s *Server) handleListTools(req *JSONRPCRequest) error {
	result := ListToolsResult{
		Tools: s.tools,
	}
	return s.sendResult(req.ID, result)
}

func (s *Server) handleCallTool(ctx context.Context, req *JSONRPCRequest) error {
	var params CallToolRequest
	paramsBytes, _ := json.Marshal(req.Params)
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return s.sendError(req.ID, -32602, "Invalid params", err.Error())
	}

	handler, ok := s.toolHandlers[params.Name]
	if !ok {
		return s.sendError(req.ID, -32602, "Tool not found", params.Name)
	}

	// Execute tool
	result, err := handler(ctx, params.Arguments)
	if err != nil {
		return s.sendResult(req.ID, CallToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Error: %v", err),
			}},
			IsError: true,
		})
	}

	// Convert result to string
	var text string
	switch v := result.(type) {
	case string:
		text = v
	default:
		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		text = string(jsonBytes)
	}

	return s.sendResult(req.ID, CallToolResult{
		Content: []Content{{
			Type: "text",
			Text: text,
		}},
		IsError: false,
	})
}

func (s *Server) handleListResources(req *JSONRPCRequest) error {
	result := ListResourcesResult{
		Resources: s.resources,
	}
	return s.sendResult(req.ID, result)
}

func (s *Server) handleReadResource(ctx context.Context, req *JSONRPCRequest) error {
	var params ReadResourceRequest
	paramsBytes, _ := json.Marshal(req.Params)
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return s.sendError(req.ID, -32602, "Invalid params", err.Error())
	}

	handler, ok := s.resourceHandlers[params.URI]
	if !ok {
		return s.sendError(req.ID, -32602, "Resource not found", params.URI)
	}

	content, mimeType, err := handler(ctx, params.URI)
	if err != nil {
		return s.sendError(req.ID, -32603, "Internal error", err.Error())
	}

	result := ReadResourceResult{
		Contents: []ResourceContents{{
			URI:      params.URI,
			MimeType: mimeType,
			Text:     content,
		}},
	}

	return s.sendResult(req.ID, result)
}

func (s *Server) handleListPrompts(req *JSONRPCRequest) error {
	result := ListPromptsResult{
		Prompts: s.prompts,
	}
	return s.sendResult(req.ID, result)
}

func (s *Server) handleGetPrompt(req *JSONRPCRequest) error {
	var params GetPromptRequest
	paramsBytes, _ := json.Marshal(req.Params)
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return s.sendError(req.ID, -32602, "Invalid params", err.Error())
	}

	// Find prompt
	var prompt *Prompt
	for i := range s.prompts {
		if s.prompts[i].Name == params.Name {
			prompt = &s.prompts[i]
			break
		}
	}

	if prompt == nil {
		return s.sendError(req.ID, -32602, "Prompt not found", params.Name)
	}

	// For now, return a simple prompt - extend later for actual templates
	result := GetPromptResult{
		Description: prompt.Description,
		Messages: []PromptMessage{{
			Role: "user",
			Content: PromptContent{
				Type: "text",
				Text: fmt.Sprintf("Execute prompt: %s", params.Name),
			},
		}},
	}

	return s.sendResult(req.ID, result)
}

func (s *Server) sendResult(id interface{}, result interface{}) error {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	return s.sendResponse(resp)
}

func (s *Server) sendError(id interface{}, code int, message string, data interface{}) error {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	return s.sendResponse(resp)
}

func (s *Server) sendResponse(resp JSONRPCResponse) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	data = append(data, '\n')

	if _, err := s.output.Write(data); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	s.logger.Printf("Sent response (id=%v)", resp.ID)
	return nil
}

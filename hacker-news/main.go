package main

import (
	"context"
	"encoding/json"

	"github.com/Djancyp/oido-studio/internal/plugins/rpc"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// ExtensionPluginServer implements the gRPC ExtensionPlugin service.
type ExtensionPluginServer struct {
	rpc.UnimplementedExtensionPluginServer
	hn *HackerNewsClient
}

// GetInfo returns plugin metadata.
func (s *ExtensionPluginServer) GetInfo(ctx context.Context, req *rpc.GetInfoRequest) (*rpc.GetInfoResponse, error) {
	return &rpc.GetInfoResponse{
		Info: &rpc.ExtensionInfo{
			Id:          "hacker-news",
			Name:        "Hacker News",
			Version:     "1.0.0",
			Description: "Fetch and browse Hacker News top stories, with configurable limit",
			Tags:        []string{"news", "hacker-news", "hn", "tech"},
		},
	}, nil
}

// GetTools returns the list of tools this plugin provides.
func (s *ExtensionPluginServer) GetTools(ctx context.Context, req *rpc.GetToolsRequest) (*rpc.GetToolsResponse, error) {
	return &rpc.GetToolsResponse{
		Tools: []*rpc.ToolDefinition{
			{
				Name:        "hn_top_stories",
				Description: "Fetch top stories from Hacker News. Returns story titles, URLs, scores, and authors. Use limit to control how many stories to return (default: 10, max: 30).",
				InputSchemaJson: `{
					"type": "object",
					"properties": {
						"limit": {
							"type": "number",
							"description": "Number of stories to return (1-30, default: 10)"
						}
					},
					"required": []
				}`,
				RequiresApproval: false,
				ApprovalMode:     "auto",
			},
			{
				Name:        "hn_story_detail",
				Description: "Get details for a specific Hacker News story by ID. Returns title, URL, score, author, text, and comment count.",
				InputSchemaJson: `{
					"type": "object",
					"properties": {
						"id": {
							"type": "number",
							"description": "Hacker News story ID"
						}
					},
					"required": ["id"]
				}`,
				RequiresApproval: false,
				ApprovalMode:     "auto",
			},
		},
	}, nil
}

// ExecuteTool handles tool execution requests.
func (s *ExtensionPluginServer) ExecuteTool(ctx context.Context, req *rpc.ExecuteToolRequest) (*rpc.ExecuteToolResponse, error) {
	switch req.Execution.ToolName {
	case "hn_top_stories":
		return s.toolTopStories(ctx, req.Execution.InputJson)
	case "hn_story_detail":
		return s.toolStoryDetail(ctx, req.Execution.InputJson)
	default:
		return &rpc.ExecuteToolResponse{
			Result: &rpc.ToolResult{
				Success: false,
				Error:   strPtr("unknown tool: " + req.Execution.ToolName),
			},
		}, nil
	}
}

// GetHooks returns hook handlers (none for this plugin).
func (s *ExtensionPluginServer) GetHooks(ctx context.Context, req *rpc.GetHooksRequest) (*rpc.GetHooksResponse, error) {
	return &rpc.GetHooksResponse{}, nil
}

// OnEvent handles event notifications (not used by this plugin).
func (s *ExtensionPluginServer) OnEvent(ctx context.Context, req *rpc.OnEventRequest) (*rpc.OnEventResponse, error) {
	return &rpc.OnEventResponse{}, nil
}

// toolTopStories fetches top HN stories with configurable limit.
func (s *ExtensionPluginServer) toolTopStories(ctx context.Context, inputJSON string) (*rpc.ExecuteToolResponse, error) {
	var input struct {
		Limit int `json:"limit"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return &rpc.ExecuteToolResponse{
			Result: &rpc.ToolResult{
				Success: false,
				Error:   strPtr("invalid input: " + err.Error()),
			},
		}, nil
	}

	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Limit > 30 {
		input.Limit = 30
	}

	stories, err := s.hn.GetTopStories(ctx, input.Limit)
	if err != nil {
		return &rpc.ExecuteToolResponse{
			Result: &rpc.ToolResult{
				Success: false,
				Error:   strPtr("failed to fetch stories: " + err.Error()),
			},
		}, nil
	}

	data, _ := json.MarshalIndent(stories, "", "  ")
	return &rpc.ExecuteToolResponse{
		Result: &rpc.ToolResult{
			Success:    true,
			OutputJson: string(data),
		},
	}, nil
}

// toolStoryDetail fetches a specific HN story by ID.
func (s *ExtensionPluginServer) toolStoryDetail(ctx context.Context, inputJSON string) (*rpc.ExecuteToolResponse, error) {
	var input struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return &rpc.ExecuteToolResponse{
			Result: &rpc.ToolResult{
				Success: false,
				Error:   strPtr("invalid input: " + err.Error()),
			},
		}, nil
	}

	if input.ID <= 0 {
		return &rpc.ExecuteToolResponse{
			Result: &rpc.ToolResult{
				Success: false,
				Error:   strPtr("id must be a positive integer"),
			},
		}, nil
	}

	story, err := s.hn.GetStory(ctx, input.ID)
	if err != nil {
		return &rpc.ExecuteToolResponse{
			Result: &rpc.ToolResult{
				Success: false,
				Error:   strPtr("failed to fetch story: " + err.Error()),
			},
		}, nil
	}

	data, _ := json.MarshalIndent(story, "", "  ")
	return &rpc.ExecuteToolResponse{
		Result: &rpc.ToolResult{
			Success:    true,
			OutputJson: string(data),
		},
	}, nil
}

// PluginGRPCPlugin implements go-plugin.Plugin and GRPCPlugin.
type PluginGRPCPlugin struct {
	plugin.Plugin
	Impl *ExtensionPluginServer
}

func (p *PluginGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	rpc.RegisterExtensionPluginServer(s, p.Impl)
	return nil
}

func (p *PluginGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return rpc.NewExtensionPluginClient(c), nil
}

// Handshake config
const (
	MagicCookieKey   = "OIDO_PLUGIN"
	MagicCookieValue = "oido_plugin"
)

func main() {
	hnClient := NewHackerNewsClient()

	extensionServer := &ExtensionPluginServer{
		hn: hnClient,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   MagicCookieKey,
			MagicCookieValue: MagicCookieValue,
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Plugins: map[string]plugin.Plugin{
			"extension": &PluginGRPCPlugin{Impl: extensionServer},
		},
	})
}

// Helper
func strPtr(s string) *string { return &s }

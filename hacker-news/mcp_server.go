package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewMCPHandler creates a new MCP handler for Hacker News tools.
func NewMCPHandler(hn *HackerNewsClient) *MCPHandler {
	return &MCPHandler{hn: hn}
}

// MCPHandler implements MCP tool handlers.
type MCPHandler struct {
	hn *HackerNewsClient
}

// TopStoriesArgs represents the arguments for hn_top_stories tool.
type TopStoriesArgs struct {
	Limit int `json:"limit" jsonschema:"Number of stories to return (1-30, default: 10)"`
}

// StoryDetailArgs represents the arguments for hn_story_detail tool.
type StoryDetailArgs struct {
	ID int `json:"id" jsonschema:"Hacker News story ID"`
}

// RunMCPServer starts the MCP server using stdio transport.
func RunMCPServer() {
	hnClient := NewHackerNewsClient()
	handler := NewMCPHandler(hnClient)

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "hacker-news",
		Version: "2.0.0",
	}, nil)

	// Register hn_top_stories tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "hn_top_stories",
		Description: "Fetch top stories from Hacker News. Returns story titles, URLs, scores, and authors. Use limit to control how many stories to return (default: 10, max: 30).",
	}, handler.HandleTopStories)

	// Register hn_story_detail tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "hn_story_detail",
		Description: "Get details for a specific Hacker News story by ID. Returns title, URL, score, author, text, and comment count.",
	}, handler.HandleStoryDetail)

	// Run server using stdio transport
	// This allows Qwen CLI to spawn it as a subprocess
	ctx := context.Background()
	log.Println("Hacker News MCP Server starting on stdio...")
	
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

// HandleTopStories fetches top HN stories.
func (h *MCPHandler) HandleTopStories(ctx context.Context, req *mcp.CallToolRequest, args TopStoriesArgs) (*mcp.CallToolResult, any, error) {
	limit := args.Limit
	
	if limit <= 0 {
		limit = 10
	}
	if limit > 30 {
		limit = 30
	}

	stories, err := h.hn.GetTopStories(ctx, limit)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to fetch stories: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	// Format stories as text
	result := fmt.Sprintf("# Hacker News Top Stories (%d stories)\n\n", len(stories))
	for i, story := range stories {
		result += fmt.Sprintf("%d. **%s**\n", i+1, story.Title)
		if story.URL != "" {
			result += fmt.Sprintf("   URL: %s\n", story.URL)
		}
		result += fmt.Sprintf("   Score: %d | Author: %s | Comments: %d\n\n", story.Score, story.Author, story.Descendants)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

// HandleStoryDetail fetches a specific HN story.
func (h *MCPHandler) HandleStoryDetail(ctx context.Context, req *mcp.CallToolRequest, args StoryDetailArgs) (*mcp.CallToolResult, any, error) {
	if args.ID <= 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Invalid story ID: must be a positive integer"},
			},
			IsError: true,
		}, nil, nil
	}

	story, err := h.hn.GetStory(ctx, args.ID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to fetch story: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	// Format story as text
	result := fmt.Sprintf("# %s\n\n", story.Title)
	if story.URL != "" {
		result += fmt.Sprintf("**URL:** %s\n\n", story.URL)
	}
	result += fmt.Sprintf("**Score:** %d  \n", story.Score)
	result += fmt.Sprintf("**Author:** %s  \n", story.Author)
	result += fmt.Sprintf("**Comments:** %d\n\n", story.Descendants)

	if story.Text != "" {
		result += fmt.Sprintf("**Text:**\n%s\n", story.Text)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

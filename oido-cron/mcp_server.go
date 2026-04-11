package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPHandler implements MCP tool handlers for cron management.
type MCPHandler struct {
	client *CronClient
}

func NewMCPHandler(client *CronClient) *MCPHandler {
	return &MCPHandler{client: client}
}

// RunMCPServer starts the MCP server using stdio transport.
func RunMCPServer() {
	client := NewCronClient()
	handler := NewMCPHandler(client)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "oido-cron",
		Version: "1.0.0",
	}, nil)

	// Register all cron tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "cron_list",
		Description: "List all scheduled cron jobs. Returns job ID, name, schedule, status, and delivery mode.",
	}, handler.HandleList)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "cron_get",
		Description: "Get details of a specific cron job by its ID.",
	}, handler.HandleGet)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "cron_add",
		Description: "Create a new scheduled cron job. Use schedule for recurring jobs (cron expression like '0 9 * * *'), interval_ms for fixed intervals (milliseconds), or at for one-shot jobs (e.g. '30m'). Delivery mode: 'none' (default), 'session', or 'channel' (requires channel and to).",
	}, handler.HandleAdd)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "cron_toggle",
		Description: "Enable or disable a cron job. Use this to pause a job without deleting it, or re-enable a paused job.",
	}, handler.HandleToggle)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "cron_run",
		Description: "Run a cron job immediately, regardless of its schedule. Useful for testing or triggering a job on demand.",
	}, handler.HandleRun)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "cron_logs",
		Description: "View execution logs for a specific cron job. Returns recent run history with status, output, and any errors.",
	}, handler.HandleLogs)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "cron_update",
		Description: "Update an existing cron job's properties. Only specified fields will be updated; others remain unchanged.",
	}, handler.HandleUpdate)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "cron_delete",
		Description: "Permanently delete a cron job. This action cannot be undone.",
	}, handler.HandleDelete)

	log.Println("OIDO Cron MCP Server starting on stdio...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

// --- Tool Handlers ---

type ListArgs struct{}

func (h *MCPHandler) HandleList(ctx context.Context, req *mcp.CallToolRequest, args ListArgs) (*mcp.CallToolResult, any, error) {
	jobs, err := h.client.ListJobs()
	if err != nil {
		return errorResult("Failed to list jobs: %v", err), nil, nil
	}

	if len(jobs) == 0 {
		return textResult("No cron jobs found."), nil, nil
	}

	result := "# Cron Jobs\n\n"
	result += fmt.Sprintf("| ID | Name | Schedule | Enabled | Delivery |\n")
	result += fmt.Sprintf("|----|------|----------|---------|----------|\n")
	for _, job := range jobs {
		enabled := "✅"
		if !job.Enabled {
			enabled = "❌"
		}
		result += fmt.Sprintf("| %d | %s | %s | %s | %s |\n", job.ID, job.Name, job.Schedule, enabled, job.Delivery)
	}

	return textResult(result), nil, nil
}

type GetArgs struct {
	ID int `json:"id" jsonschema:"The cron job ID"`
}

func (h *MCPHandler) HandleGet(ctx context.Context, req *mcp.CallToolRequest, args GetArgs) (*mcp.CallToolResult, any, error) {
	job, err := h.client.GetJob(int64(args.ID))
	if err != nil {
		return errorResult("Job not found: %v", err), nil, nil
	}

	data, _ := json.MarshalIndent(job, "", "  ")
	return textResult(string(data)), nil, nil
}

type AddArgs struct {
	Name       string `json:"name" jsonschema:"Job name (required)"`
	Message    string `json:"message" jsonschema:"Prompt for the agent (required)"`
	Schedule   string `json:"schedule" jsonschema:"Cron expression like '0 9 * * *' OR 'every|3600000' for intervals OR 'at|30m' for one-shot"`
	TZ         string `json:"tz" jsonschema:"Timezone like 'America/New_York' (default: UTC, only use with cron expressions)"`
	Model      string `json:"model" jsonschema:"Model override (optional)"`
	Session    string `json:"session" jsonschema:"Session target: isolated or main (default: isolated)"`
	Delivery   string `json:"delivery" jsonschema:"Delivery mode: none, session, or channel (default: none)"`
	Channel    string `json:"channel" jsonschema:"Channel type like 'whatsapp' (only use when delivery=channel)"`
	To         string `json:"to" jsonschema:"Recipient like '+1234567890' (only use when delivery=channel)"`
}

func (h *MCPHandler) HandleAdd(ctx context.Context, req *mcp.CallToolRequest, args AddArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errorResult("name is required"), nil, nil
	}
	if args.Message == "" {
		return errorResult("message is required"), nil, nil
	}

	createReq := CreateJobRequest{
		Name:    args.Name,
		Message: args.Message,
		Schedule: args.Schedule,
		TZ:       args.TZ,
		Model:    args.Model,
		Session:  args.Session,
		Delivery: args.Delivery,
		Channel:  args.Channel,
		To:       args.To,
	}

	id, err := h.client.CreateJob(createReq)
	if err != nil {
		return errorResult("Failed to create job: %v", err), nil, nil
	}

	return textResult(fmt.Sprintf("✅ Cron job created successfully\n\n**ID:** %d\n**Name:** %s", id, args.Name)), nil, nil
}

type ToggleArgs struct {
	ID      int  `json:"id" jsonschema:"The cron job ID"`
	Enabled bool `json:"enabled" jsonschema:"true to enable, false to disable"`
}

func (h *MCPHandler) HandleToggle(ctx context.Context, req *mcp.CallToolRequest, args ToggleArgs) (*mcp.CallToolResult, any, error) {
	err := h.client.ToggleJob(int64(args.ID), args.Enabled)
	if err != nil {
		return errorResult("Failed to toggle job: %v", err), nil, nil
	}

	action := "enabled"
	if !args.Enabled {
		action = "disabled"
	}

	return textResult(fmt.Sprintf("✅ Cron job #%d %s successfully", args.ID, action)), nil, nil
}

type RunArgs struct {
	ID int `json:"id" jsonschema:"The cron job ID"`
}

func (h *MCPHandler) HandleRun(ctx context.Context, req *mcp.CallToolRequest, args RunArgs) (*mcp.CallToolResult, any, error) {
	err := h.client.RunJob(int64(args.ID))
	if err != nil {
		return errorResult("Failed to run job: %v", err), nil, nil
	}

	return textResult(fmt.Sprintf("✅ Cron job #%d triggered successfully", args.ID)), nil, nil
}

type LogsArgs struct {
	ID    int `json:"id" jsonschema:"The cron job ID"`
	Limit int `json:"limit" jsonschema:"Number of log entries (default: 10)"`
}

func (h *MCPHandler) HandleLogs(ctx context.Context, req *mcp.CallToolRequest, args LogsArgs) (*mcp.CallToolResult, any, error) {
	limit := args.Limit
	if limit <= 0 {
		limit = 10
	}

	logs, err := h.client.GetJobLogs(int64(args.ID), limit)
	if err != nil {
		return errorResult("Failed to get logs: %v", err), nil, nil
	}

	if len(logs) == 0 {
		return textResult(fmt.Sprintf("No execution logs found for job #%d", args.ID)), nil, nil
	}

	result := fmt.Sprintf("# Execution Logs for Job #%d\n\n", args.ID)
	result += fmt.Sprintf("| Run ID | Started | Status | Duration |\n")
	result += fmt.Sprintf("|--------|---------|--------|----------|\n")
	for _, log := range logs {
		duration := "-"
		if log.CompletedAt != "" {
			duration = fmt.Sprintf("%s → %s", log.StartedAt, log.CompletedAt)
		}
		status := log.Status
		if status == "succeeded" {
			status = "✅ " + status
		} else if status == "failed" {
			status = "❌ " + status
		}
		result += fmt.Sprintf("| %d | %s | %s | %s |\n", log.ID, log.StartedAt, status, duration)
		if log.Error != "" {
			result += fmt.Sprintf("  **Error:** %s\n", log.Error)
		}
	}

	return textResult(result), nil, nil
}

type UpdateArgs struct {
	ID       int     `json:"id" jsonschema:"The cron job ID"`
	Name     *string `json:"name" jsonschema:"New job name"`
	Schedule *string `json:"schedule" jsonschema:"New cron expression"`
	TZ       *string `json:"tz" jsonschema:"New timezone"`
	Message  *string `json:"message" jsonschema:"New prompt for the agent"`
	Delivery *string `json:"delivery" jsonschema:"New delivery mode"`
}

func (h *MCPHandler) HandleUpdate(ctx context.Context, req *mcp.CallToolRequest, args UpdateArgs) (*mcp.CallToolResult, any, error) {
	updateReq := UpdateJobRequest{
		Name:     args.Name,
		Schedule: args.Schedule,
		TZ:       args.TZ,
		Message:  args.Message,
		Delivery: args.Delivery,
	}

	err := h.client.UpdateJob(int64(args.ID), updateReq)
	if err != nil {
		return errorResult("Failed to update job: %v", err), nil, nil
	}

	return textResult(fmt.Sprintf("✅ Cron job #%d updated successfully", args.ID)), nil, nil
}

type DeleteArgs struct {
	ID int `json:"id" jsonschema:"The cron job ID"`
}

func (h *MCPHandler) HandleDelete(ctx context.Context, req *mcp.CallToolRequest, args DeleteArgs) (*mcp.CallToolResult, any, error) {
	err := h.client.DeleteJob(int64(args.ID))
	if err != nil {
		return errorResult("Failed to delete job: %v", err), nil, nil
	}

	return textResult(fmt.Sprintf("✅ Cron job #%d deleted successfully", args.ID)), nil, nil
}

// --- Helper functions ---

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

func errorResult(format string, args ...interface{}) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf(format, args...)},
		},
		IsError: true,
	}
}

// Helper to convert string ID to int64
func parseID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

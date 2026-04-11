package main

import (
	"log"
)

// Main entry point for OIDO Cron MCP Server.
// This runs as a standalone process using stdio transport for Qwen CLI.
func main() {
	log.Println("Starting OIDO Cron MCP Server v1.0.0...")
	RunMCPServer()
}

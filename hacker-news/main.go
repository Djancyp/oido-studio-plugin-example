package main

import (
	"log"
)

// Main entry point for Hacker News MCP Server.
// This runs as a standalone process using stdio transport for Qwen CLI.
func main() {
	log.Println("Starting Hacker News MCP Server v2.0.0...")
	RunMCPServer()
}

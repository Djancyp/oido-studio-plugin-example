# Quick Start: Qwen CLI Extensions

## For Developers

### Build & Install Hacker News Extension

```bash
# Navigate to plugin directory
cd plugins/hacker-news

# Build MCP server (for Qwen CLI)
make build-mcp

# Package for Qwen CLI
make dist

# Link to Qwen CLI (development mode)
qwen extensions link $(pwd)/dist

# Verify installation
qwen extensions list
```

### Test the Extension

```bash
# Start Qwen CLI session
qwen

# In Qwen CLI, use the tools:
# - Ask: "What are the top HN stories?"
# - Or use command: /hn-top --limit 5
# - Or get story detail: /hn-story --id 12345
```

## For Users

### Install Extension

```bash
# From local path
qwen extensions install /path/to/plugins/hacker-news/dist

# Or from Git repository
qwen extensions install https://github.com/your-org/oido-studio-extensions

# List installed extensions
qwen extensions list
```

### Manage Extensions

```bash
# Update all extensions
qwen extensions update --all

# Disable an extension
qwen extensions disable hacker-news

# Enable an extension
qwen extensions enable hacker-news

# Uninstall an extension
qwen extensions uninstall hacker-news
```

## File Structure

```
plugins/hacker-news/
├── plugin.json              # ← oido-studio plugin (gRPC)
├── qwen-extension.json      # ← Qwen CLI extension (MCP)
├── main_grpc.go             # ← gRPC entry point
├── main_mcp.go              # ← MCP entry point
├── mcp_server.go            # ← MCP tool handlers
├── hn.go                    # ← Shared code
├── QWEN.md                  # ← Context for LLM
└── commands/                # ← Custom Qwen commands
    ├── hn-top.toml
    └── hn-story.toml
```

## Key Points

1. **Dual Purpose**: One plugin → oido-studio (gRPC) + Qwen CLI (MCP)
2. **Build Tags**: Use `//go:build mcp` to separate entry points
3. **Link Command**: `qwen extensions link <path>` for development
4. **MCP Transport**: Uses stdio (Qwen spawns the binary as subprocess)
5. **Auto-Detection**: Plugin manager detects `qwen-extension.json` automatically

## Troubleshooting

```bash
# Check if extension is detected
ls -la ~/.qwen/extensions/

# Check Qwen CLI logs
qwen --verbose

# Test MCP server manually
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | ./hacker-news-mcp
```

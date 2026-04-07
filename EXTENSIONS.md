# Qwen CLI Extension System for OIDO Studio Plugins

## Overview

OIDO Studio plugins are now **MCP-only** (Model Context Protocol) and serve as Qwen CLI extensions.

After uploading and installing a plugin to oido-studio, it **automatically** registers as a Qwen CLI extension:
```bash
# Automatic after upload - no manual step needed!
qwen extensions link <plugins-dir>/<plugin-name>
```

## Architecture

### Plugin Type

OIDO Studio uses **one plugin type** now:

| Type | Protocol | Purpose |
|------|----------|---------|
| **extension** | MCP (go-sdk) | LLM tools and hooks for Qwen CLI |

### Plugin Structure

```
plugins/hacker-news/
├── plugin.json              # oido-studio plugin manifest
├── qwen-extension.json      # Qwen CLI extension manifest
├── main.go                  # MCP server entry point
├── mcp_server.go            # MCP tool handlers
├── hn.go                    # Shared HN client code
├── Makefile                 # Build + package as zip
├── QWEN.md                  # Context file for Qwen CLI
└── commands/                # Custom Qwen CLI commands
    ├── hn-top.toml
    └── hn-story.toml
```

## Building Plugins

### Build & Package

```bash
cd plugins/hacker-news
make dist
```

Output: `dist/hacker-news.zip` ready for upload

### What's in the Zip

```
hacker-news.zip
├── plugin.json              # Manifest
├── qwen-extension.json      # Qwen manifest
├── QWEN.md                  # Context
├── hacker-news-mcp          # MCP server binary (12MB)
└── commands/
    ├── hn-top.toml
    └── hn-story.toml
```

## Installation Flow

### 1. Build the Extension

```bash
cd plugins/hacker-news
make dist
```

### 2. Upload via Plugins UI

- Open oido-studio web UI
- Navigate to **Plugins** page
- Click **Upload Plugin**
- Select `dist/hacker-news.zip`

### 3. Automatic Qwen Link

After upload, the plugin manager **automatically** runs:
```bash
qwen extensions link <plugins-dir>/hacker-news
```

No manual step needed! ✅

### 4. Verify

```bash
# Check in oido-studio UI
curl http://localhost:8080/api/plugins

# Check Qwen extensions
qwen extensions list
```

## qwen-extension.json Manifest

```json
{
  "name": "hacker-news",
  "version": "2.0.0",
  "description": "Fetch and browse Hacker News top stories via MCP server.",
  "mcpServers": {
    "hacker-news": {
      "command": "${extensionPath}/hacker-news-mcp",
      "args": [],
      "env": {}
    }
  },
  "contextFileName": "QWEN.md",
  "excludeTools": []
}
```

### Variables

| Variable | Description |
|----------|-------------|
| `${extensionPath}` | Full path to the extension directory |
| `${workspacePath}` | Full path to the current workspace |
| `${/}` or `${pathSeparator}` | OS-specific path separator (`/` or `\`) |

## MCP Server Implementation

### MCP Server Structure

```go
package main

import (
    "context"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Define argument structs with jsonschema tags
type TopStoriesArgs struct {
    Limit int `json:"limit" jsonschema:"Number of stories to return"`
}

// Tool handler with typed arguments
func (h *MCPHandler) HandleTopStories(
    ctx context.Context, 
    req *mcp.CallToolRequest, 
    args TopStoriesArgs,
) (*mcp.CallToolResult, any, error) {
    // Implementation
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: result},
        },
    }, nil, nil
}

func main() {
    server := mcp.NewServer(&mcp.Implementation{
        Name:    "hacker-news",
        Version: "2.0.0",
    }, nil)
    
    mcp.AddTool(server, &mcp.Tool{
        Name:        "hn_top_stories",
        Description: "Fetch top stories from Hacker News",
    }, handler.HandleTopStories)
    
    server.Run(context.Background(), &mcp.StdioTransport{})
}
```

## Qwen CLI Commands

Extensions can define custom commands via TOML files in the `commands/` directory.

### Example: `commands/hn-top.toml`

```toml
name = "hn-top"
description = "Fetch and display the top Hacker News stories."

prompt = """
Fetch the top {{limit}} stories from Hacker News using the hn_top_stories tool.
Format them as a numbered list with title, URL, score, author, and comment count.
"""

[parameters]
limit = { type = "number", description = "Number of stories to fetch", default = 10 }
```

### Command Namespacing

Subdirectories create namespaces:
- `commands/gcs/sync.toml` → `/gcs:sync` command
- `commands/hn-top.toml` → `/hn-top` command

### Conflict Resolution

Extension commands have the **lowest precedence**:
1. User/project commands keep original names
2. Extension commands get prefixed (e.g., `/hacker-news:hn-top`)

## Plugin Manager Integration

The plugin manager automatically:

1. Extracts and installs plugin from zip
2. Detects `qwen-extension.json`
3. Runs `qwen extensions link <path>` automatically
4. Sets `is_qwen_extension: true` in API response

### API Response Example

```json
{
  "id": "hacker-news",
  "name": "Hacker News",
  "type": "extension",
  "version": "2.0.0",
  "status": "running",
  "is_qwen_extension": true,
  "qwen_extension_path": "/path/to/plugins/hacker-news"
}
```

## Qwen Extension Management Commands

| Action | Command |
|--------|---------|
| Install from path | `qwen extensions install <url_or_path>` |
| Link for development | `qwen extensions link <path>` |
| Update | `qwen extensions update <name>` |
| Uninstall | `qwen extensions uninstall <name>` |
| Disable globally | `qwen extensions disable <name>` |
| Enable | `qwen extensions enable <name>` |
| List installed | `qwen extensions list` or `/extensions list` in CLI |

## Troubleshooting

### MCP Server Not Starting

1. Verify binary is in the zip:
   ```bash
   unzip -l hacker-news.zip | grep mcp
   ```

2. Check `qwen-extension.json` has correct command path:
   ```json
   {
     "mcpServers": {
       "hacker-news": {
         "command": "${extensionPath}/hacker-news-mcp"
       }
     }
   }
   ```

3. Check binary is executable:
   ```bash
   chmod +x hacker-news-mcp
   ```

### Qwen Extension Not Linked

1. Check if `qwen-extension.json` exists:
   ```bash
   ls plugins/hacker-news/qwen-extension.json
   ```

2. Manually link:
   ```bash
   qwen extensions link /path/to/plugins/hacker-news
   ```

3. Verify link:
   ```bash
   qwen extensions list
   ```

### qwen CLI Not Found

The plugin manager looks for `qwen` in PATH. If not found:

1. Install Qwen CLI
2. Ensure it's in PATH:
   ```bash
   which qwen
   ```

## Best Practices

1. **MCP Only**: No gRPC - use `modelcontextprotocol/go-sdk` with stdio transport
2. **Typed Arguments**: Use struct tags with `jsonschema` for schema generation
3. **Context File**: Always include `QWEN.md` for LLM instructions
4. **Commands**: Add custom commands in `commands/` directory
5. **Version Sync**: Keep `plugin.json` and `qwen-extension.json` versions in sync

## Example: Creating a New Extension

```bash
# 1. Create plugin directory
mkdir -p plugins/my-extension
cd plugins/my-extension

# 2. Initialize Go module
go mod init my-extension

# 3. Add MCP SDK dependency
go get github.com/modelcontextprotocol/go-sdk/mcp

# 4. Create files
touch plugin.json              # oido-studio manifest
touch qwen-extension.json      # Qwen manifest
touch main.go                  # MCP entry point
touch mcp_server.go            # MCP handlers
touch QWEN.md                  # Context file
mkdir commands                 # Custom commands

# 5. Write implementation
# (see hacker-news example)

# 6. Build and package
make dist

# 7. Upload via UI
# Open oido-studio UI → Plugins → Upload → select dist/my-extension.zip

# 8. Automatic link happens after upload!
```

## References

- [Qwen CLI Extensions Documentation](https://qwenlm.github.io/qwen-code-docs/en/developers/extensions/extension/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

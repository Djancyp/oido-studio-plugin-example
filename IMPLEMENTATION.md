# Plugin System - Final Implementation

## Summary of Changes

### ✅ Completed

**1. Removed gRPC Plugin System**
- Deleted `main_grpc.go` 
- Removed all `hashicorp/go-plugin` code
- Removed `GRPCLoader` from plugin manager
- Removed build tags (`//go:build mcp`)

**2. MCP-Only Plugin System**
- Using `modelcontextprotocol/go-sdk` with stdio transport
- Single binary: `hacker-news-mcp` (12MB)
- Located in `bin/` subdirectory after extraction

**3. Automatic Qwen Extension Management**
- ✅ **Install**: Auto-runs `qwen extensions link <path>` after upload (with auto-confirm)
- ✅ **Uninstall**: Auto-runs `qwen extensions uninstall <name>` when plugin is deleted
- ✅ **Detect**: Plugin manager detects `qwen-extension.json` and tracks extension status

**4. Zip Packaging with Directory Structure**
- ✅ Preserves `bin/`, `commands/`, and `skills/` directories
- ✅ Binary at `bin/hacker-news-mcp`
- ✅ Commands at `commands/*.toml`
- ✅ Skills at `skills/<name>/SKILL.md` - Teaches LLM how to use the extension

## Complete Flow

### Build

```bash
cd plugins/hacker-news
make dist
```

**Output:**
```
dist/hacker-news.zip
├── plugin.json
├── qwen-extension.json
├── QWEN.md
├── bin/
│   └── hacker-news-mcp      # MCP server binary
├── commands/
│   ├── hn-top.toml
│   └── hn-story.toml
└── skills/
    └── hacker-news/
        └── SKILL.md          # Teaches LLM how to use tools
```

### Upload & Auto-Link

1. Upload `dist/hacker-news.zip` via Plugins UI
2. Plugin manager extracts to `plugins/hacker-news/`
3. Detects `qwen-extension.json`
4. **Automatically** runs: `qwen extensions link plugins/hacker-news`
5. Extension ready in Qwen CLI! ✅

### Delete & Auto-Unlink

1. Click "Delete" on plugin in UI
2. Plugin manager removes from oido-studio
3. **Automatically** runs: `qwen extensions uninstall hacker-news`
4. Extension removed from Qwen CLI! ✅

## Key Code Changes

### 1. Auto-Link on Install

```go
// In InstallPlugin()
if HasQwenExtensionManifest(pluginDir) {
    if err := m.linkQwenExtension(pluginDir); err != nil {
        m.logger.Printf("[plugins] warning: failed to link Qwen extension: %v", err)
    }
}
```

### 2. Auto-Unlink on Uninstall

```go
// In UninstallPlugin()
if isQwenExtension {
    if err := m.unlinkQwenExtension(id); err != nil {
        m.logger.Printf("[plugins] warning: failed to unlink Qwen extension: %v", err)
    }
}
```

### 3. Auto-Confirm for Qwen CLI

```go
func (m *Manager) linkQwenExtension(pluginDir string) error {
    qwenPath, err := exec.LookPath("qwen")
    if err != nil {
        return fmt.Errorf("qwen CLI not found in PATH")
    }
    
    cmd := exec.Command(qwenPath, "extensions", "link", pluginDir)
    cmd.Stdin = bytes.NewReader([]byte("Y\n"))  // Auto-confirm
    output, err := cmd.CombinedOutput()
    // ...
}
```

### 4. Preserve Directory Structure in Zip

```go
func extractPluginFromZip(zipData []byte) (...) {
    // ...
    for _, f := range r.File {
        if f.FileInfo().IsDir() {
            continue  // Skip directories
        }
        
        // Preserve full path from zip
        cleanName := filepath.Clean(f.Name)
        files = append(files, {Name: cleanName, Data: data})
    }
}

func (m *Manager) InstallPlugin(...) {
    // Create parent directories as needed
    if dir := filepath.Dir(filePath); dir != pluginDir {
        os.MkdirAll(dir, 0755)
    }
}
```

## Plugin Structure

```
plugins/hacker-news/
├── plugin.json              # oido-studio manifest
├── qwen-extension.json      # Qwen CLI manifest
├── main.go                  # MCP server entry point
├── mcp_server.go            # MCP tool handlers
├── hn.go                    # Shared HN client code
├── Makefile                 # Build + package as zip
├── QWEN.md                  # Context file for Qwen CLI
└── commands/                # Custom Qwen CLI commands
    ├── hn-top.toml
    └── hn-story.toml
```

## qwen-extension.json

```json
{
  "name": "hacker-news",
  "version": "2.0.0",
  "description": "Fetch and browse Hacker News top stories via MCP server.",
  "mcpServers": {
    "hacker-news": {
      "command": "${extensionPath}/bin/hacker-news-mcp",
      "args": [],
      "env": {}
    }
  },
  "contextFileName": "QWEN.md",
  "excludeTools": []
}
```

## Qwen CLI Commands Reference

| Command | Description |
|---------|-------------|
| `qwen extensions install <path>` | Install extension from path |
| `qwen extensions link <path>` | Link extension (development mode) |
| `qwen extensions uninstall <name>` | Uninstall extension |
| `qwen extensions list` | List installed extensions |
| `qwen extensions update <name>` | Update extension |
| `qwen extensions disable <name>` | Disable extension |
| `qwen extensions enable <name>` | Enable extension |

## Testing

```bash
# Build plugin
cd plugins/hacker-news
make dist

# Verify zip structure
unzip -l dist/hacker-news.zip
# Should show: bin/hacker-news-mcp, commands/, etc.

# Upload via UI
# Open http://localhost:8080/plugins → Upload dist/hacker-news.zip

# Verify auto-link
qwen extensions list
# Should show: hacker-news

# Delete via UI
# Click "Delete" button

# Verify auto-unlink
qwen extensions list
# Should be empty or not show hacker-news
```

## Troubleshooting

### qwen CLI Not Found

```bash
# Check if qwen is in PATH
which qwen

# If not found, install Qwen CLI first
```

### Extension Not Auto-Linking

```bash
# Check logs
tail -f logs/oido-studio.log | grep "Qwen extension"

# Manual link
qwen extensions link /path/to/plugins/hacker-news
```

### Extension Not Auto-Unlinking

```bash
# Manual unlink
qwen extensions uninstall hacker-news
```

### Binary Path Issues

Verify `qwen-extension.json` has correct path:
```json
{
  "mcpServers": {
    "hacker-news": {
      "command": "${extensionPath}/bin/hacker-news-mcp"
    }
  }
}
```

## What's Removed

- ❌ gRPC plugin system (`hashicorp/go-plugin`)
- ❌ Build tags (`//go:build mcp`)
- ❌ Dual binaries
- ❌ Manual Qwen extension linking
- ❌ GRPCLoader from plugin manager

## What's Added

- ✅ MCP-only plugins (`modelcontextprotocol/go-sdk`)
- ✅ Auto-link on install (`qwen extensions link`)
- ✅ Auto-unlink on uninstall (`qwen extensions uninstall`)
- ✅ Directory structure preservation in zip
- ✅ `IsQwenExtension` field in PluginEntry
- ✅ Auto-confirm for Qwen CLI prompts

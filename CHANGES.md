# Plugin System Changes - Summary

## What Changed

### ✅ Removed
- **gRPC plugin system** - All `hashicorp/go-plugin` code removed
- **Build tags** - No more `//go:build mcp` separation needed
- **gRPC loader** - `GRPCLoader` struct and all related code removed
- **Dual binaries** - No more `hacker-news-plugin` (gRPC) + `hacker-news-mcp` (MCP)

### ✅ Added
- **MCP-only plugins** - Using `modelcontextprotocol/go-sdk` with stdio transport
- **Auto-link to Qwen** - After upload, automatically runs `qwen extensions link <path>`
- **Qwen extension detection** - Plugin manager detects `qwen-extension.json`
- **Zip packaging** - Single zip contains everything (binary + manifests + commands)

## Architecture

```
Before:
┌─────────────────────────────────────┐
│  Upload hacker-news.zip via UI      │
│  ↓                                   │
│  Extract to plugins/hacker-news/    │
│  ↓                                   │
│  Load via gRPC (hashicorp/go-plugin)│
│  ↓                                   │
│  Manual: qwen extensions link       │
└─────────────────────────────────────┘

After:
┌─────────────────────────────────────┐
│  Upload hacker-news.zip via UI      │
│  ↓                                   │
│  Extract to plugins/hacker-news/    │
│  ↓                                   │
│  Detect qwen-extension.json         │
│  ↓                                   │
│  Auto: qwen extensions link <path>  │
│  ↓                                   │
│  Ready in Qwen CLI! ✅              │
└─────────────────────────────────────┘
```

## File Changes

### Plugin Directory (`plugins/hacker-news/`)

| File | Status | Notes |
|------|--------|-------|
| `main_grpc.go` | ❌ Removed | gRPC entry point |
| `main_mcp.go` → `main.go` | ✅ Renamed | MCP entry point (no build tags) |
| `mcp_server.go` | ✅ Updated | Removed `//go:build mcp` tag |
| `plugin.json` | ✅ Updated | Binary: `hacker-news-mcp` |
| `qwen-extension.json` | ✅ Added | Qwen CLI manifest |
| `QWEN.md` | ✅ Added | Context file for LLM |
| `commands/` | ✅ Added | Custom Qwen commands |
| `Makefile` | ✅ Updated | Build + zip in one command |

### Plugin Manager (`app/internal/plugins/`)

| File | Changes |
|------|---------|
| `manager.go` | Removed gRPC loader, added `linkQwenExtension()` |
| `manifest.go` | Added `LoadQwenExtensionManifest()`, `HasQwenExtensionManifest()` |
| `types.go` | Added `QwenExtensionManifest`, `MCPServerConfig` structs |

## Build Flow

```bash
# Developer runs:
make dist

# Output:
dist/hacker-news.zip
├── plugin.json
├── qwen-extension.json
├── QWEN.md
├── hacker-news-mcp          # MCP server binary
└── commands/
    ├── hn-top.toml
    └── hn-story.toml

# Upload via UI → Automatic Qwen link!
```

## Key Code Changes

### 1. Plugin Manager Auto-Link

```go
// In InstallPlugin()
if HasQwenExtensionManifest(pluginDir) {
    if err := m.linkQwenExtension(pluginDir); err != nil {
        m.logger.Printf("[plugins] warning: failed to link Qwen extension: %v", err)
    }
}
```

### 2. Auto-Unlink on Delete

```go
// In UninstallPlugin()
if isQwenExtension {
    if err := m.unlinkQwenExtension(id); err != nil {
        m.logger.Printf("[plugins] warning: failed to unlink Qwen extension: %v", err)
    }
}
```

### 3. Link Function with Auto-Confirm

```go
func (m *Manager) linkQwenExtension(pluginDir string) error {
    qwenPath, err := exec.LookPath("qwen")
    if err != nil {
        return fmt.Errorf("qwen CLI not found in PATH")
    }
    
    cmd := exec.Command(qwenPath, "extensions", "link", pluginDir)
    cmd.Stdin = bytes.NewReader([]byte("Y\n"))  // Auto-confirm prompt
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("qwen extensions link failed: %w\nOutput: %s", err, string(output))
    }
    
    m.logger.Printf("[plugins] ✓ Linked Qwen extension: %s", pluginDir)
    return nil
}
```

### 4. Unlink Function

```go
func (m *Manager) unlinkQwenExtension(extensionName string) error {
    qwenPath, err := exec.LookPath("qwen")
    if err != nil {
        return fmt.Errorf("qwen CLI not found in PATH")
    }
    
    cmd := exec.Command(qwenPath, "extensions", "uninstall", extensionName)
    cmd.Stdin = bytes.NewReader([]byte("Y\n"))  // Auto-confirm prompt
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("qwen extensions uninstall failed: %w\nOutput: %s", err, string(output))
    }
    
    m.logger.Printf("[plugins] ✓ Unlinked Qwen extension: %s", extensionName)
    return nil
}
```

### 3. Plugin Entry Extension Tracking

```go
type PluginEntry struct {
    // ... existing fields ...
    IsQwenExtension     bool   `json:"is_qwen_extension,omitempty"`
    QwenExtensionPath   string `json:"qwen_extension_path,omitempty"`
}
```

## Testing

```bash
# Build plugin
cd plugins/hacker-news
make clean && make dist

# Verify zip
unzip -l dist/hacker-news.zip

# Upload via UI
# Open http://localhost:8080/plugins → Upload dist/hacker-news.zip

# Verify auto-link
qwen extensions list
# OR
ls ~/.qwen/extensions/hacker-news
```

## Migration Notes

- **No breaking changes to UI** - Plugins still display the same way
- **No breaking changes to API** - Same endpoints, just added `is_qwen_extension` field
- **Existing gRPC plugins** - Need to be migrated to MCP format
- **Qwen CLI requirement** - Must be installed and in PATH for auto-link to work

## Next Steps

1. Test upload flow with real UI
2. Verify auto-link works correctly
3. Test Qwen extension in CLI session
4. Add more MCP extensions (data-analysis, code-review, etc.)

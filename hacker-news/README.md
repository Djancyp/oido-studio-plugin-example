# Hacker News Plugin — Example Extension Plugin

This is an example **external extension plugin** for Oido Studio. It fetches top stories from Hacker News and provides tools that the LLM can use.

## Features

- **`hn_top_stories`** — Fetch top stories with configurable limit (1-30, default: 10)
- **`hn_story_detail`** — Get details for a specific story by ID

## Building

```bash
make build
# or
go build -o hacker-news-plugin .
```

## Installing

The plugin is already in the correct location (`plugins/hacker-news/`). The plugin manager will discover it automatically on startup.

## Testing

Start the Oido Studio app, then in a chat with the LLM:

```
Show me the top 5 hacker news stories
```

Or call the API directly:

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/plugins
```

## Architecture

```
hacker-news/
├── plugin.json           # Plugin manifest (discovered by host)
├── hacker-news-plugin    # Compiled binary (spawned by host via go-plugin)
├── main.go               # Plugin entry point + gRPC server
├── hn.go                 # Hacker News Firebase API client
├── Makefile              # Build helper
└── README.md             # This file
```

## How it works

1. **Discovery**: On startup, the host scans `plugins/` for directories with `plugin.json`
2. **Loading**: Host spawns `hacker-news-plugin` as a subprocess via `hashicorp/go-plugin`
3. **gRPC**: Host communicates with plugin over gRPC (auto-mTLS encrypted)
4. **Tool Registration**: Plugin reports its tools (`hn_top_stories`, `hn_story_detail`) to the host
5. **Execution**: When the LLM calls a tool, host sends gRPC request → plugin executes → returns result

## Creating Your Own Plugin

Copy this directory and modify:

1. `plugin.json` — Update `id`, `name`, `description`, `binary`, `capabilities`
2. `main.go` — Implement your tool logic in `ExecuteTool()`
3. `go.mod` — Update module path
4. Build: `go build -o <your-plugin-name>-plugin .`

# Swiss Army Knife Extension

Multi-tool extension providing **ripgrep** search, **fzf** fuzzy finding, and **git** object browsing for OIDO Studio.

## Available Tools

### `rg_search` — ripgrep Code Search

Fast recursive text search using [ripgrep](https://github.com/BurntSushi/ripgrep). Searches code with regex, file type filtering, and context lines.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `pattern` | string | **Yes** | Search pattern (regex by default) |
| `path` | string | No | Directory/file to search. Default: `.` |
| `file_type` | string | No | File type filter (e.g. `go`, `py`, `js`, `rust`, `json`). See `rg --type-list` |
| `case_sensitive` | boolean | No | Use case-sensitive matching. Default: smart case |
| `include_hidden` | boolean | No | Search hidden files/dirs. Default: false |
| `context_lines` | integer | No | Context lines before/after match (`-C`) |
| `max_count` | integer | No | Max matches to return (`-m`). Default: 100 |
| `use_pcre2` | boolean | No | Use PCRE2 engine for lookarounds/backreferences |
| `glob` | string | No | Glob filter (e.g. `*.go`, `*.{ts,tsx}`) |

**Examples:**
```
Search for function definitions in Go files:
rg_search("func.*Handler", file_type="go")

Search for TODO comments in TypeScript files with 2 context lines:
rg_search("TODO", file_type="typescript", context_lines=2)

Case-sensitive search for "password" in all files under src/:
rg_search("password", path="./src", case_sensitive=true)

Search for email pattern using PCRE2 lookbehind:
rg_search("(?<=email: )\\S+@\\S+", use_pcre2=true)

Find error messages in Python files:
rg_search("raise.*Error", file_type="py", max_count=20)
```

---

### `fzf_find_files` — Fuzzy File Finder

Fuzzy-find files in a directory using [fzf](https://github.com/junegunn/fzf).

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `path` | string | No | Root directory. Default: `.` |
| `query` | string | No | Fuzzy query. Supports `^prefix`, `$suffix`, `'exact'`, `!exclude` |
| `max_results` | integer | No | Max results. Default: 50 |
| `include_hidden` | boolean | No | Include hidden files. Default: false |

**Examples:**
```
Find all config files:
fzf_find_files(query="config")

Find TypeScript files in src/:
fzf_find_files(path="./src", query="ts")

Find files ending in .test.tsx:
fzf_find_files(query=".test.tsx$")

Find files but exclude node_modules:
fzf_find_files(query="index !node_modules")
```

---

### `fzf_filter_lines` — Fuzzy Line Filter

Filter lines of text using fzf fuzzy matching. Useful for filtering command output, logs, or any text.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `text` | string | **Yes** | Text to filter (each line is a candidate) |
| `query` | string | **Yes** | Fuzzy query to match against |
| `max_results` | integer | No | Max matching lines. Default: 50 |
| `exact` | boolean | No | Use exact substring matching. Default: false |

**Examples:**
```
Filter log lines containing errors:
fzf_filter_lines(text=errorLogOutput, query="ERROR")

Find specific imports from a file listing:
fzf_filter_lines(text=fileContents, query="import.*fmt")

Exact match for a specific line:
fzf_filter_lines(text=output, query="Build succeeded", exact=true)
```

---

### `fzf_git_browse` — Git Object Browser

Browse git objects: branches, tags, commits, stashes, remotes, reflogs, worktrees, and tracked files.

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | string | **Yes** | Object type: `branches`, `tags`, `commits`, `stashes`, `remotes`, `reflogs`, `worktrees`, `files` |
| `path` | string | No | Git repo path. Default: `.` |
| `query` | string | No | Fuzzy filter for results |
| `max_results` | integer | No | Max results. Default: 50 |

**Examples:**
```
List all branches:
fzf_git_browse(type="branches")

Find a specific commit:
fzf_git_browse(type="commits", query="fix auth bug")

Browse stashes:
fzf_git_browse(type="stashes")

Find tracked TypeScript files:
fzf_git_browse(type="files", query="ts")

List remotes:
fzf_git_browse(type="remotes")

Browse reflogs:
fzf_git_browse(type="reflogs", max_results=20)
```

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `RG_PATH` | Path to ripgrep binary (defaults to PATH lookup) |
| `FZF_PATH` | Path to fzf binary (defaults to PATH lookup) |

## Dependencies

- **ripgrep** (`rg`) — must be installed and on PATH
- **fzf** — must be installed and on PATH
- **git** — for `fzf_git_browse` tool

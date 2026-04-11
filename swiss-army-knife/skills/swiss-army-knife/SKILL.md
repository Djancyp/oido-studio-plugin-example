---
name: swiss-army-knife
description: Multi-tool extension for code search and file discovery
version: 1.0.0
tools:
  - name: rg_search
    description: Search code with ripgrep — fast regex search across files
  - name: fzf_find_files
    description: Fuzzy-find files in directories
  - name: fzf_filter_lines
    description: Filter text lines with fuzzy matching
  - name: fzf_git_browse
    description: Browse git objects — branches, tags, commits, stashes, etc
---

# Swiss Army Knife

A swiss army knife for code navigation. Provides four tools for searching and browsing code and git repositories.

## Tools

### rg_search
Search code with ripgrep. Supports regex, file type filtering, context lines, and PCRE2.
Use when: searching for patterns in code, finding function definitions, locating TODOs/errors.

### fzf_find_files
Fuzzy-find files in directories.
Use when: looking for files by name, browsing project structure, finding config files.

### fzf_filter_lines
Filter text lines with fuzzy matching.
Use when: filtering command output, finding specific lines in logs or large text.

### fzf_git_browse
Browse git objects — branches, tags, commits, stashes, remotes, reflogs, worktrees, files.
Use when: exploring git history, finding branches, listing tracked files, browsing stashes.

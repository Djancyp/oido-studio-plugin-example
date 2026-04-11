package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ---- Tool Definitions ----

func getToolDefinitions() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "rg_search",
			"description": "Search for text patterns in files using ripgrep (rg). Fast recursive code search with regex support, file type filtering, and context lines. Supports JSON output for structured results.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pattern": map[string]interface{}{
						"type":        "string",
						"description": "The search pattern (regex by default). Use whole-word matching with -w flag.",
					},
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Directory or file path to search in. Defaults to current directory.",
					},
					"file_type": map[string]interface{}{
						"type":        "string",
						"description": "Filter by file type (e.g. 'go', 'py', 'js', 'rust', 'json', 'yaml'). Use -T to exclude. See 'rg --type-list' for all types.",
					},
					"case_sensitive": map[string]interface{}{
						"type":        "boolean",
						"description": "If true, use case-sensitive matching (-s). Defaults to smart case (auto).",
					},
					"include_hidden": map[string]interface{}{
						"type":        "boolean",
						"description": "If true, search hidden files and directories (--hidden). Defaults to false.",
					},
					"context_lines": map[string]interface{}{
						"type":        "integer",
						"description": "Number of context lines before and after each match (-C). Defaults to 0.",
					},
					"max_count": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of matches to return (-m). Defaults to 100.",
					},
					"use_pcre2": map[string]interface{}{
						"type":        "boolean",
						"description": "If true, use PCRE2 regex engine (-P) for lookarounds and backreferences.",
					},
					"glob": map[string]interface{}{
						"type":        "string",
						"description": "Glob pattern to filter files (e.g. '*.go', '*.{ts,tsx}'). Multiple patterns separated by comma.",
					},
				},
				"required": []string{"pattern"},
			},
		},
		{
			"name":        "fzf_find_files",
			"description": "Fuzzy-find files in a directory using fzf. Interactive fuzzy matching file search with support for query filtering, multi-selection, and path scoping.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Root directory to search. Defaults to current directory.",
					},
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Initial fuzzy query to filter results. Supports fzf extended syntax: ^prefix, $suffix, 'exact', !exclude.",
					},
					"max_results": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results to return. Defaults to 50.",
					},
					"include_hidden": map[string]interface{}{
						"type":        "boolean",
						"description": "Include hidden files/directories in results. Defaults to false.",
					},
				},
			},
		},
		{
			"name":        "fzf_filter_lines",
			"description": "Filter lines of text using fzf fuzzy matching. Takes text input and returns lines matching a fuzzy query. Useful for filtering command output, logs, or any text stream.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "The text to filter. Each line is a candidate for fuzzy matching.",
					},
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Fuzzy query to match against. Supports fzf extended syntax.",
					},
					"max_results": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of matching lines to return. Defaults to 50.",
					},
					"exact": map[string]interface{}{
						"type":        "boolean",
						"description": "If true, use exact substring matching instead of fuzzy. Defaults to false.",
					},
				},
				"required": []string{"text", "query"},
			},
		},
		{
			"name":        "fzf_git_browse",
			"description": "Browse git objects interactively: branches, tags, commits, stashes, remotes, reflogs, worktrees, and files. Returns a list of git objects with details.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"type": map[string]interface{}{
						"type":        "string",
						"description": "Type of git object to browse: 'branches', 'tags', 'commits', 'stashes', 'remotes', 'reflogs', 'worktrees', 'files'.",
						"enum":        []string{"branches", "tags", "commits", "stashes", "remotes", "reflogs", "worktrees", "files"},
					},
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the git repository. Defaults to current directory.",
					},
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Fuzzy query to filter results. Defaults to no filter.",
					},
					"max_results": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results to return. Defaults to 50.",
					},
				},
				"required": []string{"type"},
			},
		},
	}
}

// ---- Tool Execution ----

func executeTool(name string, args map[string]interface{}) (map[string]interface{}, error) {
	switch name {
	case "rg_search":
		return execRgSearch(args)
	case "fzf_find_files":
		return execFzfFindFiles(args)
	case "fzf_filter_lines":
		return execFzfFilterLines(args)
	case "fzf_git_browse":
		return execFzfGitBrowse(args)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// ---- ripgrep Search ----

func execRgSearch(args map[string]interface{}) (map[string]interface{}, error) {
	pattern, _ := args["pattern"].(string)
	if pattern == "" {
		return nil, fmt.Errorf("pattern is required")
	}

	searchPath := "."
	if p, ok := args["path"].(string); ok && p != "" {
		searchPath = p
	}

	maxCount := 100
	if mc, ok := args["max_count"].(float64); ok && mc > 0 {
		maxCount = int(mc)
	}

	contextLines := 0
	if cl, ok := args["context_lines"].(float64); ok {
		contextLines = int(cl)
	}

	rgPath := getRgPath()

	cmdArgs := []string{
		"--json", "--no-heading", "--color=never",
		"-n", // line numbers
		fmt.Sprintf("-m=%d", maxCount),
	}

	if caseSensitive, _ := args["case_sensitive"].(bool); caseSensitive {
		cmdArgs = append(cmdArgs, "-s")
	}
	if includeHidden, _ := args["include_hidden"].(bool); includeHidden {
		cmdArgs = append(cmdArgs, "--hidden")
	}
	if contextLines > 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("-C=%d", contextLines))
	}
	if usePcre2, _ := args["use_pcre2"].(bool); usePcre2 {
		cmdArgs = append(cmdArgs, "-P")
	}
	if fileType, ok := args["file_type"].(string); ok && fileType != "" {
		cmdArgs = append(cmdArgs, "-t", fileType)
	}
	if glob, ok := args["glob"].(string); ok && glob != "" {
		for _, g := range strings.Split(glob, ",") {
			g = strings.TrimSpace(g)
			if g != "" {
				cmdArgs = append(cmdArgs, "-g", g)
			}
		}
	}

	cmdArgs = append(cmdArgs, pattern, searchPath)

	cmd := exec.Command(rgPath, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// rg returns exit code 1 when no matches, 2 for errors
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("No matches found for pattern '%s' in %s", pattern, searchPath)},
				},
			}, nil
		}
		return nil, fmt.Errorf("rg failed: %v\n%s", err, string(output))
	}

	// Parse JSON lines output
	type RgMatch struct {
		Type string `json:"type"`
		Data struct {
			Path struct {
				Text string `json:"text"`
			} `json:"path"`
			Line        json.Number `json:"line_number"`
			SubMatches  []struct {
				Match struct { Text string } `json:"match"`
			} `json:"submatches"`
			Lines struct {
				Text string `json:"text"`
			} `json:"lines"`
		} `json:"data"`
	}

	var results []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var match RgMatch
		if err := json.Unmarshal([]byte(line), &match); err != nil {
			continue
		}
		if match.Type == "match" {
			lineNum := match.Data.Line.String()
			filePath := match.Data.Path.Text
			content := match.Data.Lines.Text
			results = append(results, fmt.Sprintf("%s:%s:%s", filePath, lineNum, strings.TrimSpace(content)))
		}
	}

	if len(results) == 0 {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": fmt.Sprintf("No matches found for pattern '%s' in %s", pattern, searchPath)},
			},
		}, nil
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": fmt.Sprintf("Found %d matches for '%s':\n\n%s", len(results), pattern, strings.Join(results, "\n"))},
		},
	}, nil
}

// ---- fzf Find Files ----

func execFzfFindFiles(args map[string]interface{}) (map[string]interface{}, error) {
	searchPath := "."
	if p, ok := args["path"].(string); ok && p != "" {
		searchPath = p
	}

	query := ""
	if q, ok := args["query"].(string); ok {
		query = q
	}

	maxResults := 50
	if mr, ok := args["max_results"].(float64); ok && mr > 0 {
		maxResults = int(mr)
	}

	includeHidden, _ := args["include_hidden"].(bool)

	// Use find or fd to list files, pipe to fzf
	fzfPath := getFzfPath()

	// Build file listing command
	var listCmd *exec.Cmd
	if includeHidden {
		listCmd = exec.Command("find", searchPath, "-type", "f")
	} else {
		listCmd = exec.Command("find", searchPath, "-type", "f", "-not", "-path", "*/.*")
	}

	// Build fzf command
	fzfArgs := []string{"--filter", query, fmt.Sprintf("--limit=%d", maxResults), "--scheme=path"}
	if query == "" {
		// Without query, just take first N lines
		fzfArgs = []string{fmt.Sprintf("--limit=%d", maxResults), "--scheme=path"}
	}

	fzfCmd := exec.Command(fzfPath, fzfArgs...)

	// Create pipe between find and fzf
	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("pipe failed: %v", err)
	}
	listCmd.Stdout = pipeWriter
	fzfCmd.Stdin = pipeReader

	var stdout bytes.Buffer
	fzfCmd.Stdout = &stdout

	// Run pipeline
	if err := listCmd.Start(); err != nil {
		pipeReader.Close()
		pipeWriter.Close()
		return nil, fmt.Errorf("find failed: %v", err)
	}
	pipeWriter.Close()
	if err := fzfCmd.Start(); err != nil {
		pipeReader.Close()
		return nil, fmt.Errorf("fzf failed: %v", err)
	}
	if err := listCmd.Wait(); err != nil {
		// find errors are non-fatal if fzf got some output
	}
	pipeReader.Close()
	fzfErr := fzfCmd.Wait()

	files := strings.TrimSpace(stdout.String())
	if files == "" {
		if fzfErr == nil {
			return map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("No files found in %s", searchPath)},
				},
			}, nil
		}
		return nil, fmt.Errorf("fzf failed: %v", fzfErr)
	}

	fileList := strings.Split(files, "\n")
	if len(fileList) > maxResults {
		fileList = fileList[:maxResults]
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": fmt.Sprintf("Found %d files%s:\n\n%s", len(fileList),
				func() string { if query != "" { return " matching '" + query + "'" }; return "" }(),
				strings.Join(fileList, "\n"))},
		},
	}, nil
}

// ---- fzf Filter Lines ----

func execFzfFilterLines(args map[string]interface{}) (map[string]interface{}, error) {
	text, _ := args["text"].(string)
	query, _ := args["query"].(string)

	maxResults := 50
	if mr, ok := args["max_results"].(float64); ok && mr > 0 {
		maxResults = int(mr)
	}

	exact, _ := args["exact"].(bool)

	fzfPath := getFzfPath()
	fzfArgs := []string{"--filter", query, fmt.Sprintf("--limit=%d", maxResults)}
	if exact {
		fzfArgs = append(fzfArgs, "--exact")
	}

	cmd := exec.Command(fzfPath, fzfArgs...)
	cmd.Stdin = strings.NewReader(text)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// fzf returns 1 when no match
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": fmt.Sprintf("No lines matching '%s'", query)},
			},
		}, nil
	}

	lines := strings.TrimSpace(string(output))
	if lines == "" {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": fmt.Sprintf("No lines matching '%s'", query)},
			},
		}, nil
	}

	matchedLines := strings.Split(lines, "\n")
	if len(matchedLines) > maxResults {
		matchedLines = matchedLines[:maxResults]
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": fmt.Sprintf("Found %d lines matching '%s':\n\n%s", len(matchedLines), query, strings.Join(matchedLines, "\n"))},
		},
	}, nil
}

// ---- fzf Git Browse ----

func execFzfGitBrowse(args map[string]interface{}) (map[string]interface{}, error) {
	objType, _ := args["type"].(string)
	if objType == "" {
		return nil, fmt.Errorf("type is required (branches, tags, commits, stashes, remotes, reflogs, worktrees, files)")
	}

	searchPath := "."
	if p, ok := args["path"].(string); ok && p != "" {
		searchPath = p
	}

	query := ""
	if q, ok := args["query"].(string); ok {
		query = q
	}

	maxResults := 50
	if mr, ok := args["max_results"].(float64); ok && mr > 0 {
		maxResults = int(mr)
	}

	fzfPath := getFzfPath()

	// Build git command based on type
	var gitCmd *exec.Cmd
	switch objType {
	case "branches":
		gitCmd = exec.Command("git", "-C", searchPath, "branch", "--list", "--all", "--format=%(refname:short) %(committerdate:short) %(subject)")
	case "tags":
		gitCmd = exec.Command("git", "-C", searchPath, "tag", "-l", "--format=%(refname:short) %(creatordate:short) %(subject)")
	case "commits":
		gitCmd = exec.Command("git", "-C", searchPath, "log", "--all", "--oneline", "--no-decorate", fmt.Sprintf("-n=%d", maxResults*2))
	case "stashes":
		gitCmd = exec.Command("git", "-C", searchPath, "stash", "list")
	case "remotes":
		gitCmd = exec.Command("git", "-C", searchPath, "remote", "-v")
	case "reflogs":
		gitCmd = exec.Command("git", "-C", searchPath, "reflog", "--format=%gd %gs", fmt.Sprintf("-n=%d", maxResults*2))
	case "worktrees":
		gitCmd = exec.Command("git", "-C", searchPath, "worktree", "list")
	case "files":
		gitCmd = exec.Command("git", "-C", searchPath, "ls-files")
	default:
		return nil, fmt.Errorf("unknown type: %s (use branches, tags, commits, stashes, remotes, reflogs, worktrees, files)", objType)
	}

	var gitOutput bytes.Buffer
	gitCmd.Stdout = &gitOutput
	if err := gitCmd.Run(); err != nil {
		return nil, fmt.Errorf("git command failed: %v", err)
	}

	gitText := strings.TrimSpace(gitOutput.String())
	if gitText == "" {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": fmt.Sprintf("No %s found in %s", objType, searchPath)},
			},
		}, nil
	}

	// Apply fzf filtering
	if query != "" {
		fzfArgs := []string{"--filter", query, fmt.Sprintf("--limit=%d", maxResults)}
		fzfCmd := exec.Command(fzfPath, fzfArgs...)
		fzfCmd.Stdin = strings.NewReader(gitText)
		fzfOutput, err := fzfCmd.CombinedOutput()
		if err != nil {
			return map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("No %s matching '%s'", objType, query)},
				},
			}, nil
		}
		gitText = strings.TrimSpace(string(fzfOutput))
	}

	if gitText == "" {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": fmt.Sprintf("No %s matching '%s'", objType, query)},
			},
		}, nil
	}

	lines := strings.Split(gitText, "\n")
	if len(lines) > maxResults {
		lines = lines[:maxResults]
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": fmt.Sprintf("Found %d %s:\n\n%s", len(lines), objType, strings.Join(lines, "\n"))},
		},
	}, nil
}

// ---- Helpers ----

func getRgPath() string {
	if env := os.Getenv("RG_PATH"); env != "" {
		return env
	}
	if path, err := exec.LookPath("rg"); err == nil {
		return path
	}
	return "rg"
}

func getFzfPath() string {
	if env := os.Getenv("FZF_PATH"); env != "" {
		return env
	}
	if path, err := exec.LookPath("fzf"); err == nil {
		return path
	}
	return "fzf"
}

// Suppress unused import warnings
var _ = strconv.Itoa
var _ = filepath.Join

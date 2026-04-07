# Hacker News Extension

Fetch and browse Hacker News top stories. This extension provides MCP tools for the LLM to fetch HN content with configurable limits.

## Available Tools

### `hn_top_stories`
Fetch the top stories from Hacker News.

**Parameters:**
- `limit` (number, optional): Number of stories to return (1-30, default: 10)

**Returns:** Formatted list of stories with title, URL, score, author, and comment count.

### `hn_story_detail`
Get details for a specific Hacker News story by ID.

**Parameters:**
- `id` (number, required): Hacker News story ID

**Returns:** Detailed story information including title, URL, score, author, text, and comments.

## Example Usage

```
User: What are the top 5 stories on Hacker News right now?

Assistant: I'll fetch the top 5 HN stories for you.

[Uses hn_top_stories tool with limit=5]

Here are the top 5 stories:
1. Story Title
   URL: https://...
   Score: 100 | Author: user123 | Comments: 50
...
```

## When to Use

- User asks about Hacker News, HN, or Y Combinator news
- User wants to see trending tech stories
- User asks for specific HN story details by ID
- User wants to browse tech news or startup updates

## Notes

- Default limit is 10 stories, maximum is 30
- Story IDs can be used with `hn_story_detail` to fetch individual stories
- Stories are fetched from the public Hacker News Firebase API

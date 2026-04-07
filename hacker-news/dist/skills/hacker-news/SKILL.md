---
name: hacker-news
description: Fetch and browse Hacker News top stories, get story details, and stay updated with tech news and startup updates
---

# Hacker News Extension

## Overview

The Hacker News extension provides tools to fetch and browse content from Hacker News (news.ycombinator.com). Use these tools when users ask about tech news, trending stories, Y Combinator updates, or specific HN story details.

## Available Tools

### `hn_top_stories`

Fetch the top stories from Hacker News.

**Parameters:**
- `limit` (number, optional): Number of stories to return (1-30, default: 10)

**When to use:**
- User asks "What's on HN?"
- User wants to see trending tech stories
- User asks for tech news or startup updates
- User wants to browse Y Combinator content

**Example Usage:**
```
User: "What are the top stories on Hacker News?"
→ Call hn_top_stories with limit: 10

User: "Show me top 5 HN stories"
→ Call hn_top_stories with limit: 5
```

**Response Format:**
Returns a formatted list with:
- Story title
- URL (if available)
- Score (upvotes)
- Author
- Comment count

### `hn_story_detail`

Get details for a specific Hacker News story by ID.

**Parameters:**
- `id` (number, required): Hacker News story ID

**When to use:**
- User provides a specific HN story ID
- User asks for details about a particular story
- User wants to read comments or full story content

**Example Usage:**
```
User: "Tell me about story 12345"
→ Call hn_story_detail with id: 12345

User: "What's the details of that HN post?"
→ Ask for the story ID, then call hn_story_detail
```

**Response Format:**
Returns detailed story information:
- Title
- URL
- Score
- Author
- Text content (if available)
- Comment count

## Best Practices

1. **Default Limits**: When user doesn't specify, use 10 stories
2. **Maximum**: Don't exceed 30 stories in one request
3. **Formatting**: Present stories in a clean, numbered list
4. **Context**: Mention that stories are from Hacker News
5. **Follow-up**: Offer to get details for any story the user is interested in

## Example Interactions

### Browsing Top Stories

```
User: "What's trending on HN today?"

Assistant: I'll fetch the top Hacker News stories for you.
[Calls hn_top_stories with limit: 10]

Here are the top 10 stories on Hacker News:

1. **Story Title One**
   URL: https://example.com/story1
   Score: 150 | Author: user123 | Comments: 45

2. **Another Tech Story**
   URL: https://techblog.com/post
   Score: 120 | Author: dev456 | Comments: 32

...

Would you like more details on any of these stories?
```

### Getting Story Details

```
User: "Tell me more about story 38472615"

Assistant: Let me fetch the details for that story.
[Calls hn_story_detail with id: 38472615]

**Story Title**
URL: https://example.com/article

Score: 200
Author: techwriter
Comments: 67

Story content and discussion summary...
```

## Limitations

- Stories are fetched from the public HN Firebase API
- No authentication or private content access
- Maximum 30 stories per request
- Story IDs must be valid HN item IDs
- Text content may be truncated for long stories

## Related Commands

- `/hn-top` - Fetch top stories (custom command)
- `/hn-story` - Get story details (custom command)

## Triggers

Use these tools when you see:
- "Hacker News" or "HN"
- "tech news" or "technology news"
- "trending stories" or "top stories"
- "Y Combinator" or "YC news"
- Specific HN story IDs or URLs

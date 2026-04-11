---
name: oido-cron
description: Create, manage, and monitor scheduled cron jobs for automated agent tasks in OIDO Studio
---

# OIDO Cron Manager Extension

## Overview

The OIDO Cron Manager extension enables AI agents to create and manage scheduled cron jobs. These are automated tasks that run on a schedule — they can perform periodic checks, generate reports, process data, or trigger agent turns without user interaction.

Use this extension when users want to:
- Automate recurring tasks
- Set up scheduled reports
- Create periodic checks
- Schedule automated agent turns
- Manage existing scheduled jobs

## Required Fields for Creating Jobs

When using `cron_add`, you MUST provide:

| Field | Required | Description |
|-------|----------|-------------|
| **name** | ✅ YES | Job name (any descriptive name) |
| **message** | ✅ YES | The prompt/instruction for the agent to execute |
| **schedule** | ✅ YES | When to run (see formats below) |

All other fields are optional and have sensible defaults.

## Schedule Formats

The `schedule` field accepts THREE formats:

### 1. Cron Expression (Recurring)
```
schedule: "0 9 * * *"        # Daily at 9 AM UTC
schedule: "0 9 * * 1-5"      # Weekdays at 9 AM UTC
schedule: "*/5 * * * *"      # Every 5 minutes
schedule: "0 * * * *"        # Every hour
```

With timezone:
```
schedule: "0 9 * * *|America/New_York"
```

### 2. Fixed Interval (Recurring)
```
schedule: "every|3600000"    # Every hour (3600000ms)
schedule: "every|300000"     # Every 5 minutes
schedule: "every|86400000"   # Every day
```

### 3. One-Shot (Single Execution)
```
schedule: "at|30m"           # In 30 minutes
schedule: "at|2025-12-31T23:59:59Z"  # At specific time
```

## Available Tools

### `cron_add` - Create a New Job

**Required Parameters:**
- `name` (string): Job name
- `message` (string): Prompt for the agent
- `schedule` (string): When to run (use formats above)

**Optional Parameters:**
- `tz` (string): Timezone for cron expressions (default: UTC)
- `delivery` (string): "none", "session", or "channel" (default: none)
- `channel` (string): Channel type like "whatsapp" (only if delivery=channel)
- `to` (string): Recipient like "+1234567890" (only if delivery=channel)
- `model` (string): Model override
- `session` (string): "isolated" or "main" (default: isolated)

**Example - Daily Report:**
```json
{
  "name": "Daily Morning Report",
  "message": "Generate a summary of yesterday's activities",
  "schedule": "0 9 * * *",
  "tz": "America/New_York",
  "delivery": "session"
}
```

**Example - Alert Check Every 5 Minutes:**
```json
{
  "name": "Error Monitor",
  "message": "Check system logs for errors",
  "schedule": "every|300000",
  "delivery": "channel",
  "channel": "whatsapp",
  "to": "+1234567890"
}
```

**Example - One-Shot Reminder:**
```json
{
  "name": "Meeting Reminder",
  "message": "Remind user about upcoming meeting",
  "schedule": "at|30m"
}
```

### `cron_list` - List All Jobs

**Parameters:** None

**When to use:**
- User asks "What cron jobs do I have?"
- Show overview of scheduled tasks

### `cron_get` - Get Job Details

**Parameters:**
- `id` (number, required): Job ID

### `cron_toggle` - Enable/Disable Job

**Parameters:**
- `id` (number, required): Job ID
- `enabled` (boolean, required): true to enable, false to disable

### `cron_run` - Run Job Immediately

**Parameters:**
- `id` (number, required): Job ID

### `cron_logs` - View Execution Logs

**Parameters:**
- `id` (number, required): Job ID
- `limit` (number, optional): Number of entries (default: 10)

### `cron_update` - Update Job

**Parameters:**
- `id` (number, required): Job ID
- `name` (string, optional): New name
- `schedule` (string, optional): New schedule
- `message` (string, optional): New prompt
- `delivery` (string, optional): New delivery mode

### `cron_delete` - Delete Job

**Parameters:**
- `id` (number, required): Job ID

⚠️ **WARNING:** This action cannot be undone!

## Common Mistakes to Avoid

### ❌ WRONG - Missing required fields
```json
{
  "name": "My Job"
  // ERROR: missing 'message' and 'schedule'
}
```

### ✅ CORRECT - All required fields
```json
{
  "name": "My Job",
  "message": "Do something useful",
  "schedule": "0 9 * * *"
}
```

### ❌ WRONG - Conflicting schedule fields
The old API had separate `interval_ms` and `at` fields. **DON'T use them!**
Use the single `schedule` field with the formats shown above.

### ✅ CORRECT - Single schedule field
```json
{
  "schedule": "every|3600000"
}
```

## Delivery Modes Explained

| Mode | Where Results Go | Use Case |
|------|------------------|----------|
| `none` | Job logs only | Background tasks |
| `session` | Agent session | User sees results |
| `channel` | WhatsApp/Telegram | Notifications |

## Example Workflows

### 1. Daily News Digest
```
User: "Create a daily HN digest at 8 AM"

Assistant: I'll create a cron job for that.
→ cron_add({
    name: "HN Daily Digest",
    message: "Fetch top 20 HN stories and summarize",
    schedule: "0 8 * * *",
    delivery: "session"
  })
✓ Created job #123
```

### 2. Error Monitoring
```
User: "Alert me on WhatsApp if there are errors every 5 min"

Assistant: I'll set up error monitoring.
→ cron_add({
    name: "Error Monitor",
    message: "Check logs for critical errors",
    schedule: "every|300000",
    delivery: "channel",
    channel: "whatsapp",
    to: "+1234567890"
  })
✓ Created job #456
```

### 3. Check and Manage Jobs
```
User: "Show me my jobs"
→ cron_list()

User: "Disable job 123"
→ cron_toggle({ id: 123, enabled: false })

User: "Run job 456 now"
→ cron_run({ id: 456 })
```

## Triggers

Use these tools when you see:
- "schedule" or "scheduled task"
- "cron job" or "cron"
- "recurring" or "every day/hour"
- "automate" or "automation"
- "remind me at..."
- "periodic check"
- "monitor every..."

# OIDO Cron Manager Extension

## Overview

This extension provides tools to create, manage, and monitor scheduled cron jobs in OIDO Studio. Cron jobs are automated tasks that run on a schedule — they can perform periodic checks, generate reports, process data, or trigger agent turns without user interaction.

## Available Tools

### `cron_list`
List all scheduled cron jobs. Returns job ID, name, schedule, status, and delivery mode.

### `cron_get`
Get details of a specific cron job by its ID.

### `cron_add`
Create a new scheduled cron job.
- **schedule**: Cron expression (e.g., '0 9 * * 1-5')
- **interval_ms**: Fixed interval in milliseconds
- **at**: One-shot time (e.g., '30m' or RFC3339)
- **delivery**: 'none' (default), 'session', or 'channel'

### `cron_toggle`
Enable or disable a cron job. Use to pause without deleting.

### `cron_run`
Run a cron job immediately, regardless of schedule.

### `cron_logs`
View execution logs for a specific cron job.

### `cron_update`
Update an existing cron job's properties.

### `cron_delete`
Permanently delete a cron job. Cannot be undone.

## Common Cron Expressions

| Expression | Description |
|------------|-------------|
| `*/5 * * * *` | Every 5 minutes |
| `0 * * * *` | Every hour |
| `0 9 * * *` | Daily at 9 AM UTC |
| `0 9 * * 1-5` | Weekdays at 9 AM UTC |
| `0 0 * * 0` | Weekly on Sunday |
| `0 0 1 * *` | Monthly on 1st |

## Delivery Modes

- **none**: Results stay in job logs
- **session**: Results sent to agent session
- **channel**: Results sent via WhatsApp, Telegram, etc.

## Example Usage

```
User: "Create a daily report at 9 AM"
→ Use cron_add with schedule="0 9 * * *", message="Generate daily summary..."

User: "Show me my cron jobs"
→ Use cron_list

User: "Run job 123 now"
→ Use cron_run with id=123
```

## Environment Variables

- `OIDO_API_BASE`: API base URL (default: http://localhost:8080)
- `OIDO_API_TOKEN`: Authentication token

## When to Use

- User wants to automate recurring tasks
- Set up scheduled reports or checks
- Create periodic agent turns
- Manage existing scheduled jobs
- Monitor job execution history

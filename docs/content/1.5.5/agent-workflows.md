# Agent Workflows

Andurel's CLI exposes stable structured output for agents and automation.

## Discover the project

Start with read-only commands:

```bash
andurel --agent --help
andurel commands --json
andurel project info --json
andurel routes --json
andurel models --json
andurel migrations --json
andurel controllers --json
andurel views --json
andurel jobs --json
```

## Output modes

Use `--json` for the response envelope, `--agent` for structured agent output, `--md` for Markdown where supported, and `--quiet` to suppress non-essential progress. Use `--jq`, `--ids-only`, or `--count` when a command supports a smaller projection.

## Safe mutations

Preview writes before applying them:

```bash
andurel generate scaffold Product --dry-run --diff --json
andurel upgrade --dry-run --diff --json
```

Structured mutation reports identify created, updated, and deleted files, routes, commands, warnings, and useful next commands.

## Install the skill

```bash
andurel skill show
andurel skill install --harness pi
```

Automation must pass one or more harness values so the command never waits for interactive selection.

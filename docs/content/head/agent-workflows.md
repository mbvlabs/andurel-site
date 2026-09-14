# Agent Workflows

Andurel's CLI exposes structured discovery and mutation reports for agents and automation.

## Discover before acting

```bash
andurel --agent --help
andurel commands --json
andurel project info --json
andurel config show --json
andurel routes --json
andurel models --json
andurel migrations --json
andurel controllers --json
andurel views --json
andurel jobs --json
```

Use the returned command tree and project metadata instead of assuming a v1 command or directory still exists.

## Output and projections

Use `--json` for the response envelope, `--agent` for structured agent output, `--md` for Markdown where supported, and `--quiet` to suppress non-essential progress. `--jq`, `--ids-only`, and `--count` provide smaller projections where supported.

## Preview safe mutations

```bash
andurel new orbit --dry-run --diff --json
andurel generate model Product --dry-run --diff --json
andurel generate query ProductReport --dry-run --diff --json
andurel generate scaffold Product --dry-run --diff --json
andurel extension add docker --dry-run --diff --json
```

Review created, updated, and deleted files, route additions, commands, warnings, and breadcrumbs before applying changes. A v2 migration is intentionally manual; do not point `andurel upgrade` at a v1 project to cross the major version.

## Install the embedded skill

```bash
andurel skill show --json
andurel skill install --harness codex,claude
```

Automation must pass at least one harness so installation never waits for an interactive selection.

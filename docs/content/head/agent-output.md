# Agent Output

Andurel's CLI exposes structured discovery and mutation reports for agents and automation.

## Discover before acting

```bash
andurel commands --json
andurel inspect project --json
andurel inspect routes --json
andurel inspect models --json
andurel inspect migrations --json
andurel inspect controllers --json
andurel inspect views --json
andurel inspect jobs --json
andurel doctor --json
```

Use the returned command tree and project metadata instead of assuming a v1 command or directory still exists.

## Output and projections

| Flag | Purpose |
| --- | --- |
| `--json` | JSON response envelope |
| `--agent` | Structured agent-oriented output |
| `--md` | Markdown where supported |
| `--quiet` | Suppress non-essential progress |
| `--jq` | Select a simple field path from command data |
| `--ids-only` | One resource identifier per line when supported |
| `--count` | One raw resource count when supported |

## Preview safe mutations

```bash
andurel new orbit --dry-run --diff --json
andurel generate model Product --dry-run --diff --json
andurel generate query ProductReport --dry-run --diff --json
andurel generate scaffold Product --dry-run --diff --json
```

Review created, updated, and deleted files, route additions, commands, warnings, and breadcrumbs before applying changes. A v2 migration is intentionally manual; do not point `andurel upgrade` at a v1 project to cross the major version.

## Install the embedded skill

```bash
andurel skill show --json
andurel skill install --harness codex,claude
```

Automation must pass at least one harness so installation never waits for an interactive selection. See [skill](/docs/head/skill).

# Agentic Development

Andurel is built so humans and agents share one project pad: generated application code you own, an agent-ready CLI, and an installable skill.

## Install the skill

```bash
andurel skill show
andurel skill install --harness claude,codex
```

`skill install` writes the embedded Andurel skill into the chosen harnesses (`claude`, `codex`, `pi`, `opencode`, `crush`). Prefer the skill over inventing paths, URLs, or generator names. Automation must pass `--harness` so installation never waits for interactive selection.

## Project guidance

Generated applications include `AGENTS.md`. Agents should:

1. Run `andurel commands --json` before guessing subcommands
2. Prefer `andurel inspect … --json` for read-only discovery
3. Use `--dry-run` / structured output on mutating generators when available
4. Import routes from `@/routes` after `andurel sync routes` in Inertia apps

```bash
andurel commands --json
andurel inspect project --json
andurel doctor --json
andurel generate scaffold Product --dry-run --json
andurel sync queries --json
andurel inspect routes --json
```

## Structured CLI output

Global flags make command results machine-readable:

| Flag | Purpose |
| --- | --- |
| `--json` | JSON response envelope |
| `--md` | Markdown where supported |
| `--agent` | Structured agent-oriented output |
| `--jq` | Select a simple field path from command data |
| `--ids-only` | One resource identifier per line when supported |
| `--count` | One raw resource count when supported |
| `--quiet` | Suppress non-essential human progress |

See [Agent Output](/docs/head/agent-output) for mutation previews and [Overview](/docs/head/cli) for the full command surface.

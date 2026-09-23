# Agentic Development

Andurel is built so humans and agents share one project pad: generated application code you own, an agent-ready CLI, and an installable skill.

## Install the skill

```bash
andurel skill install
andurel skill show
```

`skill install` writes the embedded Andurel skill into the agent harness configured for the project. Prefer the skill over inventing paths, URLs, or generator names.

## Project guidance

Generated applications include `AGENTS.md` (or restore it with the CLI when needed). Agents should:

1. Run `andurel commands --json` before guessing subcommands
2. Prefer `andurel inspect … --json` for read-only discovery
3. Use `--dry-run` / structured output on mutating generators when available
4. Import routes from `@/routes` after `andurel sync routes` in Inertia apps

## Structured CLI output

Global flags such as `--json`, `--md`, `--agent`, `--jq`, `--ids-only`, and `--count` make command results machine-readable. See [Agent Output](/docs/head/agent-output) for patterns and [Overview](/docs/head/cli) for the full command surface.

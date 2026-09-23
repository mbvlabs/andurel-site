# skill

`andurel skill` shows or installs the embedded Andurel agent skill.

```bash
andurel skill show
andurel skill show --json
andurel skill install --harness claude,codex
andurel skill install --harness pi --harness opencode
```

Supported harnesses: `claude`, `codex`, `pi`, `opencode`, `crush`. Install once per harness so agents discover commands, generators, and project layout without inventing paths. Automation must pass `--harness` so installation never waits for interactive selection.

## Show

```bash
andurel skill show
```

Prints the skill markdown (front matter plus body):

```text
---
name: andurel
description: Use this skill for Andurel framework projects when deciding where code belongs...
---

# Andurel

Use this skill when working in an Andurel project or generating Andurel code.
...
```

Structured form:

```bash
andurel skill show --json
```

```json
{
  "ok": true,
  "data": {
    "name": "andurel",
    "body": "---\nname: andurel\n..."
  },
  "summary": "Loaded Andurel skill"
}
```

## Install

```bash
andurel skill install --harness claude,codex --json
```

```json
{
  "ok": true,
  "data": {
    "name": "andurel",
    "installations": [
      {
        "harness": "claude",
        "path": "/home/you/orbit/.claude/skills/andurel/SKILL.md"
      },
      {
        "harness": "codex",
        "path": "/home/you/orbit/.codex/skills/andurel/SKILL.md"
      }
    ]
  },
  "summary": "Installed Andurel skill for Claude, Codex"
}
```

A single harness also sets `data.path` to that `SKILL.md`. Re-run install after CLI upgrades to refresh the embedded skill text.

See [Agentic Development](/docs/head/agentic-development) and [Agent Output](/docs/head/agent-output).

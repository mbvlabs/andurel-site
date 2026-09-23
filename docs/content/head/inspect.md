# inspect

`andurel inspect` is read-only project shape discovery. Use it before generating or syncing so agents and humans see the same pad.

## Commands

```bash
andurel inspect project --json
andurel inspect routes --json
andurel inspect models --json
andurel inspect migrations --json
andurel inspect controllers --json
andurel inspect views --json
andurel inspect jobs --json
```

Prefer `--json` (or `--md`) for automation. Mutating work belongs to [generate](/docs/head/generate), [sync](/docs/head/sync), and [db](/docs/head/db).

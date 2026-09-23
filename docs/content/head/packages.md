# packages

`andurel packages` lists or updates Andurel `pkg/*` module versions in `go.mod`.

```bash
andurel packages list --json
andurel packages update
andurel packages update storage inertia kiks
andurel packages update --dry-run --json
```

## List

```bash
andurel packages list --json
```

```json
{
  "ok": true,
  "data": {
    "action": "list",
    "packages": [
      {
        "name": "email",
        "path": "github.com/mbvlabs/andurel/pkg/email",
        "current": "v0.3.3",
        "latest": "v0.3.3",
        "status": "current"
      },
      {
        "name": "inertia",
        "path": "github.com/mbvlabs/andurel/pkg/inertia",
        "current": "v0.5.0",
        "latest": "v0.6.0",
        "status": "outdated"
      },
      {
        "name": "storage",
        "path": "github.com/mbvlabs/andurel/pkg/storage",
        "current": "v0.8.0",
        "latest": "v0.8.0",
        "status": "current"
      }
    ]
  },
  "summary": "1 Andurel package has an update available"
}
```

Human lines look like:

```text
  email     v0.3.3          current
  inertia   v0.5.0 -> v0.6.0  outdated
  storage   v0.8.0          current
```

## Update

```bash
andurel packages update --dry-run --json
andurel packages update --json
```

Dry-run keeps `status: outdated` without rewriting `go.mod`. A successful update flips matching entries to `status: updated` and prints a summary such as `Updated 1 Andurel package`.

Updating the CLI does not rewrite application imports. `andurel.toml` / `andurel.lock` pin tool binaries, not Go modules. Read each package changelog, compile the owning process, and run tests after bumps.

See [Framework Packages](/docs/head/framework-packages) and [Project Lock](/docs/head/project-lock).

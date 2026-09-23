# packages

`andurel packages` lists or updates Andurel `pkg/*` module versions in `go.mod`.

```bash
andurel packages list --json
andurel packages update
andurel packages update storage inertia
```

Updating the CLI does not rewrite application imports. Read each package changelog, compile the owning process, and run tests after bumps. See [Framework Packages](/docs/head/framework-packages) and [Project Lock](/docs/head/project-lock).

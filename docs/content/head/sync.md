# sync

`andurel sync` refreshes **derived** files from a source of truth. It does not create application-owned domain files; that is [generate](/docs/head/generate).

## When to sync versus generate

| Need | Command |
| --- | --- |
| New model, controller, job, email, scaffold, migration | `andurel generate ...` |
| Refresh `*_templ.go`, narsilc clients, routes.ts, payloads, email CSS, factories | `andurel sync ...` |
| Drift check without writing | `andurel sync factory NAME --check` / `sync factories --check` |

Generate creates or updates owned Go/SQL/frontend sources. Sync recompiles or regenerates artifacts those sources imply. After editing `models/queries/*.sql`, run `sync queries`. After adding an Inertia route marker, run `sync routes`. After changing payload structs, run `sync payloads`.

`andurel tool sync` downloads pinned binaries from `andurel.toml` / `andurel.lock` into `bin/`. It is not this group.

## Commands

```bash
andurel sync views
andurel sync queries
andurel sync routes --json
andurel sync payloads --json
andurel sync email
andurel sync factory User --check
andurel sync factory User --sync
andurel sync factories --check --json
andurel sync factories --sync
```

| Command | Source → output |
| --- | --- |
| `views` | `.templ` → `*_templ.go` (and email compile when inputs exist) |
| `queries` | `models/queries/*.sql` → `models/internal/queries` Go clients |
| `routes` | route declarations → `resources/js/routes.ts` |
| `payloads` | controller payload structs → `resources/js/types/payloads.ts` |
| `email` | authored email templates → inlined Tailwind / renderers |
| `factory` | one model Entity → `models/factories` declaration |
| `factories` | every managed Entity → factories (`--check` or `--sync` required) |

## Sample output

Routes:

```bash
andurel sync routes --json
```

```json
{
  "ok": true,
  "data": {
    "generated_file": "resources/js/routes.ts",
    "generated_helpers": 13,
    "skipped_count": 4
  },
  "summary": "Generated 13 route helpers to resources/js/routes.ts (4 skipped)"
}
```

Payloads:

```bash
andurel sync payloads --json
```

```json
{
  "ok": true,
  "data": {
    "generated_file": "resources/js/types/payloads.ts",
    "type_count": 5,
    "skipped_count": 0
  },
  "summary": "Generated 5 payload types to resources/js/types/payloads.ts"
}
```

Queries:

```bash
andurel sync queries --json
```

```json
{
  "ok": true,
  "data": {
    "action": "generate queries",
    "files_created": [],
    "files_updated": [],
    "routes_added": [],
    "commands_run": [
      "narsilc generate",
      "go fmt ./models/internal/queries/..."
    ]
  },
  "summary": "generate queries completed with no file changes"
}
```

When SQL changed, `files_updated` lists regenerated clients under `models/internal/queries/`.

Factory check:

```bash
andurel sync factory Product --check --json
```

Pass `--sync` to rewrite managed factory regions when the Entity drifted. Prefer `--json` so agents can read `ok`, `summary`, and any error `hint` without scraping human text.

See [TypeScript Sync](/docs/head/inertia-typescript-sync) for the Inertia-focused route and payload workflow, and [Factories](/docs/head/factories) for factory drift checks.

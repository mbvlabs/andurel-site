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

| Command | Reports |
| --- | --- |
| `project` | Manifest metadata, UI choice, database settings |
| `routes` | Route manifest entries |
| `models` | Model files under `models/` |
| `migrations` | SQL migrations under `migrations/` |
| `controllers` | Controller files |
| `views` | Templ view files |
| `jobs` | Job and worker files |

Prefer `--json` (or `--md`) for automation. Combine with `--jq`, `--ids-only`, or `--count` when you only need identifiers or a count.

## Sample: project

```bash
andurel inspect project --json
```

```json
{
  "ok": true,
  "data": {
    "root": "/home/you/orbit",
    "module": "orbit",
    "go_version": "1.27.1",
    "andurel_version": "v2.0.0-alpha",
    "scaffold_config": {
      "projectName": "orbit",
      "inertia": "react",
      "javascriptPackageManager": "pnpm",
      "javascriptSSRRuntime": "node"
    },
    "database_config": {
      "engine": "postgresql",
      "nullType": "pgtype.Null"
    }
  },
  "summary": "Project info loaded"
}
```

## Sample: routes

```bash
andurel inspect routes --json
```

```json
{
  "ok": true,
  "data": {
    "routes": [
      {
        "variable": "SessionNew",
        "name": "users.new_user_session",
        "path": "/users/sign-in",
        "constructor": "NewSimpleRoute",
        "kind": "simple",
        "is_inertia": true,
        "source_file": "router/routes/users.go",
        "line": 9
      },
      {
        "variable": "DocumentationShow",
        "name": "documentations.show",
        "path": "/docs/:version/:slug",
        "constructor": "NewRouteWithParams",
        "kind": "params",
        "is_inertia": true,
        "params": [
          { "name": "version", "type": "string" },
          { "name": "slug", "type": "string" }
        ],
        "source_file": "router/routes/documentations.go",
        "line": 30
      }
    ]
  },
  "summary": "Listed 20 routes (4 skipped)"
}
```

Mutating work belongs to [generate](/docs/head/generate), [sync](/docs/head/sync), and [db](/docs/head/db).

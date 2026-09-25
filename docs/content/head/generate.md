# generate

`andurel generate` creates **application-owned** code you edit after creation: migrations, models, controllers, scaffolds, queries, jobs, and email templates. Framework behavior stays in versioned packages; project policy and domain code remain local.

Use the child pages for each subcommand's flags and file outputs. This page covers the group contract: when to generate versus sync, shared preview flags, and how to discover the live CLI surface.

## Start here

| Command | Use it when |
| --- | --- |
| [scaffold](/docs/head/generate-scaffold) | Add a full resource from an existing table migration |
| [migration](/docs/head/generate-migration) | Add a schema change as SQL under `migrations/` |
| [model](/docs/head/generate-model) | Create or update a model API from migration history |
| [controller](/docs/head/generate-controller) | Add HTTP handlers, views/pages, and routes for an existing model |
| [query](/docs/head/generate-query) | Author a new narsilc SQL file under `models/queries/` |
| [job](/docs/head/generate-job) | Add a River job plus worker registration |
| [email](/docs/head/generate-email) | Author a transactional or marketing email template |

## Generate versus sync

| Need | Command |
| --- | --- |
| New or updated owned Go/SQL/frontend sources | `andurel generate ...` |
| Refresh derived artifacts those sources imply | `andurel sync ...` |

Generate writes migrations, model packages, controllers, route files, query SQL, job args, and email templates. Sync recompiles Templ, narsilc clients, TypeScript routes/payloads, email CSS, and factories. See [sync](/docs/head/sync).

## Shared flags

Most generate subcommands accept the same preview and agent flags:

| Flag | What it does |
| --- | --- |
| `--dry-run` | Preview file changes without applying |
| `--diff` | Include a text diff preview in structured output |
| `--json` | Emit the structured `{ok,data,summary,breadcrumbs}` envelope |
| `--md` | Emit Markdown where supported |
| `--agent` | Structured output with non-essential human progress suppressed |
| `--jq` | Select a simple field path from command data |
| `--quiet` / `--verbose` | Suppress or expand human-only output |

Prefer `--dry-run --diff --json` before broad mutations. Child pages document command-local flags only.

## Discover the live surface

The generate group evolves with the CLI. Prefer discovery over copying a static flag list:

```bash
andurel commands --json
andurel generate --help
andurel generate scaffold --help
```

## Typical resource workflow

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel generate scaffold Product --dry-run --diff --json
andurel generate scaffold Product
```

Generation follows the UI recorded in `andurel.toml` / `andurel.lock` (Inertia pages for Vue, React, or Svelte; Templ otherwise). Pass `--api` on scaffold or controller for JSON under `controllers/api`. Inertia payload, TypeScript, and route-helper details live under [Generators](/docs/head/inertia-generators).

After Inertia routes or payload structs change, refresh derived TypeScript with `andurel sync routes` and `andurel sync payloads`. See [Agent Output](/docs/head/agent-output) for structured envelopes and dry-run conventions.

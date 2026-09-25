# scaffold

`andurel generate scaffold` creates a full resource from an existing table migration: constructed model API, entity and data types, optional factory, controller, typed routes, and views or Inertia pages.

## When to use

- Add a full resource when the matching migration already exists under `migrations/`.
- Prefer [model](/docs/head/generate-model) or [controller](/docs/head/generate-controller) when you only need one layer.

## Usage

```bash
andurel generate scaffold NAME [flags]
```

`NAME` is CamelCase. One lowercase namespace segment is allowed (`admin/Widget`). Namespaced scaffolds place controllers under `controllers/admin`, use `admin.*` route names, and Admin-prefixed symbols.

## Flags

| Flag | What it does |
| --- | --- |
| `--api` | Generate a JSON API controller under `controllers/api` instead of UI pages |
| `--primary-key` | Specify the primary key column (skips interactive detection) |
| `--table-name` | Override the default table name derived from `NAME` |
| `--skip-factory` | Skip generating a factory for the model |
| `--dry-run` / `--diff` | Preview without writing; see [generate](/docs/head/generate) for shared flags |

Without `--api`, generation follows the project UI in `andurel.lock` (Inertia adapter or Templ).

## Examples

```bash
andurel generate scaffold Product --dry-run --diff --json
andurel generate scaffold Product
andurel generate scaffold admin/Widget
andurel generate scaffold Product --api
andurel generate scaffold Product --table-name catalog_products --primary-key id
```

A typical scaffold writes:

- `models/product.go` and related query SQL
- `models/factories` entry (unless `--skip-factory`)
- `controllers/products.go` (or `controllers/api` / namespaced path)
- `router/routes/products.go`
- Templ views or Inertia pages for the configured adapter

## Next

```bash
andurel sync routes --json
andurel sync payloads --json
andurel doctor --json
```

See [Generators](/docs/head/inertia-generators) for Inertia payload and TypeScript flow, and [migration](/docs/head/generate-migration) when the table does not exist yet.

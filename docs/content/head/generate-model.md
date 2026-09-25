# model

`andurel generate model` creates or updates a model API from SQL migration history. It constructs the plural model type, entity and data structs, annotated query SQL, Fx module wiring, and an optional factory.

## When to use

- Create `models/NAME.go` from an existing table migration.
- Pass `--update` after a migration changes columns.
- Do not invent columns by hand; write a [migration](/docs/head/generate-migration) first.

## Usage

```bash
andurel generate model NAME [flags]
```

## Flags

| Flag | What it does |
| --- | --- |
| `--mode` | Operation mode: `crud` (default), `read-only`, or `create-only` |
| `--update` | Update an existing model from migration changes |
| `--yes` | Apply changes without prompting for confirmation |
| `--primary-key` | Specify the primary key column (skips interactive detection) |
| `--table-name` | Override the default table name |
| `--skip-factory` | Skip generating or updating the matching factory |
| `--dry-run` / `--diff` | Preview without writing; see [generate](/docs/head/generate) |

CRUD modes require a supported primary key. Model generation updates `models.Module` so Fx can supply the new plural model API.

## Examples

```bash
andurel generate model Product --dry-run --json
andurel generate model Product
andurel generate model Product --mode read-only
andurel generate model Product --update --yes
andurel generate model Product --table-name catalog_products --primary-key id
```

## Next

```bash
andurel sync factory Product --check --json
andurel sync factory Product --sync --json
andurel inspect models --json
```

For HTTP layers on top of an existing model, use [controller](/docs/head/generate-controller) or [scaffold](/docs/head/generate-scaffold). Factory sync details live under [Factories](/docs/head/factories) and [sync](/docs/head/sync).

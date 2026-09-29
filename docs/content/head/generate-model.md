# model

`andurel generate model` creates a model API from SQL migration history, or a custom query model from explicit field specs.

## When to use

- Create `models/NAME.go` from an existing table migration.
- Pass `--custom field:type ...` for a non-table model used with narsilc custom queries.
- After a migration changes columns on an existing table model, use `andurel sync model NAME` — not `generate model --update`.
- Do not invent table columns by hand; write a [migration](/docs/head/generate-migration) first.
- Do not use `--custom` when the entity maps to a real table.

## Usage

```bash
andurel generate model NAME [flags]
```

## Flags

| Flag | What it does |
| --- | --- |
| `--custom` | Generate a non-table-backed model from `field:type` specs |
| `--mode` | Operation mode: `crud` (default), `read-only`, or `create-only` |
| `--primary-key` | Specify the primary key column (skips interactive detection) |
| `--table-name` | Override the default table name |
| `--skip-factory` | Skip generating the matching factory |
| `--dry-run` / `--diff` | Preview without writing; see [generate](/docs/head/generate) |

CRUD modes require a supported primary key. Table model generation updates `models.Module` so Fx can supply the new plural model API.

## Custom models

`andurel generate model NAME --custom field:type ...` writes a struct-only model marked `// andurel:custom`. There is no table, no migration, and no factory. Do not run `andurel sync model` or `andurel sync factory` after `--custom`; those commands require a table-backed model. After authoring SQL, run `andurel sync queries`.

Supported `field:type` specs follow the generator's field type list (for example `id:uuid`, `name:string`, `total:int64`).

## Examples

```bash
andurel generate model Product --dry-run --json
andurel generate model Product
andurel generate model Product --mode read-only
andurel generate model Product --table-name catalog_products --primary-key id
andurel generate model AggregateResult --custom id:uuid name:string currency:int64
```

## Next

Table models:

```bash
andurel sync model Product
andurel sync factory Product --check --json
andurel inspect models --json
```

Custom models:

```bash
andurel sync queries --json
andurel doctor --json
```

For HTTP layers on top of an existing model, use [controller](/docs/head/generate-controller) or [scaffold](/docs/head/generate-scaffold). Factory sync details live under [Factories](/docs/head/factories) and [sync](/docs/head/sync).

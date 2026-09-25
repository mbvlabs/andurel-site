# query

`andurel generate query` creates a narsilc SQL query file under `models/queries/`. It does not compile SQL into Go clients; that is `andurel sync queries`.

## When to use

- Add a new annotated SQL file for typed persistence.
- Do not use this to regenerate clients after editing SQL; run [sync](/docs/head/sync) `queries`.

## Usage

```bash
andurel generate query NAME [flags]
```

`NAME` is CamelCase (for example `UserReport` or `SalesReport`).

## Flags

| Flag | What it does |
| --- | --- |
| `--table` | Existing table name for a starter annotated query |
| `--dry-run` / `--diff` | Preview without writing; see [generate](/docs/head/generate) |

## Examples

```bash
andurel generate query UserReport --dry-run --json
andurel generate query UserReport --table users
andurel generate query SalesReport --table products
```

Edit the authored SQL under `models/queries/`, then sync. Generated clients land in `models/internal/queries/`; only the owning `models` package should import that internal package.

## Next

```bash
andurel sync queries --json
```

See [Queries](/docs/head/queries) for narsilc conventions and package boundaries.

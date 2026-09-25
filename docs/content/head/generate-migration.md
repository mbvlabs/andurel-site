# migration

`andurel generate migration` creates a new SQL migration file under root `migrations/`. It does not apply the migration.

## When to use

- Add a schema change as Goose SQL before generating models or scaffolds.
- Do not use this to apply or roll back schema; that is [db](/docs/head/db) `migrate`.

## Usage

```bash
andurel generate migration NAME
```

`NAME` is a snake_case description such as `create_products_table` or `add_status_to_products`.

## Flags

Migration has no command-local flags beyond the shared agent and output flags documented on [generate](/docs/head/generate).

## Examples

```bash
andurel generate migration create_products_table
andurel generate migration add_status_to_products
```

The command writes a timestamped SQL file under `migrations/` with up and down sections for you to edit.

## Next

```bash
andurel db migrate up
andurel generate model Product
# or
andurel generate scaffold Product
```

See [Migrations & Seeding](/docs/head/migrations) for apply, rollback, and seed workflows.

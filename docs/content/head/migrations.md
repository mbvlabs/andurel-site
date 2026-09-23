# Migrations & Seeding

SQL migrations under root `migrations/` are the schema source of truth. Named seed compositions live in root `seeds/`. Generated projects embed migration SQL so CLI commands, tests, and production binaries share one history.

## Create and apply migrations

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel db migrate status
```

Generated file shape (`migrations/00015_create_products_table.sql`):

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id uuid NOT NULL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    name VARCHAR(255) NOT NULL,
    price_cents INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
```

| Command | Behavior |
| --- | --- |
| `andurel db migrate up` | Apply pending migrations |
| `andurel db migrate down` | Roll back the most recent migration |
| `andurel db migrate up-to [version]` | Apply through a version |
| `andurel db migrate down-to [version]` | Roll back down to a version |
| `andurel db migrate reset` | Roll back all, then re-apply |
| `andurel db migrate fix` | Re-number migrations to fix gaps |
| `andurel db migrate status` | Show applied / pending state |

`migrate up` and `status` print Goose text:

```text
OK   00015_create_products_table.sql (8.1ms)
```

```text
    Applied At                  Migration
    =======================================
    2026-09-23 14:02:11         00001_create_river_migration_table.sql
    2026-09-23 14:05:02         00015_create_products_table.sql
```

narsilc output does not replace migrations. Write schema changes as SQL first, apply them, then `andurel generate model` / `scaffold` from the migrated table. See [db](/docs/head/db).

## Seeds

```bash
andurel db seed
andurel db seed development
andurel db seed --list
andurel db rebuild --seed test
```

```bash
andurel db seed --list --json
```

```json
{
  "ok": true,
  "data": { "names": ["default", "development", "test"] },
  "summary": "Found 3 seed sets"
}
```

Seeds should call exported model and factory APIs. Keep development, test, and production-safe seed intent separate. `andurel db rebuild` drops and recreates the database, migrates, and seeds (pass `--skip-seed` to stop after migrate).

See [Factories](/docs/head/factories).

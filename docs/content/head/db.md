# db

Andurel wraps PostgreSQL lifecycle, embedded Goose migrations, and named seed operations under `andurel db`.

## Database lifecycle

```bash
andurel db create
andurel db drop
andurel db nuke
andurel db rebuild
andurel db rebuild --skip-seed
```

Human output for create:

```text
Database "orbit" created successfully.
```

`rebuild` drops and recreates the configured database, applies migrations, and runs a seed. Destructive commands prompt by default; `--force` permits protected system-database names and should be reserved for explicit automation.

```text
Database "orbit" rebuilt successfully.
```

## Migrations

Create a new SQL migration with the generator, then apply it:

```bash
andurel generate migration add_status_to_products
andurel db migrate up
andurel db migrate down
andurel db migrate status
andurel db migrate reset
andurel db migrate up-to 12
andurel db migrate down-to 8
andurel db migrate fix
```

Creating the SQL file is [generate migration](/docs/head/generate-migration); applying and rolling back is `andurel db migrate`.

`migrate` shells out to the pinned `bin/goose` against root `migrations/`. Goose status looks like:

```text
    Applied At                  Migration
    =======================================
    2026-09-23 14:02:11         00001_create_river_migration_table.sql
    2026-09-23 14:02:11         00002_create_river_job_and_leader_tables.sql
    Pending                     00015_add_status_to_products.sql
```

After `migrate up`, Goose prints each applied version:

```text
OK   00015_add_status_to_products.sql (12.4ms)
```

Migrations are SQL files in root `migrations/` and are embedded by the `migrations` package. See [Migrations & Seeding](/docs/head/migrations).

## Named seeds

```bash
andurel db seed
andurel db seed development
andurel db seed --list
andurel db rebuild --seed test
```

`--list` and named runs emit structured summaries when you pass `--json`:

```bash
andurel db seed --list --json
```

```json
{
  "ok": true,
  "data": {
    "names": ["default", "development", "test"]
  },
  "summary": "Found 3 seed sets"
}
```

```bash
andurel db seed development --json
```

```json
{
  "ok": true,
  "data": {
    "name": "development",
    "output": ["seeded development users"]
  },
  "summary": "Ran \"development\" seed"
}
```

Root `seeds/` owns the registry and composes model factories. Keep test, development, and production-safe seed intent separate.

## Inspect the database

```bash
andurel db console
andurel tool dblab
```

Both use the current project connection settings. `db console` opens usql; `dblab` opens the terminal database UI.

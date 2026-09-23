# Database Commands

Andurel wraps PostgreSQL lifecycle, embedded Goose migrations, and named seed operations under `andurel db`.

## Database lifecycle

```bash
andurel db create
andurel db drop
andurel db nuke
andurel db rebuild
andurel db rebuild --skip-seed
```

`rebuild` drops and recreates the configured database, applies migrations, and runs a seed. Destructive commands prompt by default; `--force` permits protected system-database names and should be reserved for explicit automation.

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

Migrations are SQL files in root `migrations/` and are embedded by the `migrations` package.

## Named seeds

```bash
andurel db seed
andurel db seed development
andurel db seed --list
andurel db rebuild --seed test
```

Root `seeds/` owns the registry and composes model factories. Keep test, development, and production-safe seed intent separate.

## Inspect the database

```bash
andurel db console
andurel tool dblab
```

Both use the current project connection settings. `db console` opens usql; `dblab` opens the terminal database UI.

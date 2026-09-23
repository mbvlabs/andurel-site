# Database Commands

Andurel wraps PostgreSQL lifecycle, embedded Goose migrations, and named seed operations.

## Database lifecycle

```bash
andurel database create
andurel database drop
andurel database nuke
andurel database rebuild
andurel database rebuild --skip-seed
```

`rebuild` drops and recreates the configured database, applies migrations, and runs a seed. Destructive commands prompt by default; `--force` permits protected system-database names and should be reserved for explicit automation.

## Migrations

```bash
andurel database migrate new add_status_to_products
andurel database migrate up
andurel database migrate down
andurel database migrate status
andurel database migrate reset
andurel database migrate up-to 12
andurel database migrate down-to 8
andurel database migrate fix
```

Migrations are SQL files in root `migrations/` and are embedded by the `migrations` package.

## Named seeds

```bash
andurel database seed
andurel database seed development
andurel database seed --list
andurel database rebuild --seed test
```

Root `seeds/` owns the registry and composes model factories. Keep test, development, and production-safe seed intent separate.

## Inspect the database

```bash
andurel console
andurel tool dblab
```

Both use the current project connection settings. `console` opens usql; `dblab` opens the terminal database UI.

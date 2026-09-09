# Database Commands

Andurel wraps common PostgreSQL lifecycle and migration operations.

## Database lifecycle

```bash
andurel database create
andurel database drop
andurel database nuke
andurel database rebuild
andurel database seed
```

`rebuild` drops and recreates the database, applies migrations, and runs seeds. Pass `--skip-seed` when seed data is not needed.

## Migrations

```bash
andurel database migrate new add_status_to_products
andurel database migrate up
andurel database migrate down
andurel database migrate status
andurel database migrate reset
```

Use `up-to` and `down-to` to target a version. Use `fix` to renumber migration gaps.

## Database console

Open a console using the current `.env` connection settings:

```bash
andurel console
```

The project toolchain provides `usql` and `dblab` for command-line and browser-based inspection.

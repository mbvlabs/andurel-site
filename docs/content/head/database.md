# Database

Andurel v2 uses PostgreSQL through the standalone storage package. One shared `database/sql` pool backed by pgx serves Bun, sqlc, and River.

## Connect through storage

The application constructs `*storage.Postgres` from a validated `storage.Config`, then exposes it as the small `storage.Connection` boundary. Models receive that connection through constructors.

```go
type Products struct {
    db storage.Connection
}

func NewProducts(db storage.Connection) Products {
    return Products{db: db}
}
```

Use `db.Executor()` for Bun operations and `db.DB()` for standard SQL consumers. Keep entities, validation, relationships, and persistence methods together in their owning model file.

## Schema and migrations

SQL migrations under root `migrations/` are the source of truth. The package embeds `*.sql` files for CLI, tests, and deployment tooling. Bun tags and sqlc output do not replace migrations.

```bash
andurel database migrate new create_products_table
andurel database migrate up
andurel generate model Product
```

## Transactions

Transactions bridge Bun and standard SQL. Use `storage.RunInTransaction` for multi-step workflows. Within the transaction, `Executor()` serves Bun model work and `SQL()` can be passed to `queries.WithTx` for sqlc. This keeps all operations on one transaction and pool.

## Factories and seeds

Factories live beside models under `models/factories`. Named seed compositions live in root `seeds/` and are run through `cmd/seeds`. Check generated factories for drift with `andurel generate factories --check --json` and sync them explicitly.

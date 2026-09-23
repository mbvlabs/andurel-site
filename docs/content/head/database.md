# Database

Andurel v2 uses PostgreSQL through the standalone storage package. One shared **pgx** pool serves narsilc-backed models and River.

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

Pass `db` to narsilc with `queries.New(db)`. Keep entities, validation, relationships, and persistence methods together in their owning model file. See [SQL Queries](/docs/head/sql-queries).

## Schema and migrations

SQL migrations under root `migrations/` are the source of truth. The package embeds `*.sql` files for CLI, tests, and deployment tooling. narsilc output does not replace migrations.

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel generate model Product
```

## Transactions

Use `storage.RunInTransaction` for multi-step workflows. Within the transaction, pass `tx` to `queries.New(tx)` so every statement shares one PostgreSQL transaction and the same pool.

## Factories and seeds

Factories live beside models under `models/factories`. Named seed compositions live in root `seeds/` and are run through `cmd/seeds` / `andurel db seed`. Check generated factories for drift with `andurel sync factories --check --json` and refresh them with `--sync`.

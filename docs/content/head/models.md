# Models

Models are application-owned Go types. They receive `storage.Connection` through constructors and are registered with Fx via `models.Module`.

## Construct and inject

```go
type Products struct {
    db storage.Connection
}

func NewProducts(db storage.Connection) Products {
    return Products{db: db}
}
```

Keep entities, validation, relationships, and persistence methods together in the owning model file. Pass `queries.New(db)` or `queries.New(tx)` for narsilc access — see [Queries](/docs/head/queries).

## Generate from migrations

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel generate model Product
```

Model generation reads the migrated schema and writes (or updates) the model API you can edit afterward.

## Transactions

For multi-step work, use `storage.RunInTransaction` and pass the transaction into model methods or query clients so every statement shares one PostgreSQL transaction. See [Getting Started](/docs/head/database).

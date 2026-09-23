# Queries

Andurel v2 persists through **narsilc**: you write annotated SQL under `models/queries/`, and the tool generates Go that fills application-owned result structs under `models/internal/queries`.

## Generate a query file

```bash
andurel generate query UserReport
andurel generate query UserReport --table users
andurel generate query UserReport --dry-run --json
```

The command creates `models/queries/user_report.sql`. Without `--table`, the file contains commented annotation examples. With `--table`, it includes a starter annotated query for the existing table.

After editing the SQL, compile Go code:

```bash
andurel sync queries
```

Output is written under `models/internal/queries`. The command is a no-op when `models/queries` contains no annotated SQL files. `andurel run`, `andurel build`, scaffolding, and related flows also regenerate narsilc output when query files exist. The narsilc binary version comes from `andurel.toml` `[tools]`; digests live in `andurel.lock`.

## Annotate SQL

narsilc queries use name annotations and result cardinality:

```sql
-- name: GetProduct :one
SELECT *
FROM products
WHERE id = $1
LIMIT 1;

-- name: ListProducts :many
-- @order id
SELECT *
FROM products;

-- name: CountProducts :one
SELECT count(*)
FROM products;

-- name: CreateProduct :one
INSERT INTO products (id, sku, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
```

| Annotation | Meaning |
| --- | --- |
| `:one` | Exactly one row (or error) |
| `:many` | Zero or more rows |
| `:exec` | Statement with no result set |
| `-- @order` | Stable ordering hint for list queries |

Model generators emit a matching `*.sql` file beside CRUD methods. Hand-written reports follow the same annotation style.

## Keep generated queries behind models

Only model packages should import the internal generated query package. Controllers and services depend on application model APIs and projection types, not database-shaped rows or parameters.

```go
func (p Products) Find(ctx context.Context, id uuid.UUID) (Product, error) {
    entity, err := p.queries.GetProduct[Product](ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return Product{}, ErrNotFound
        }
        return Product{}, err
    }
    return entity, nil
}
```

Map generated rows to an application-owned projection inside the model method when the public API should differ from the SQL shape.

## Share transactions

narsilc clients accept `storage.Connection` and `storage.Transaction` directly because both implement pgx `DBTX`:

```go
err := storage.RunInTransaction(ctx, connection,
    func(ctx context.Context, tx storage.Transaction) error {
        return queries.New(tx).RecordProductCreated(ctx, productID)
    },
)
```

See [Getting Started](/docs/head/database) and [Models](/docs/head/models).

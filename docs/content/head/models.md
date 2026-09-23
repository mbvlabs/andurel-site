# Models

Models are application-owned Go types. They receive `storage.Connection` through constructors, wrap a narsilc client, and register with Fx via `models.Module`. Controllers and services call model APIs. They never import `models/internal/queries` directly.

## Entity versus query client

Two layers sit in the same package:

| Layer | What it is | Who uses it |
| --- | --- | --- |
| Entity | Domain row type (`Product`, `User`) with `andurel` tags | Controllers map to payloads; factories sync fields |
| Plural API | `Products`, `Users` holding `*queries.Queries` | Controllers, services, seeds |
| narsilc client | Generated under `models/internal/queries` | Only the owning model package |

Hand-written SQL lives in `models/queries/*.sql`. `andurel sync queries` compiles it into the internal client. The plural API constructs that client with `queries.New(db)` or `queries.New(tx)`.

Mark the table above the imports so generators and factory sync can find the Entity:

```go
// andurel:table products

type Product struct {
    ID         uuid.UUID          `andurel:"id"`
    Name       string             `andurel:"name"`
    PriceCents int32              `andurel:"price_cents"`
    CreatedAt  pgtype.Timestamptz `andurel:"created_at"`
    UpdatedAt  pgtype.Timestamptz `andurel:"updated_at"`
}
```

Nullable columns use `pgtype.*` when `database.nullType` is `pgtype.Null` (the default), or pointers when set to `pointer` in `andurel.toml`.

## Construct and inject

Generated plural APIs look like this:

```go
type Products struct {
    queries *queries.Queries
}

func NewProducts(db storage.Connection) Products {
    return Products{queries: queries.New(db)}
}

func (p Products) WithTx(tx storage.Transaction) Products {
    return Products{queries: queries.New(tx)}
}
```

`NewProducts` binds the pool. `WithTx` returns a value copy whose narsilc client uses the same PostgreSQL transaction. Persistence methods on `Products` call generated query methods and map rows into `Product` entities.

Register constructors in `models.Module` so Fx can supply them to controllers and services:

```go
var Module = fx.Module(
    "models",
    fx.Provide(
        NewUsers,
        NewProducts,
    ),
)
```

Keep entities, validation, relationships, and persistence methods together in the owning model file.

## Generate from migrations

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel generate model Product
andurel generate model Product --mode read-only
andurel generate model Product --update --yes
```

| Flag | Meaning |
| --- | --- |
| `--mode` | `crud` (default), `read-only`, or `create-only` |
| `--table-name` | Override the default table name |
| `--primary-key` | Skip interactive primary-key detection |
| `--skip-factory` | Do not create or update `models/factories` |
| `--update` | Refresh an existing model from migration changes |
| `--dry-run` / `--diff` / `--json` | Preview structured mutations |

Model generation reads migration history, writes (or updates) the model API and matching `models/queries/*.sql`, and updates `models.Module`. You own the resulting Go afterward. After editing SQL by hand, run `andurel sync queries --json` so the internal client matches.

## Transactions

For multi-step work, use `storage.RunInTransaction` and `WithTx` so every statement shares one PostgreSQL transaction:

```go
err := storage.RunInTransaction(ctx, db,
    func(ctx context.Context, tx storage.Transaction) error {
        products := models.NewProducts(db).WithTx(tx)
        _, err := products.Create(ctx, data)
        return err
    },
)
```

Prefer injecting `models.Products` and calling `products.WithTx(tx)` over constructing a new API mid-handler when Fx already supplied one.

See [Getting Started](/docs/head/database) and [Queries](/docs/head/queries).

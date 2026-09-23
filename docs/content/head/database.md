# Getting Started

Andurel v2 uses PostgreSQL through `github.com/mbvlabs/andurel/pkg/storage`. One shared **pgx/v5** pool serves narsilc-backed models and River via `storage.Connection`.

## Connection lifecycle

The web and queue composition roots construct `*storage.Postgres` and publish it as `storage.Connection`. Construction validates configuration, builds a pgx pool, applies pool limits, and pings PostgreSQL. A failed ping prevents startup; an Fx lifecycle hook closes the pool during shutdown.

```go
db, err := storage.NewPostgres(ctx, cfg)
if err != nil {
    return nil, err
}
lifecycle.Append(fx.Hook{
    OnStop: func(context.Context) error { return db.Close() },
})
```

`Connection` implements pgx `DBTX` (`Exec`, `Query`, `QueryRow`) plus `Health(ctx)` and `BeginTransaction(ctx)`. Pass it directly to narsilc: `queries.New(db)`. Do not create separate pools for models, migrations tooling, and River.

## PostgreSQL configuration

Start from `storage.DefaultConfig()` and replace application identity and credentials. Application `config.NewDatabase` maps environment variables into this type. `Config.Validate` requires a PostgreSQL kind, host, port, name, user, and supported SSL mode.

```go
cfg := storage.DefaultConfig()
cfg.Name = "orbit_production"
cfg.User = "orbit"
cfg.Password = secret
cfg.SSLMode = "verify-full"
cfg.ApplicationName = "orbit-web"
db, err := storage.NewPostgres(ctx, cfg)
```

Programmatic options override config fields. `WithMaxOpenConnections` maps to pgxpool `MaxConns`. Prefer max open connections and `ConnectionMaxIdleTime` over idle-connection knobs that are not applied to the pgx pool. See [Configuration](/docs/head/configuration) for the full `DB_*` table.

## Models and queries

Application-owned models depend on `storage.Connection`, wrap a narsilc client, and expose `WithTx` for transactional work:

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

See [Models](/docs/head/models) and [Queries](/docs/head/queries).

## Transactions

`RunInTransaction` commits only when its callback succeeds and otherwise rolls back. The `Transaction` value also implements pgx `DBTX`, so narsilc takes it directly:

```go
err := storage.RunInTransaction(ctx, connection,
    func(ctx context.Context, tx storage.Transaction) error {
        products := models.NewProducts(connection).WithTx(tx)
        product, err := products.Create(ctx, data)
        if err != nil {
            return err
        }
        _, err = queue.InsertTx(ctx, tx, jobs.ProductCreatedArgs{ID: product.ID}, nil)
        return err
    },
)
```

Calling a model method that uses the original connection from inside the callback escapes the transaction. Pass `tx` (or a model constructed with `WithTx`) into participating operations instead.

## Migrations and test databases

`RunMigrations` applies Goose migrations from an `fs.FS`. Generated projects embed SQL under root `migrations/`. See [Migrations & Seeding](/docs/head/migrations).

`NewTestCluster` starts a PostgreSQL 17 Alpine container for isolated test databases. Prefer one cluster per test package. See [Testing](/docs/head/testing).

## River insertion and processing

`storage.NewQueueInsert` builds an insert-only River client on the shared pgx pool (`riverpgxv5`). Publish it as `storage.InsertQueue` in the web process. `storage.NewQueueProcessor` owns processing and belongs in `cmd/queue`.

Use `InsertTx` / `InsertManyTx` when the job must share a PostgreSQL transaction with domain writes. See [Queues](/docs/head/queues).

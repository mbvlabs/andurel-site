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

Start from `storage.DefaultConfig()` and replace application identity and credentials. `Config.Validate` requires a PostgreSQL kind, host, port, name, user, and supported SSL mode.

```go
cfg := storage.DefaultConfig()
cfg.Name = "orbit_production"
cfg.User = "orbit"
cfg.Password = secret
cfg.SSLMode = "verify-full"
cfg.ApplicationName = "orbit-web"
db, err := storage.NewPostgres(ctx, cfg)
```

Programmatic options override config fields. `WithMaxOpenConnections` maps to pgxpool `MaxConns`. Prefer max open connections and `ConnectionMaxIdleTime` over idle-connection knobs that are not applied to the pgx pool.

## Models and queries

Application-owned models depend on `storage.Connection` and call narsilc-generated clients. See [Models](/docs/head/models) and [Queries](/docs/head/queries).

```go
type Products struct {
    db storage.Connection
}

func NewProducts(db storage.Connection) Products {
    return Products{db: db}
}
```

## Transactions

`RunInTransaction` commits only when its callback succeeds and otherwise rolls back. The `Transaction` value also implements pgx `DBTX`, so narsilc takes it directly:

```go
err := storage.RunInTransaction(ctx, connection,
    func(ctx context.Context, tx storage.Transaction) error {
        q := queries.New(tx)
        if _, err := q.InsertProduct(ctx, params); err != nil {
            return err
        }
        return q.RecordProductCreated(ctx, productID)
    },
)
```

Calling a model method that uses the original connection from inside the callback escapes the transaction. Pass `tx` (or a model constructed with it) into participating operations instead.

## Migrations and test databases

`RunMigrations` applies Goose migrations from an `fs.FS`. Generated projects embed SQL under root `migrations/`. See [Migrations & Seeding](/docs/head/migrations).

`NewTestCluster` starts a PostgreSQL 17 Alpine container for isolated test databases. Prefer one cluster per test package. See [Testing](/docs/head/testing).

## River insertion and processing

`QueueInsert` belongs in the web process. `QueueProcessor` owns processing and belongs in `cmd/queue`. River uses `riverpgxv5` on the same pgx pool. See [Queues](/docs/head/queues).

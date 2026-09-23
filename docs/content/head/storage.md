# Storage

`github.com/mbvlabs/andurel/pkg/storage` is the shared data-infrastructure boundary for PostgreSQL, narsilc-friendly connections, migrations, tests, and River. It wraps a **pgx/v5** pool and exposes `storage.Connection` so narsilc clients and River share one pool.

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

`Connection` implements pgx `DBTX` (`Exec`, `Query`, `QueryRow`) plus `Health(ctx)` and `BeginTransaction(ctx)`. Pass it directly to narsilc: `queries.New(db)`. Do not create separate pools for models, migrations tooling, and River — their limits, telemetry, shutdown, and transactions would diverge.

## PostgreSQL configuration and options

Start from `storage.DefaultConfig()` and replace application identity and credentials. Defaults include a five-second connect timeout, statement and description caches of 512, 25 max open connections, a one-hour maximum lifetime, and a 30-minute maximum idle time.

`Config.Validate` requires a PostgreSQL kind, host, port, name, user, and supported SSL mode. Counts and durations cannot be negative. `DatabaseURL()` safely constructs the URL.

```go
cfg := storage.DefaultConfig()
cfg.Name = "orbit_production"
cfg.User = "orbit"
cfg.Password = secret
cfg.SSLMode = "verify-full"
cfg.ApplicationName = "orbit-web"
db, err := storage.NewPostgres(ctx, cfg)
```

Programmatic options override config fields. They include `WithDatabaseURL`, `WithTLSConfig`, `WithRuntimeParameters`, connect and pool limits, statement-cache sizes, application name, and `WithOpenTelemetry`. A custom URL retains pgx URL parameters; explicit options still win.

`WithMaxOpenConnections` maps to pgxpool `MaxConns`. `WithMaxIdleConnections` / `DB_MAX_IDLE_CONNECTIONS` remain for compatibility and are **not** applied to the pgx pool — prefer max open connections and `ConnectionMaxIdleTime`.

## narsilc and transactions

Application-owned models normally depend on `storage.Connection` and call narsilc-generated clients:

```go
func (m Products) Find(ctx context.Context, id uuid.UUID) (Product, error) {
    row, err := queries.New(m.db).FindProduct(ctx, id)
    if err != nil {
        return Product{}, err
    }
    return mapProduct(row), nil
}
```

For multi-step work, `RunInTransaction` commits only when its callback succeeds and otherwise rolls back. The `Transaction` value also implements pgx `DBTX`, so narsilc takes it directly:

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

`RunMigrations` applies Goose migrations from an `fs.FS`. Generated projects embed SQL migrations so tests, CLI commands, and production binaries use the same schema history. narsilc output does not replace migrations. Goose opens a short-lived `database/sql` adapter inside `RunMigrations` only.

`NewTestCluster` starts a PostgreSQL 17 Alpine container. A cluster can create isolated databases and produce configs for them. Prefer one cluster per test package and isolated databases per test or suite over a new container for every case.

## River insertion and processing

`QueueInsert` satisfies the narrow `InsertQueue` interface and belongs in the web process. `QueueProcessor` owns processing and belongs in `cmd/queue` with workers and periodic jobs. River uses `riverpgxv5` on the same pgx pool.

```go
inserter, err := storage.NewQueueInsert(connection, insertCfg)
processor, err := storage.NewQueueProcessor(
    connection,
    workerCfg,
    storage.WithRiverWorkers(workers),
    storage.WithRiverLogger(logger),
    storage.WithRiverPeriodicJobs(periodicJobs...),
)
```

`DefaultQueueConfig()` carries operational defaults but intentionally has no processor queues. Generated worker configuration adds the default queue; insertion configuration clears queues. `QueueConfig.Clone()` copies maps and slices before mutation. Validation covers concurrency, queue names, schemas, retention, polling, timeouts, and inconsistent rescue thresholds.

All River integration points remain available through `WithRiver*` options, including hooks, middleware, retry policy, error and stuck-job handlers, queues, schema, periodic jobs, test config, and polling policy. Transactional inserts that must share a boundary use `InsertTx` with the same `storage.Transaction`.

## Telemetry and failure behavior

`Config.OpenTelemetry` enables default pgx instrumentation. `WithOpenTelemetry` accepts explicit tracer and meter providers, attributes, SQL span naming, and query-parameter capture. Query parameters may contain secrets or personal data, so enable capture deliberately.

Constructor failures are wrapped with `storage:` context. Health failures close the partially created pool. Transaction helpers propagate callback and commit errors; a deferred rollback protects all non-committed paths. Keep environment parsing in application `config/` so package construction always receives validated typed policy.

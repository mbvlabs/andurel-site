# Storage

`github.com/mbvlabs/andurel/pkg/storage` is the shared data-infrastructure boundary for PostgreSQL, Bun, sqlc, migrations, tests, and River. It creates one underlying `database/sql` pool and exposes the view needed by each consumer.

## Connection lifecycle

The web and queue composition roots construct `*storage.Postgres` and publish it as `storage.Connection`. Construction validates configuration, builds a pgx-backed pool, wraps it with Bun, applies pool limits, and pings PostgreSQL. A failed ping prevents startup; an Fx lifecycle hook closes the pool during shutdown.

```go
db, err := storage.NewPostgres(ctx, cfg)
if err != nil {
    return nil, err
}
lifecycle.Append(fx.Hook{
    OnStop: func(context.Context) error { return db.Close() },
})
```

`Connection` exposes `Executor()` for Bun, `DB()` for sqlc, River, and other `database/sql` consumers, `Health(ctx)` for an explicit runtime check, and `BeginTransaction()` for a shared transaction. Do not create separate pools for each library: their limits, telemetry, shutdown, and transactions would diverge.

## PostgreSQL configuration and options

Start from `storage.DefaultConfig()` and replace application identity and credentials. Defaults include a five-second connect timeout, statement and description caches of 512, 25 open and idle connections, a one-hour maximum lifetime, and a 30-minute maximum idle time.

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

## Bun, sqlc, and transactions

Application-owned models normally depend on `storage.Connection` and call `Executor()`:

```go
func (m Products) Find(ctx context.Context, id uuid.UUID) (Product, error) {
    product := Product{ID: id}
    err := m.db.Executor().NewSelect().Model(&product).WherePK().Scan(ctx)
    return product, err
}
```

Construct sqlc queries with `dbqueries.New(connection.DB())`. For multi-step work, `RunInTransaction` commits only when its callback succeeds and otherwise rolls back. Its `Transaction` exposes Bun and standard SQL over the same PostgreSQL transaction:

```go
err := storage.RunInTransaction(ctx, connection, nil,
    func(ctx context.Context, tx storage.Transaction) error {
        if _, err := tx.Executor().NewInsert().Model(&product).Exec(ctx); err != nil {
            return err
        }
        return queries.WithTx(tx.SQL()).RecordProductCreated(ctx, product.ID)
    },
)
```

Calling a model method that uses the original connection from inside the callback escapes the transaction. Pass `tx.Executor()` into participating model operations instead.

## Migrations and test databases

`RunMigrations` applies Goose migrations from an `fs.FS`. Generated projects embed SQL migrations so tests, CLI commands, and production binaries use the same schema history. Bun tags and sqlc output do not replace migrations.

`NewTestCluster` starts a PostgreSQL 17 Alpine container. A cluster can create isolated databases and produce configs for them. Prefer one cluster per test package and isolated databases per test or suite over a new container for every case.

## River insertion and processing

`QueueInsert` satisfies the narrow `InsertQueue` interface and belongs in the web process. `QueueProcessor` owns processing and belongs in `cmd/queue` with workers and periodic jobs.

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

All River integration points remain available through `WithRiver*` options, including hooks, middleware, retry policy, error and stuck-job handlers, queues, schema, periodic jobs, test config, and polling policy.

## Telemetry and failure behavior

`Config.OpenTelemetry` enables default pgx instrumentation. `WithOpenTelemetry` accepts explicit tracer and meter providers, attributes, SQL span naming, and query-parameter capture. Query parameters may contain secrets or personal data, so enable capture deliberately.

Constructor failures are wrapped with `storage:` context. Health failures close the partially created pool. Transaction helpers propagate callback and commit errors; a deferred rollback protects all non-committed paths. Keep environment parsing in application `config/` so package construction always receives validated typed policy.

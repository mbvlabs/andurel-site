# Dependency Injection

Andurel applications use [Fx](https://uber-go.github.io/fx/) for dependency injection and process lifecycle. Framework packages expose ordinary constructors; Fx belongs to the generated application in `cmd/app`, `cmd/queue`, and `cmd/ssr`.

## Composition roots

Each process builds its own graph. Load validated environment before Fx, then provide only the collaborators that process should own:

```go
if err := config.LoadEnvironment(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}

app := fx.New(
    fx.Provide(func() context.Context { return ctx }),
    config.Module,
    databaseModule,
    queueInsertModule,
    models.Module,
    controllers.Module,
    cookies.Module,
    router.Module,
    fx.Invoke(startServer),
)
```

Constructors return concrete types or small interfaces (`storage.Connection`, `storage.InsertQueue`, `*inertia.Renderer`, `*kiks.Jar`). Do not hide dependencies in Echo request context.

Publish the pool as both the interface and the concrete type when lifecycle needs `Close`:

```go
var databaseModule = fx.Module(
    "database",
    fx.Provide(fx.Annotate(
        newDatabase,
        fx.As(new(storage.Connection)),
        fx.As(fx.Self()),
    )),
)
```

Application packages register focused modules. Generated models look like:

```go
var Module = fx.Module(
    "models",
    fx.Provide(
        NewUsers,
        NewTokens,
        NewProducts,
    ),
)
```

## Lifecycle

Register start/stop hooks for resources that must clean up: database pools, queue processors, SSR runtimes, and HTTP servers. A failed ping or constructor error should prevent process start.

```go
lc.Append(fx.Hook{
    OnStart: func(ctx context.Context) error {
        return db.Health(ctx)
    },
    OnStop: func(ctx context.Context) error {
        return db.Close()
    },
})
```

Web needs HTTP and queue insertion. Queue needs a River processor and workers, but not the Echo server. Keep provider sets distinct so a process cannot accidentally acquire a capability it should not own.

## Replacing implementations

Prefer the smallest public interface at the composition root. Tests can provide fake email senders or in-memory collaborators without rewriting controllers. Package updates should not silently introduce new environment contracts; environment mapping stays in application `config/`.

See [Configuration](/docs/head/configuration), [Request Lifecycle](/docs/head/request-lifecycle), and [Framework Packages](/docs/head/framework-packages).

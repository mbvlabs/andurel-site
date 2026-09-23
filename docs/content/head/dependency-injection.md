# Dependency Injection

Andurel applications use [Fx](https://uber-go.github.io/fx/) for dependency injection and process lifecycle. Framework packages expose ordinary constructors; Fx belongs to the generated application.

## Composition roots

Each process (`cmd/app`, `cmd/queue`, `cmd/ssr` as applicable) has its own Fx graph. Constructors return concrete types or small interfaces (`storage.Connection`, email senders, `*inertia.Renderer`). Do not hide dependencies in Echo request context.

```go
fx.Provide(
    config.NewDatabase,
    storage.NewPostgres,
    models.Module,
    controllers.Module,
)
```

## Lifecycle

Register start/stop hooks for resources that must clean up — database pools, queue processors, SSR runtimes, and HTTP servers. A failed ping or constructor error should prevent process start.

## Replacing implementations

Prefer the smallest public interface. Tests can provide fake senders or in-memory collaborators at the composition root. Package updates should not silently introduce new environment contracts — environment mapping stays in application `config/`. See [Configuration](/docs/head/configuration) and [Framework Packages](/docs/head/framework-packages).

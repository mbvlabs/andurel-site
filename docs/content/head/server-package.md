# Server

`github.com/mbvlabs/andurel/pkg/server` is a small wrapper around `net/http.Server`. It supplies bounded defaults, request base contexts, coordinated shutdown components, and consistent start errors. Routing, middleware, TLS termination, health policy, and process supervision remain application concerns.

## Construction

`server.New` receives the process context, listen host and port, environment, an `http.Handler`, optional shutdown components, and functional options.

```go
srv := server.New(
    appCtx, httpCfg.Host, httpCfg.Port, appCfg.Environment,
    router.Handler, nil,
    server.WithTimeouts(
        httpCfg.IdleTimeout, httpCfg.ReadTimeout, httpCfg.WriteTimeout,
    ),
)
```

The context becomes the base context for accepted connections. `Start` blocks while serving, so the generated application runs it in a supervised background goroutine.

## Timeout policy

`DefaultConfig()` sets a 120-second idle timeout, 10-second read timeout, and 30-second write timeout. `WithTimeouts` replaces all three. Validation rejects negative values; zero disables a timeout.

Read timeout bounds request ingestion, write timeout bounds response writing, and idle timeout bounds keep-alive connections. Long-lived SSE may require a separate write-timeout strategy. Do not disable global bounds casually to accommodate one route.

## Start and shutdown

`Server.Start` calls `ListenAndServe`, treats `http.ErrServerClosed` as normal, and wraps other failures as `server: listen`. `Server.Shutdowners` contains the HTTP server plus caller-supplied values implementing:

```go
type Shutdowner interface {
    Shutdown(context.Context) error
}
```

The generated Fx stop hook invokes them with a bounded context, joins failures, and waits for the serve loop. This makes graceful shutdown explicit: stop accepting connections, drain active requests, stop owned components, and return control to process supervision.

## Fx lifecycle integration

The package does not import Fx. The composition root owns hooks and total deadlines:

```go
lifecycle.Append(fx.Hook{
    OnStart: func(context.Context) error {
        done = startInBackground(appCtx, "server", func(ctx context.Context) error {
            return srv.Start(ctx, appCfg.Environment)
        })
        return nil
    },
    OnStop: func(ctx context.Context) error {
        return stopAndWait(ctx, shutdownAll, done)
    },
})
```

Asynchronous listen failures must reach supervision and readiness state. The generated background helper owns that policy while the package remains usable without Fx.

## Deployment boundary

Host and port describe the local listener, not the public URL. Configure protocol and domain separately. TLS commonly terminates at a proxy or platform; a Go-owned TLS listener belongs in application composition rather than environment-dependent package behavior.

The package does not create readiness endpoints. Build them from meaningful dependencies such as `storage.Connection.Health` with short timeouts. Liveness should normally report process health without failing for every transient dependency outage.

## Testing and extension

Test handlers directly with `httptest`; bind a real port only for connection, timeout, or shutdown behavior. Extra long-lived resources can be supplied as shutdowners, but preserve ownership: the component constructor should define how it stops, while the composition root defines ordering and the overall deadline.

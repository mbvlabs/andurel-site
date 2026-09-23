# Request Lifecycle

An Andurel v2 application is one codebase with multiple process graphs. Understanding which process owns a request (or a job) keeps composition roots honest.

## Web process (`cmd/app`)

Typical path for an HTTP request:

```text
.env / process env
        |
        v
config.LoadEnvironment()
        |
        v
Fx build: config, storage.Connection, models,
          services, controllers, cookies.Jar, router
        |
        v
pkg/server starts Echo (config.HTTP timeouts)
        |
        v
middleware: kiks bag -> CSRF -> auth -> Inertia (when set)
        |
        v
controller (injected collaborators)
        |
        +---> inertia.Renderer Page / Redirect / Location
        |
        +---> hypermedia Templ / Datastar
        |
        v
kiks persist dirty cookies + flashes
        |
        v
response; on signal: graceful shutdown, close pool
```

1. `config.LoadEnvironment()` loads `.env` and process environment before Fx.
2. Fx constructs validated config providers, `storage.Postgres` as `storage.Connection`, models, services, controllers, `cookies.NewJar`, and the router.
3. `pkg/server` starts Echo with HTTP timeouts from `config.HTTP`.
4. Middleware runs in order: kiks bag load (`jar.EchoMiddleware`), CSRF, auth gates, Inertia middleware when configured.
5. A controller handler receives injected collaborators (not ambient globals).
6. The handler validates input, calls model APIs, and either renders Inertia via `*inertia.Renderer` (`Page` / redirects) or returns Templ / Datastar through `pkg/hypermedia`.
7. kiks persists dirty bagged cookies and flashes before the response commits.
8. On signal, Fx stops the server with a graceful shutdown timeout and closes the pool.

Queue insertion belongs here: inject `storage.InsertQueue` (or `*storage.QueueInsert`) and call `Insert` / `InsertTx` after domain writes. Do not start River workers in the web process.

## Queue process (`cmd/queue`)

`cmd/queue` constructs the same `storage.Connection`, a River `QueueProcessor`, registered workers, telemetry, and email transports. It does not construct Echo or Inertia.

```text
Fx build: storage.Connection, QueueProcessor,
          workers, telemetry, email senders
        |
        v
processor starts on shared pgx pool (riverpgxv5)
        |
        v
job.Work(ctx, ...) with injected models / email
        |
        v
soft-stop from config.NewQueueWorker
        |
        v
shutdown: stop processor, close pool
```

1. Fx starts the processor against the shared pgx pool (`riverpgxv5`).
2. Workers run job `Work` methods with injected model and email collaborators.
3. Soft-stop and lifecycle config come from `config.NewQueueWorker`.
4. Shutdown stops the processor, then closes the pool.

Commit domain state and a River job together with `storage.RunInTransaction` and `InsertTx` in the web process so workers never see half-written rows. See [Queues](/docs/head/queues).

## SSR process (`cmd/ssr`)

Inertia SSR is optional and per-page.

```text
cmd/app Page(...).SSR()
        |
        +-- development --> Vite /__inertia_ssr
        |
        +-- otherwise ----> INERTIA_SSR_URL
                                    |
                                    v
                            cmd/ssr (Node bind
                            on INERTIA_SSR_LISTEN)
                                    |
                                    v
                            HTML fragment back to
                            renderer; client fallback
                            unless fail-fast
```

`cmd/ssr` owns the Node runtime bind (`INERTIA_SSR_LISTEN`). The web renderer posts to `INERTIA_SSR_URL` (or Vite's `/__inertia_ssr` in development). Pages still opt into SSR individually. See [SSR](/docs/head/inertia-ssr).

## Where state lives

| Concern | Owner |
| --- | --- |
| PostgreSQL rows | `storage.Connection` / `storage.Transaction` in the owning process |
| Session and flash | kiks cookies on the HTTP path ([Cookies & Sessions](/docs/head/cookies-sessions)) |
| Background side effects | River jobs processed in `cmd/queue` |
| Browser document / visits | Inertia renderer in `cmd/app`, optional Node SSR in `cmd/ssr` |

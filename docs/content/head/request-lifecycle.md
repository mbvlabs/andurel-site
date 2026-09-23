# Request Lifecycle

An Andurel v2 application is one codebase with multiple process graphs. Understanding which process owns a request (or a job) keeps composition roots honest.

## Web process

Typical path for an HTTP request:

1. `cmd/app` starts Fx and constructs validated config
2. `storage.Postgres` opens the shared pgx pool
3. Models, services, and controllers receive injected dependencies
4. Echo routes the request through middleware (including kiks session/cookie bags)
5. A controller renders Templ/Datastar or Inertia via `*inertia.Renderer`
6. `pkg/server` serves until graceful shutdown

## Queue process

`cmd/queue` constructs the same storage connection, a River processor, registered workers, telemetry, and email transports — but no HTTP server. Jobs inserted by the web process are processed here.

## SSR process

Inertia SSR is optional and per-page. `cmd/ssr` owns the Node runtime bind (`INERTIA_SSR_LISTEN`). The web renderer posts to `INERTIA_SSR_URL` (or Vite’s `/__inertia_ssr` in development). See [SSR](/docs/head/inertia-ssr).

## Where state lives

- Database work uses `storage.Connection` / transactions in the owning process
- Session and flash state use kiks cookies on the HTTP path — see [Cookies & Sessions](/docs/head/cookies-sessions)
- Background side effects (email, slow work) should commit domain state and a River job together, then run in `cmd/queue`

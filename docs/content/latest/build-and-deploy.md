# Build and Deploy

Andurel v2 builds generated database code, email styles, Templ components, CSS, optional Inertia assets, and the Go web application.

## Build

```bash
andurel build --version v2.0.0-dev
```

The command runs sqlc when query files exist, compiles email templates, generates Templ code, minifies Tailwind CSS, installs frontend dependencies and builds Vite assets for Inertia projects, downloads Go dependencies, and compiles the application release.

The JavaScript package manager comes from `andurel.lock`. It is independent from the Node runtime used by managed Inertia SSR.

## Deploy both process types

`andurel build` produces the web application from `cmd/app`. Build the queue entry point separately when the application processes jobs:

```bash
CGO_ENABLED=0 GOOS=linux go build -o orbit-queue ./cmd/queue
```

Deploy the two binaries as separate process types. They may share a database and telemetry backend, but each owns its Fx lifecycle and shutdown.

Apply root `migrations/` as a deliberate release step before traffic reaches code that needs the new schema. Back up important data before destructive changes.

## Environment and health

Provide production database, HTTP, session, auth, email, queue, telemetry, and optional Inertia SSR settings through the deployment environment. Set `ENVIRONMENT=production` so secure cookie and server policies apply.

The storage health check verifies the database. Configure the platform's startup and liveness checks to cover required dependencies and allow enough termination time for HTTP requests and queue jobs to stop cleanly.

## Development releases

The current `master` branch is the v2 development line. Pin the CLI commit and the independent package versions in `go.mod` for reproducible pre-release deployments. Do not treat the stable v1 `@latest` tag as v2 until `v2.0.0` is published.

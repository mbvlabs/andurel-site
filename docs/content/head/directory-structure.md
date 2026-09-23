# Directory Structure

Andurel v2 separates application-owned code, reusable framework modules, and executable process boundaries.

## Application layout

```text
cmd/app/                  web process and Fx composition root
cmd/queue/                queue worker process and lifecycle
cmd/seeds/                seed command entry point
config/                   environment loading and typed providers
controllers/              HTTP handlers and route registration
router/routes/            typed route declarations
router/appctx/            typed request metadata helpers
views/                    Templ components and optional Inertia root
models/                   entities, model APIs, and Fx module
models/factories/         generated test and seed builders
models/queries/           application-written sqlc SQL
models/internal/queries/  generated sqlc implementation
services/                 application workflows
queue/jobs/               River argument types
queue/workers/            River workers and registration
email/                    Templ email components
migrations/               embedded Goose SQL migrations
seeds/                    named seed sets and registry
telemetry/                application observability wiring
assets/                   embedded compiled assets
```

## Reusable framework packages

Generated applications import `github.com/mbvlabs/andurel/pkg/*` modules for routing, server, storage, validation, hypermedia, Inertia, and email behavior. Their versions are pinned in the application's `go.mod` and do not have to match the CLI version.

Unlike v1, these implementations are not copied into an application-owned `internal/` tree. Configuration, middleware policy, controllers, model behavior, and presentation remain application-owned.

## Project metadata

`andurel.lock` records the framework version, tools, extensions, database conventions, frontend adapter, JavaScript package manager, and the separate Inertia SSR runtime. Commit it with the project. `go.mod` remains the source of truth for framework package versions.

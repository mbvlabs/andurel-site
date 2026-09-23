# Directory Structure

Andurel v2 separates application-owned code, reusable framework modules, and executable process boundaries.

## Application layout

```text
cmd/app/                  web process and Fx composition root
cmd/queue/                queue worker process and lifecycle
cmd/ssr/                  Node SSR process owner (Inertia)
cmd/seeds/                seed command entry point
config/                   environment loading and typed providers
controllers/              HTTP handlers and route registration
router/cookies/           kiks jar, session payload, cookie module
router/middleware/        auth and related Echo middleware
router/routes/            typed route declarations
views/                    Templ components and Inertia root document
resources/js/             Inertia pages, layouts, Vite entry (when enabled)
models/                   entities, model APIs, and Fx module
models/factories/         generated test and seed builders
models/queries/           application-written narsilc SQL
models/internal/queries/  generated narsilc implementation
services/                 application workflows
queue/jobs/               River argument types
queue/workers/            River workers and registration
email/                    Templ email components
migrations/               embedded Goose SQL migrations
seeds/                    named seed sets and registry
telemetry/                application observability wiring
assets/                   embedded compiled assets
andurel.toml              project manifest (UI, database, tool pins)
andurel.lock              tool download digests
```

## Reusable framework packages

Generated applications import `github.com/mbvlabs/andurel/pkg/*` modules for routing, server, storage, validation, hypermedia, Inertia, email, kiks, and telemetry behavior. Their versions are pinned in the application's `go.mod` and do not have to match the CLI version.

Unlike v1, these implementations are not copied into an application-owned `internal/` tree. Configuration, middleware policy, controllers, model behavior, and presentation remain application-owned.

## Project metadata

`andurel.toml` records scaffold choices (frontend adapter, package manager, SSR runtime, database null strategy) and pinned tool versions. `andurel.lock` stores per-platform SHA-256 digests for those tools. Commit both with the project. `go.mod` remains the source of truth for framework package versions; use `andurel packages update` to bump them deliberately. See [Project Lock](/docs/head/project-lock).

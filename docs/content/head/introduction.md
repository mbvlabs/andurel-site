# Introduction

Andurel is a Rails-like web framework for Go. These **head** docs document Andurel **v2.0.0-alpha** — the current pre-release line (same commit as `master` today). APIs may still change before a stable `v2.0.0`.

## What v2 emphasizes

Andurel v2 keeps generated application code explicit while moving reusable infrastructure into independently versioned `pkg/*` modules. A new application uses Fx for dependency injection and lifecycle management, PostgreSQL through pgx and `storage.Connection`, narsilc for typed SQL into application-owned model structs, River for jobs, and **Inertia v3** (React/pnpm by default) for its UI. Templ with Datastar remains available for hypermedia pages.

The framework still favors one-time generation: a scaffold creates models, factories, controllers, routes, and views that belong to your application and can be edited normally.

## Development status

The `head` documentation URL slug stays `head`, but the product line named here is **v2.0.0-alpha**. Use the [latest documentation](/docs/latest/introduction) for the last v1 release line.

There is no automated upgrade path from v1 to v2. Create new v2 projects with the v2 CLI and migrate existing applications deliberately. See [Upgrade Guide](/docs/head/upgrade).

## Platform and requirements

v2.0.0-alpha requires **Go 1.27.1** or newer and targets Linux and macOS on amd64 and arm64. PostgreSQL is the supported database.

Continue to [Installation](/docs/head/installation) to install the CLI and create a project.

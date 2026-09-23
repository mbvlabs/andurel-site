# Upgrade Guide

Andurel v2 does not provide an automated upgrade path from v1. The `upgrade` command is not a v1-to-v2 migration tool.

## Choose a migration strategy

Keep production v1 applications on the `1-5-stable` line until you can create a fresh v2 scaffold and move application behavior deliberately. Install the v2 CLI separately (`github.com/mbvlabs/andurel/v2@v2.0.0-alpha`), generate a temporary v2 project with the same frontend and extensions, and compare its composition root and generated conventions with your application.

## Important changes

- Go **1.27.1** is the minimum version.
- Reusable routing, server, storage, validation, hypermedia, Inertia, email, and telemetry code is imported from independently versioned `pkg/*` modules instead of copied into application `internal/` packages.
- Configuration is split into validated typed providers rather than a global aggregate.
- Persistence uses **pgx** through `storage.Connection` and **narsilc**-generated queries into application-owned model structs (Bun and sqlc are gone).
- Model APIs are constructed with `storage.Connection` and injected through `models.Module`.
- Migrations and seeds live in root `migrations/` and `seeds/` packages.
- Queue insertion remains available to the web process, while processing runs from `cmd/queue`.
- Inertia v3 is the default UI (`--ui react/pnpm`); Templ/Datastar is `--ui templ/datastar`. Generators follow `andurel.lock` — there is no `--inertia` flag. Gonertia integrations require manual replacement; applications own `views/root.templ`.
- The lock records the JavaScript package manager separately from the Inertia SSR runtime (`cmd/ssr`).
- `andurel.toml` / `andurel.lock` pin tools and scaffold choices; use `andurel packages list` / `andurel packages update` for `pkg/*` upgrades in `go.mod`.

## Migrate by behavior

Move SQL migrations first, then model entities and their tests, services, jobs, routes/controllers, and views. Adapt imports and constructor dependencies as each layer moves. Replace Bun/sqlc call sites with narsilc clients (`queries.New(db)` / `queries.New(tx)`). Use the generated v2 authentication and middleware flow as the reference instead of copying v1 session plumbing forward.

Run both applications against disposable databases and compare user-visible behavior. Treat generated v2 files as a new baseline; do not overwrite a v1 project with a v2 scaffold or run broad scripted replacements without review.

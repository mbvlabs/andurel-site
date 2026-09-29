# What's New in v2

These head docs describe **v2.0.0-alpha.2**, the current tagged v2 pre-release. Requires Go **1.27.1** and module path `github.com/mbvlabs/andurel/v2`. Pin `github.com/mbvlabs/andurel/v2@v2.0.0-alpha.2` until a later tag.

v1 remains on [latest](/docs/latest/whats-new) and [1.5.x](/docs/1.5.x/whats-new). Do not mix v1 commands (`andurel new`, `andurel make`) with a v2 project.

## v2.0.0-alpha.2

Released 2026-09-29. Full notes: [GitHub release](https://github.com/mbvlabs/andurel/releases/tag/v2.0.0-alpha.2). Diff from the first alpha: [v2.0.0-alpha...v2.0.0-alpha.2](https://github.com/mbvlabs/andurel/compare/v2.0.0-alpha...v2.0.0-alpha.2).

This tag includes five changes that affect generated apps and these docs:

1. Multiple hosts ([#801](https://github.com/mbvlabs/andurel/pull/801))
2. Custom models ([#800](https://github.com/mbvlabs/andurel/pull/800))
3. Storage 0.8.1 ([#802](https://github.com/mbvlabs/andurel/pull/802))
4. `internal/runtime` ([#803](https://github.com/mbvlabs/andurel/pull/803))
5. `cmd/migrate` ([#804](https://github.com/mbvlabs/andurel/pull/804))

### Multiple hosts

Public hostnames live on `App.Hosts`. The default identity is `HOST_PRIMARY` in `.env` (`config.App.Domain` still holds that value). Extra hosts need a `routing.HostName` const in `config/hosts.go` plus a matching `App.Hosts` entry (load `HOST_<NAME>` yourself — not scaffolded). Missing hostnames fail at boot.

`andurel generate controller` and `andurel generate scaffold` take `--host` and `--prefix` independently:

- `--host` selects the Echo / `routing.HostName` (`primary` always works; others come from app `routing.HostName` consts).
- `--prefix` is only the path and package namespace.

Declare the host on the route with `routing.Host(config.HostAdmin)`. Controllers call `r.AddRoute(route, echo.Route{Method, Handler})`; Fx stays `RegisterRoutes(r)`. `AddRoute` fills `Path` and `Name` from the hosted route. Assets and `/api` stay on primary unless those routes set `Host(...)`.

Inertia apps must import `resources/js/routes.ts` after `andurel sync routes --json`. Helpers return full URLs once `configureRouteHosts` runs from the shared `hosts` prop at boot. Use `routes.name.path()` for a relative path.

See [Routing](/docs/head/routing), [Configuration](/docs/head/configuration), [controller](/docs/head/generate-controller), and [TypeScript Sync](/docs/head/inertia-typescript-sync).

### Custom models

`andurel generate model NAME --custom field:type ...` writes a struct-only model with `// andurel:custom`. There is no table, no migration, and no `sync model` / `sync factory`. Table-backed models stay the default (`generate model NAME` still reads migrations). After a migration changes columns, run `andurel sync model NAME` — not `generate model --update`.

See [model](/docs/head/generate-model) and [Models](/docs/head/models).

### Storage 0.8.1

`github.com/mbvlabs/andurel/pkg/storage` **v0.8.1** treats a queries file as present only when a `-- name:` annotation has SQL after it, not when the file is merely non-empty. Empty or comment-only query stubs no longer look like query sources. Run `andurel sync queries --json` after editing `models/queries/*.sql`.

See [Queries](/docs/head/queries).

### `internal/runtime`

Fx process graphs live in `internal/runtime`. `cmd/app` calls `runtime.App`; `cmd/queue` calls `runtime.Queue`. Register new services, controllers, and workers in their package `Module`. Add process-wide providers (email drivers, extra Fx provides) in `internal/runtime`.

`cmd/app` and `cmd/queue` stay thin `main` packages. `andurel generate` still wires new types into the right `fx.Module`.

See [Directory Structure](/docs/head/directory-structure), [Dependency Injection](/docs/head/dependency-injection), and [Request Lifecycle](/docs/head/request-lifecycle).

### `cmd/migrate`

`cmd/migrate` applies pending SQL with `storage.RunMigrations` on `migrations.Migrations`. It is a one-shot process, not part of `cmd/app` start. Development still uses `andurel db migrate up`. Production and CI should run the migrate binary (or image) before starting app and queue.

See [Migrations & Seeding](/docs/head/migrations), [db](/docs/head/db), and [Deployment](/docs/head/deployment).

## v2.0.0-alpha highlights

The first v2 alpha still defines the line:

- Independently versioned `pkg/*` modules (`storage`, `inertia`, `hypermedia`, `routing`, `server`, `email`, `validation`, `telemetry`, `kiks`) instead of copied `internal/` packages
- Fx composition with validated typed configuration providers
- Persistence through **pgx/v5** and **narsilc** (Bun and sqlc are gone)
- **Inertia v3** as the default UI (`andurel new --ui react/pnpm|vue/pnpm|svelte/pnpm`), with optional Templ/Datastar (`--ui templ/datastar`)
- Application-owned Inertia root document, Vite assets, and `cmd/ssr`
- Cookies and sessions via **kiks**
- Root `migrations/` and `seeds/`, plus `andurel.toml` / `andurel.lock`
- CLI groups: `generate`, `sync`, `inspect`, `db`, `packages`, `skill`, `tool`

First-alpha notes: [v2.0.0-alpha on GitHub](https://github.com/mbvlabs/andurel/releases/tag/v2.0.0-alpha).

Continue with the [Upgrade Guide](/docs/head/upgrade) if you are moving from v1 or the first alpha, or [Installation](/docs/head/installation) to create a new v2 application.

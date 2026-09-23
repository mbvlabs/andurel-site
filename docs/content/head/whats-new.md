# What's New in v2

Andurel **v2.0.0-alpha** is the first tagged pre-release of the v2 line. These head docs describe that release. Requires Go **1.27.1** and module path `github.com/mbvlabs/andurel/v2`.

## Highlights

- Independently versioned `pkg/*` modules (`storage`, `inertia`, `hypermedia`, `routing`, `server`, `email`, `validation`, `telemetry`, `kiks`) instead of copied `internal/` packages
- Fx composition roots with validated typed configuration providers
- Persistence through **pgx/v5** and **narsilc** (Bun and sqlc are gone)
- **Inertia v3** as the default UI (`andurel new --ui react/pnpm|vue/pnpm|svelte/pnpm`), with optional Templ/Datastar (`--ui templ/datastar`)
- Application-owned Inertia root document, Vite assets, and `cmd/ssr`
- Cookies and sessions via **kiks**
- Root `migrations/` and `seeds/`, plus `andurel.toml` / `andurel.lock`
- CLI groups: `generate`, `sync`, `inspect`, `db`, `packages`, `skill`, `tool`

## Full changelog

For every pull request and package bump in the alpha cut, see the GitHub release:

[v2.0.0-alpha on GitHub](https://github.com/mbvlabs/andurel/releases/tag/v2.0.0-alpha)

Continue with the [Upgrade Guide](/docs/head/upgrade) if you are moving from v1, or [Installation](/docs/head/installation) to create a new v2 application.

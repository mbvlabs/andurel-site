# Overview

The `andurel` command creates projects, generates application code, manages databases and tools, compiles email and frontend assets, builds releases, and exposes structured discovery.

## Discover the current surface

The v2 CLI evolves with each release, so discover commands instead of copying a static list:

```bash
andurel --help
andurel commands --json
andurel inspect project --json
andurel packages list --json
andurel tool list --json
```

## Daily workflow

```bash
andurel run
andurel fmt
andurel doctor --verbose
andurel build --version v2.0.0-alpha
```

`run` coordinates live reload and code generation. `fmt` formats Go and Templ. `doctor` checks project health, package migration issues, and generated code. `build` compiles narsilc when query files exist, email, Templ, CSS, optional Vite/SSR assets, and the Go application.

## Create a project

```bash
andurel new orbit --ui vue/pnpm
andurel new orbit --ui templ/datastar
andurel new orbit --dry-run --json
```

`--ui` accepts `react|vue|svelte` combined with `pnpm|bun|npm` (default `react/pnpm`), or `templ/datastar`. After creation, run `andurel tool sync` and `andurel doctor --json`.

## Generate and sync

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel generate scaffold Product
andurel generate model Product
andurel generate query ProductReport --table products
andurel sync queries
andurel generate controller Product
andurel generate job SendReceipt
andurel generate email Receipt
```

Keep derived TypeScript and Templ output current with `andurel sync routes`, `andurel sync payloads`, and `andurel sync views`.

## Preview mutations

Use structured dry runs where supported:

```bash
andurel new orbit --dry-run --diff --json
andurel generate scaffold Product --dry-run --diff --json
```

The v2 line does not automate upgrading a v1 application. See [Upgrade Guide](/docs/head/upgrade). Command groups: [generate](/docs/head/generate), [sync](/docs/head/sync), [inspect](/docs/head/inspect), [db](/docs/head/db), [packages](/docs/head/packages), [skill](/docs/head/skill), [Agent Output](/docs/head/agent-output).

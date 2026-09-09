# CLI Overview

The `andurel` command creates projects, generates application code, manages databases and tools, compiles email and frontend assets, builds releases, and exposes structured discovery.

## Discover the current surface

The master CLI is evolving toward v2, so discover commands instead of copying a static list:

```bash
andurel --help
andurel commands --json
andurel project info --json
andurel config show --json
```

## Daily workflow

```bash
andurel run
andurel fmt
andurel doctor --verbose
andurel build --version v2.0.0-dev
```

`run` coordinates live reload and code generation. `fmt` formats Go and Templ. `doctor` checks project health, package migration issues, and generated code. `build` compiles sqlc, email, Templ, CSS, optional Vite/SSR assets, and the Go application.

## Generate application code

```bash
andurel generate model Product
andurel generate query ProductReport --table products
andurel generate queries
andurel generate controller Product
andurel generate scaffold Product
andurel generate job SendReceipt
andurel generate email Receipt
```

## Preview mutations

Use structured dry runs where supported:

```bash
andurel new orbit --dry-run --diff --json
andurel generate scaffold Product --dry-run --diff --json
andurel extension add docker --dry-run --diff --json
```

The v2 line does not automate upgrading a v1 application. See [Moving to v2](/docs/latest/v2-migration).

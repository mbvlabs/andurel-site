# CLI Overview

The `andurel` command creates projects, generates application code, manages databases and tools, builds releases, and exposes structured discovery for automation.

## Discover commands

```bash
andurel --help
andurel commands --json
andurel project info --json
```

Use command discovery instead of relying on a copied command list when building tooling.

## Daily workflow

```bash
andurel run
andurel fmt
andurel doctor
andurel build
```

`andurel run` starts live reload. `fmt` formats Go and Templ. `doctor` checks project health. `build` compiles production assets and the Go binary.

## Preview mutations

Commands that change project files support dry-run output where applicable:

```bash
andurel generate scaffold Product --dry-run --diff --json
andurel extension add docker --dry-run --json
andurel upgrade --dry-run --diff --json
```

Review the reported created, updated, and deleted files before applying broad changes.

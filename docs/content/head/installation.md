# Installation

Install the Andurel v2 CLI, create an application, and start its development server. These head docs match **v2.0.0-alpha**.

## Requirements

You need:

- Go 1.27.1 or newer
- PostgreSQL
- Linux or macOS on amd64 or arm64
- Node.js and a package manager when using Inertia (the default); managed Inertia SSR requires Node.js 22 or newer

## Install the v2 CLI

v2 lives on the `/v2` module path. v1 remains on the unversioned path; that path’s `@latest` never jumps to v2.

```bash
# Documented alpha line (matches these head docs)
go install github.com/mbvlabs/andurel/v2@v2.0.0-alpha
andurel --version
```

While `v2.0.0-alpha` is the current `/v2` tag, `github.com/mbvlabs/andurel/v2@latest` resolves to the same line. Pin `@v2.0.0-alpha` when you want the exact release these docs describe.

Nightlies that track unreleased `master` ship as **prebuilt binaries** on the GitHub `nightly` release — do not use `go install …@master` for nightly metadata. v1 installs stay on `github.com/mbvlabs/andurel@v1.x.y`.

## Create an application

```bash
# Default: PostgreSQL + Inertia React + pnpm
andurel new orbit
cd orbit
andurel tool sync
cp .env.example .env
andurel db create
andurel db migrate up
andurel run
```

Edit `.env` before creating the database. The server is available at `http://localhost:8080` by default and reloads Go, Templ, CSS, narsilc output, and Vite as needed.

## Choose a frontend

The default project uses **Inertia React with pnpm**. Pass `--ui` for another adapter or for Templ/Datastar:

```bash
andurel new orbit --ui vue/bun
andurel new orbit --ui svelte/npm
andurel new orbit --ui templ/datastar
```

Supported Inertia adapters are Vue, React, and Svelte. Supported package managers are npm, pnpm, bun, and yarn. The choice is recorded in `andurel.lock` and drives installs, builds, and later generators.

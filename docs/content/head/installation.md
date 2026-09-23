# Installation

Install the CLI built from Andurel's `master` branch, create a v2 application, and start its development server.

## Requirements

You need:

- Go 1.27 or newer
- PostgreSQL
- Linux or macOS on amd64 or arm64
- Node.js and a package manager only when using Inertia; managed Inertia SSR requires Node.js 22 or newer

## Install the v2 development CLI

The stable `@latest` tag remains v1 until v2 is released. Install `master` to work with the version described by these docs:

```bash
go install github.com/mbvlabs/andurel@master
andurel --version
```

Because `master` is a development branch, pin a commit when reproducibility matters.

## Create an application

```bash
andurel new orbit
cd orbit
andurel tool sync
cp .env.example .env
andurel database create
andurel database migrate up
andurel run
```

Edit `.env` before creating the database. The server is available at `http://localhost:8080` by default and reloads Go, Templ, CSS, and generated sqlc code as needed.

## Choose a frontend

The default project uses Templ and Datastar. Select an Inertia adapter and, optionally, its JavaScript package manager when creating the project:

```bash
andurel new orbit --inertia vue
andurel new orbit --inertia react/pnpm
andurel new orbit --inertia svelte/bun
```

The supported adapters are Vue, React, and Svelte. The supported package managers are npm, pnpm, bun, and yarn.

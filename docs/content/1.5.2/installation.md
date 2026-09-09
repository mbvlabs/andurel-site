# Installation

Install the Andurel CLI, create an application, and start its development server.

## Requirements

You need:

- Go 1.26 or newer
- PostgreSQL
- Linux or macOS on amd64 or arm64

## Install Andurel

Install the latest stable CLI with Go:

```bash
go install github.com/mbvlabs/andurel@latest
andurel --version
```

For reproducible automation, install an explicit stable v1 tag instead of `latest`.

## Create an application

Create a project, synchronize its tools, configure the environment, and migrate the database:

```bash
andurel new orbit
cd orbit
andurel tool sync
cp .env.example .env
andurel database create
andurel database migrate up
andurel run
```

The development server is available at `http://localhost:8080` with live reload for Go, Templ, and CSS.

## Choose a frontend

The default project uses Templ and Datastar. Select an Inertia adapter when creating the project if you want a reactive frontend:

```bash
andurel new orbit --inertia vue
andurel new orbit --inertia react/pnpm
andurel new orbit --inertia svelte
```

See [Frontend Options](/docs/1.5.2/frontend-options) for the available adapters and JavaScript runtimes.

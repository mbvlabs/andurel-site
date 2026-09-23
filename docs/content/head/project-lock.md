# Project Lock

Andurel records scaffold choices and tool pins in two project files: `andurel.toml` and `andurel.lock`.

## andurel.toml

Human-editable project metadata for the framework tooling: UI choice, extensions, and related scaffold options. Treat it as the intent you want generators and upgrades to respect.

## andurel.lock

The lock records the framework version, scaffold choices, extensions, and pinned tools used to maintain the project. Examples include the JavaScript package manager for Vite builds and the Inertia SSR runtime expectations.

`go.mod` remains the source of truth for Go module versions. Updating the CLI does not rewrite application imports — use `andurel packages list` / `andurel packages update` deliberately. See [packages](/docs/head/packages).

## Why both exist

- **toml** — what you chose when creating or configuring the app
- **lock** — what the tooling last synchronized and which binaries/tools were pinned

Keep both under version control so agents and CI reproduce the same pad.

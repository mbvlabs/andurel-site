# Directory Structure

Andurel generates a layered application with conventional locations for HTTP, data, services, and presentation code.

## Application layers

```text
cmd/app/          application entry point and Fx wiring
config/           environment-backed configuration
controllers/      HTTP handlers and route registration
router/routes/    typed route declarations
views/            Templ components
models/           entities and persistence logic
services/         application workflows
queue/            jobs and workers
email/            email components
assets/           embedded compiled assets
database/         migrations, seeds, and database setup
```

## Generated framework code

The `internal` directory contains framework elements copied into the project, including routing, rendering, validation, storage contracts, and the HTTP server. The `andurel upgrade` command updates framework-owned files while preserving application-owned code.

Run a dry-run before upgrading:

```bash
andurel upgrade --dry-run --diff --json
```

## Project metadata

`andurel.lock` records the framework version, frontend adapter, JavaScript runtime, extensions, tools, and database conventions used by the project. Commit it to source control.

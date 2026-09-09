# Framework Packages

Andurel v2 is a project generator plus seven independently versioned Go modules. The generated application is not a thin wrapper around a hidden runtime: it owns its composition root, environment names, domain code, routes, controllers, models, and views. Framework packages provide reusable infrastructure at those boundaries.

## Package map

| Module | Responsibility | Deep dive |
| --- | --- | --- |
| `pkg/storage` | PostgreSQL, Bun/sqlc interop, transactions, migrations, test databases, and River clients | [Storage](/docs/latest/storage) |
| `pkg/inertia` | Echo-native Inertia v3 protocol, Vite assets, props, redirects, flash, and SSR | [Inertia](/docs/latest/inertia) |
| `pkg/hypermedia` | Templ rendering and Datastar/SSE response helpers | [Hypermedia](/docs/latest/hypermedia) |
| `pkg/routing` | Typed URL construction shared by Go and generated TypeScript | [Routing Package](/docs/latest/routing-package) |
| `pkg/server` | Bounded HTTP server construction and graceful shutdown | [Server](/docs/latest/server-package) |
| `pkg/email` | Transport-neutral email payloads, validation, retry classification, and Mailpit | [Email Package](/docs/latest/email-package) |
| `pkg/validation` | Structured field errors and composable validation rules | [Validation](/docs/latest/validation) |

The packages expose ordinary Go constructors rather than an Andurel service locator. Fx belongs to the generated application and connects constructors and process lifecycles.

## Composition and process boundaries

The web process normally follows this graph:

```text
config constructors
    -> storage.Postgres + queue insertion
    -> models and services
    -> Templ controllers or inertia.Renderer controllers
    -> Echo router
    -> server.Server
```

The queue process has a separate graph. It constructs the same storage connection, a queue processor, registered River workers, telemetry, and email transports, but no HTTP server. Configuration types preserve that separation: insertion and worker settings are distinct even though both eventually contain `storage.QueueConfig`.

## Configuration boundary

Package configuration types contain operational policy and validation. They do not choose environment-variable names. For example, `storage.DefaultConfig()` supplies pool defaults, while generated `config.NewDatabase()` maps `DB_*` values onto the type and validates the result.

This means packages work without Fx or environment variables, applications can adopt another configuration source without forking the framework, and a package update cannot silently introduce a new environment contract. See [Configuration](/docs/latest/configuration) for the generated mappings.

## Versioning and upgrades

Each package has its own module, semantic version, and changelog. Versions may differ from one another and from the CLI release. `go.mod` is the source of truth for code dependencies; `andurel.lock` records the framework version, scaffold choices, extensions, and tools used to maintain the project.

Updating the CLI does not rewrite existing imports. Update packages deliberately, read each changelog, compile generated code, and test the owning process. Storage changes require database, transaction, migration, queue-insertion, and worker coverage. Inertia changes require initial HTML, client visits, partial reloads, redirects, validation, Vite development, production assets, and configured SSR modes.

## Extending and replacing implementations

Depend on the smallest public interface available. Models normally accept `storage.Connection`, services accept email sender interfaces, and only Inertia controllers accept `*inertia.Renderer`. Application-specific adapters stay in the application.

Replace implementations at the composition root: tests can provide fake senders, a deployment can install a custom email transport, and a non-Fx program can construct storage directly. Importing implementation details throughout the application defeats these boundaries and makes independent upgrades harder.

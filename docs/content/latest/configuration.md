# Configuration

Andurel v2 keeps configuration in the application. The generated `config` package reads environment values once, validates them, and supplies independent typed values through Fx.

## Configuration providers

`config.Module` registers constructors for application, HTTP, session, authentication, database, queue insertion, queue workers, telemetry, email, and optional Inertia settings. Fx evaluates only the constructors needed by the active process.

```go
var Module = fx.Module("config",
    fx.Provide(
        NewApp,
        NewHTTP,
        NewSession,
        NewDatabase,
        NewQueueInsert,
        NewQueueWorker,
        NewTelemetry,
        NewMail,
        NewMailTransport,
    ),
)
```

Framework packages define useful configuration types and defaults, but do not read your environment. The application decides how values are loaded and overridden.

## Core environment

```dotenv
ENVIRONMENT=development
PROJECT_NAME=orbit
DOMAIN=localhost:8080
PROTOCOL=http
HOST=localhost
PORT=8080

DB_KIND=postgres
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=orbit_development
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable
```

Database pool, connection timeout, HTTP timeout, queue, email, and telemetry defaults can also be overridden. See the generated files in `config/` for the authoritative variables and validation rules for your scaffold.

## Secrets and security

`SESSION_KEY`, `SESSION_ENCRYPTION_KEY`, `TOKEN_SIGNING_KEY`, and `PEPPER` are generated for a new project. Never reuse development values in production or commit `.env`.

The default CSRF strategy is `header_only`. Additional CORS and CSRF origins must be explicit; credentialed wildcard origins are rejected during startup.

## Fail fast

Each constructor returns an error when parsing or validation fails. Fx prevents dependent components from starting, so an invalid duration, port, key, database setting, queue setting, or Inertia SSR configuration fails at the composition root instead of surfacing during a request.

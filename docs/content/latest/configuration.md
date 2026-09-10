# Configuration

Andurel keeps configuration in the generated application. The `config` package loads local `.env`, parses process environment values into typed structs, validates them, and supplies them through Fx. Framework packages receive typed values and never read environment variables themselves.

## Loading, parsing, and precedence

`config.LoadEnvironment()` runs before Fx construction. A constructor default applies only when a key is absent, so an explicitly empty string differs from an unset value. Parsers support required and optional strings, integers, booleans, floats, Go durations, and comma-separated lists.

Durations use values such as `250ms`, `10s`, `5m`, or `1h30m`. Lists are split on commas, trimmed, and empty items removed. Each constructor joins conversion and semantic errors and returns context such as `config: database: ...`.

## Typed providers and process scope

Providers are independent: `NewApp`, `NewAuth`, `NewHTTP`, `NewSession`, `NewDatabase`, `NewQueueInsert`, `NewQueueWorker`, `NewTelemetry`, `NewMail`, `NewMailTransport`, and optional `NewInertia`.

Fx calls only providers used by a process. Web needs HTTP and queue insertion; queue needs worker policy but not HTTP. Keep types distinct so a process cannot accidentally acquire a capability it should not own.

## Application and HTTP

| Variable | Default | Meaning |
| --- | --- | --- |
| `ENVIRONMENT` | `development` | Runtime policy; `production` enables generated production behavior |
| `PROJECT_NAME` | project name | Application and telemetry identity |
| `DOMAIN` | `localhost:8080` | Public host, optionally with port |
| `PROTOCOL` | `http` when empty | Public HTTP or HTTPS scheme |
| `HOST` | `localhost` | Local listen host, not public domain |
| `PORT` | `8080` | Listen port, 1 through 65535 |
| `HTTP_IDLE_TIMEOUT` | `120s` | Keep-alive idle bound; zero disables |
| `HTTP_READ_TIMEOUT` | `10s` | Request-read bound; zero disables |
| `HTTP_WRITE_TIMEOUT` | `30s` | Response-write bound; consider SSE duration |
| `CORS_ALLOWED_ORIGINS` | empty | Comma-separated CORS origins |
| `CSRF_STRATEGY` | `header_only` | Or `header_or_legacy_token` compatibility mode |
| `CSRF_TRUSTED_ORIGINS` | empty | Explicit additional trusted origins |

Behind a proxy, a normal split is `DOMAIN=app.example.com`, `PROTOCOL=https`, `HOST=0.0.0.0`, and an internal port. Credentialed wildcards are rejected.

## Database

| Variable | Default | Meaning |
| --- | --- | --- |
| `DB_KIND` | `postgres` | `postgres` or `postgresql` |
| `DB_HOST` / `DB_PORT` | `127.0.0.1` / `5432` | PostgreSQL address |
| `DB_NAME` | project database | Database name |
| `DB_USER` / `DB_PASSWORD` | `postgres` / `postgres` | Credentials |
| `DB_SSL_MODE` | `disable` | pgx SSL mode through `verify-full` |
| `DB_APPLICATION_NAME` | project name | PostgreSQL activity identity |
| `DB_CONNECT_TIMEOUT` | `5s` | Initial connection timeout |
| `DB_STATEMENT_CACHE_CAPACITY` | `512` | Statement-cache entries; zero disables |
| `DB_DESCRIPTION_CACHE_CAPACITY` | `512` | Description-cache entries; zero disables |
| `DB_MAX_OPEN_CONNECTIONS` | `25` | Open pool limit |
| `DB_MAX_IDLE_CONNECTIONS` | `25` | Idle pool limit |
| `DB_CONNECTION_MAX_LIFETIME` | `1h` | Reusable connection lifetime |
| `DB_CONNECTION_MAX_IDLE_TIME` | `30m` | Idle connection lifetime |
| `DB_OPEN_TELEMETRY` | `true` | pgx instrumentation |

Budget connections across all web and queue replicas, leaving room for migrations and administration. See [Storage](/docs/latest/storage).

## Sessions and authentication

| Variable | Required/default | Meaning |
| --- | --- | --- |
| `SESSION_KEY` | required | Hex cookie authentication key |
| `SESSION_ENCRYPTION_KEY` | required | Hex cookie encryption key |
| `SESSION_MAX_AGE` | `604800` | Positive lifetime in seconds |
| `TOKEN_SIGNING_KEY` | required | Generated account-token signing key |
| `PEPPER` | required | Current password pepper |
| `PREVIOUS_PEPPERS` | empty | Comma-separated rotation fallback values |

Production values belong in deployment secrets, never committed `.env`. Pepper rotation moves the old current value into `PREVIOUS_PEPPERS` for a bounded migration window. Session and signing-key rotation may invalidate live cookies or tokens. The cookie name includes project and environment to avoid cross-environment collisions.

## Queue configuration

Insertion reads only `QUEUE_MAX_ATTEMPTS` and `QUEUE_SCHEMA`. Workers clone that base and add processing policy:

| Variable | Default | Meaning |
| --- | --- | --- |
| `QUEUE_MAX_WORKERS` | `100` | Default queue concurrency |
| `QUEUE_ADVISORY_LOCK_PREFIX` | `0` | River lock namespace |
| `QUEUE_CANCELLED_JOB_RETENTION_PERIOD` | `24h` | Cancelled retention |
| `QUEUE_COMPLETED_JOB_RETENTION_PERIOD` | `24h` | Completed retention |
| `QUEUE_DISCARDED_JOB_RETENTION_PERIOD` | `168h` | Discarded retention |
| `QUEUE_FETCH_COOLDOWN` / `QUEUE_POLL_INTERVAL` | River defaults | Fetch timing; poll cannot be shorter than cooldown |
| `QUEUE_ID` | empty | Optional client identifier |
| `QUEUE_JOB_CLEANER_TIMEOUT` | `30s` | Cleaner timeout |
| `QUEUE_JOB_STUCK_THRESHOLD` / `QUEUE_JOB_TIMEOUT` | River defaults | Stuck and execution bounds |
| `QUEUE_REINDEXER_TIMEOUT` | `1m` | Reindex bound |
| `QUEUE_RESCUE_STUCK_JOBS_AFTER` | `1h` | Must not be shorter than job timeout |
| `QUEUE_SOFT_STOP_TIMEOUT` | River default | Shutdown grace period |
| `QUEUE_POLL_ONLY` | `false` | Disable notification-assisted fetching |
| `QUEUE_SKIP_JOB_KIND_VALIDATION` | `false` | Skip registered-kind validation |
| `QUEUE_SKIP_UNKNOWN_JOB_CHECK` | `false` | Skip unknown stored-kind check |

Edit `NewQueueWorker` for named queues or per-queue concurrency. Validation covers names, schemas, ranges, durations, and relationships.

## Email and telemetry

`DEFAULT_SENDER_SIGNATURE` defaults to `noreply@DOMAIN`. `EMAIL_PROVIDER` selects an installed transport. Development defaults to `mailpit`, with `MAILPIT_HOST=0.0.0.0` and `MAILPIT_PORT=1025`. Production requires an installed provider such as the `aws-ses` extension and its own credential configuration.

| Telemetry variable | Default |
| --- | --- |
| `TELEMETRY_SERVICE_NAME` | slugged project name |
| `TELEMETRY_SERVICE_VERSION` | `1.0.0` |
| `OTLP_LOGS_ENDPOINT`, `OTLP_METRICS_ENDPOINT`, `OTLP_TRACES_ENDPOINT` | empty |
| `OTLP_HEADERS` | empty |
| `TRACE_SAMPLE_RATE` | `1.0`, valid 0 through 1 |
| `TELEMETRY_BATCH_SIZE` | `512` |
| `TELEMETRY_BATCH_TIMEOUT_MS` | `5000` |

## Inertia and SSR

| Variable | Default | Meaning |
| --- | --- | --- |
| `INERTIA_CONTAINER_ID` | `app` | DOM mount ID |
| `INERTIA_VITE_DEV_URL` | `http://localhost:5173/assets/dist` | Development assets |
| `INERTIA_ENTRY_POINT` | `resources/js/app.ts` or `.tsx` | Vite entry |
| `INERTIA_PROTOCOL_DEBUG` | `true` in generated apps | Safe metadata diagnostics |
| `INERTIA_SSR_RUNTIME` | `node` | `cmd/ssr` executable |
| `INERTIA_SSR_BUNDLE` | `assets/dist/ssr/ssr.js` | `cmd/ssr` bundle |
| `INERTIA_SSR_LISTEN` | `http://127.0.0.1:13714` | Node bind URL (IP or localhost) |
| `INERTIA_SSR_URL` | `http://127.0.0.1:13714` | Where `cmd/app` POSTs `/render` |
| `INERTIA_SSR_STARTUP_TIMEOUT` | `10s` | Health deadline for `cmd/ssr` |
| `INERTIA_SSR_REQUEST_TIMEOUT` | `2s` | Render deadline |
| `INERTIA_SSR_MAX_RESPONSE_BYTES` | `2097152` | Response-size bound |
| `INERTIA_SSR_MINIMUM_MAJOR` | `22` | Minimum Node major |
| `INERTIA_SSR_FAIL_FAST` | `true` in generated apps | Error instead of client fallback |

`INERTIA_SSR_LISTEN` is the address Node binds. `INERTIA_SSR_URL` is the HTTP client URL and may be a service hostname. Pages still opt into SSR individually. See [Inertia](/docs/latest/inertia) and [SSR](/docs/latest/inertia-ssr).

## Adding application settings

Add a focused type and constructor, validate it, register the constructor in `config.Module`, and inject the type only where needed:

```go
type Search struct { Endpoint string; Timeout time.Duration }

func NewSearch() (Search, error) {
    env := newEnvironment()
    cfg := Search{
        Endpoint: env.RequiredString("SEARCH_ENDPOINT"),
        Timeout: env.Duration("SEARCH_TIMEOUT", 2*time.Second),
    }
    if err := env.Err(); err != nil { return Search{}, err }
    if cfg.Timeout <= 0 { return Search{}, errors.New("SEARCH_TIMEOUT must be positive") }
    return cfg, nil
}
```

Provide a `Clone` method for mutable maps or slices. Startup errors should remain fatal; continuing with partial defaults converts a clear deployment failure into request-time corruption or insecure behavior.

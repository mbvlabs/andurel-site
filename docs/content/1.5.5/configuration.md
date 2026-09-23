# Configuration

Andurel applications read configuration from environment variables. Copy `.env.example` to `.env` for local development and keep secrets out of source control.

## Application settings

The generated application configures its host, port, domain, session lifetime, token signing key, CORS origins, and CSRF strategy through the `config` package.

```dotenv
HOST=localhost
PORT=8080
DOMAIN=localhost:8080
PROTOCOL=http
PROJECT_NAME=orbit
```

## Database settings

Configure PostgreSQL before running database commands or starting modules that use models and queues.

```dotenv
DB_KIND=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=orbit_development
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable
```

## Security settings

Generate unique values for session authentication, session encryption, and token signing. Production cookies are secure when the application runs in the production environment.

Never reuse development secrets in production. Restrict CORS origins to known application origins and keep CSRF protection enabled for browser requests.

## Telemetry settings

Logs work by default. Set OTLP endpoints to export metrics and traces:

```dotenv
OTLP_METRICS_ENDPOINT=
OTLP_TRACES_ENDPOINT=
TRACE_SAMPLE_RATE=1.0
```

# Telemetry

Generated applications include structured logging, metrics, traces, and HTTP instrumentation through OpenTelemetry.

## Process-local wiring

Both `cmd/app` and `cmd/queue` include the application-owned telemetry module. Fx starts and stops exporters with each process. Application and request logs use `log/slog`; pass `context.Context` so trace metadata can be attached.

## Export configuration

```dotenv
TELEMETRY_SERVICE_NAME=orbit
TELEMETRY_SERVICE_VERSION=1.0.0
OTLP_LOGS_ENDPOINT=
OTLP_METRICS_ENDPOINT=https://collector.example.com
OTLP_TRACES_ENDPOINT=https://collector.example.com
OTLP_HEADERS=Authorization=Bearer token
TRACE_SAMPLE_RATE=1.0
```

Batch size and timeout are configurable. The generated provider validates service identity, batch settings, and the trace sample rate before exporters start.

## Instrument useful boundaries

Create spans around meaningful application operations and record stable route templates rather than raw, high-cardinality values. Never put passwords, tokens, complete email bodies, or unbounded user input into logs or telemetry attributes.

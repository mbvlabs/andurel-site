# Telemetry

Generated applications include structured logging, metrics, traces, and HTTP instrumentation through OpenTelemetry.

## Logs

Application and request logs use `log/slog`. Include request context when available so trace metadata can be attached.

## Metrics and traces

Set OTLP endpoints to export telemetry:

```dotenv
OTLP_METRICS_ENDPOINT=https://collector.example.com
OTLP_TRACES_ENDPOINT=https://collector.example.com
OTLP_HEADERS=Authorization=Bearer token
```

Control trace volume with `TRACE_SAMPLE_RATE`.

## Application spans

Create spans around meaningful application operations, not every helper call. Record stable route templates and status codes rather than high-cardinality raw values.

## Sensitive data

Do not place passwords, tokens, email bodies, or unbounded user input in logs, metric attributes, or span attributes.

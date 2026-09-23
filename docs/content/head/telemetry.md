# Telemetry

Generated applications include structured logging, metrics, traces, and HTTP instrumentation through OpenTelemetry via `github.com/mbvlabs/andurel/pkg/telemetry`.

## Process-local wiring

Both `cmd/app` and `cmd/queue` include the application-owned telemetry module. Fx starts and stops exporters with each process. Composition roots call `telemetry.New` with service identity and `With*` options from `config.Telemetry`. Do not call `otel.SetTracerProvider` or `slog.SetDefault`; use the injected `*telemetry.Telemetry` and package helpers that read it from context.

```go
tel, err := telemetry.New(ctx, cfg.ServiceName, cfg.ServiceVersion,
    telemetry.WithConsole(),
    telemetry.WithLogLevel(slog.LevelInfo),
    telemetry.WithOTLPTraces(cfg.OTLPTracesEndpoint, cfg.OTLPHeaders),
    telemetry.WithOTLPMetrics(cfg.OTLPMetricsEndpoint, cfg.OTLPHeaders),
    telemetry.WithOTLPLogs(cfg.OTLPLogsEndpoint, cfg.OTLPHeaders),
    telemetry.WithTraceSampleRate(cfg.TraceSampleRate),
    telemetry.WithBatchConfig(cfg.BatchSize, cfg.BatchTimeout, cfg.QueueSize),
)
```

Attach the provider to request or job contexts so helpers resolve it:

```go
ctx = tel.Context(ctx)
```

## Export configuration

```dotenv
TELEMETRY_SERVICE_NAME=orbit
TELEMETRY_SERVICE_VERSION=1.0.0
OTLP_LOGS_ENDPOINT=
OTLP_METRICS_ENDPOINT=https://collector.example.com
OTLP_TRACES_ENDPOINT=https://collector.example.com
OTLP_HEADERS=Authorization=Bearer token
TRACE_SAMPLE_RATE=1.0
TELEMETRY_BATCH_SIZE=512
TELEMETRY_BATCH_TIMEOUT_MS=5000
```

The generated provider validates service identity, batch settings, and the trace sample rate (0 through 1) before exporters start. Empty OTLP endpoints disable that signal export. See [Configuration](/docs/head/configuration).

## Spans in controllers and services

Controllers start a child span from Echo and write it back onto the request:

```go
func (s Sessions) Create(etx *echo.Context) error {
    ctx, span := telemetry.From(etx, "sessions.create")
    defer span.End()

    user, err := s.identity.AuthenticateUser(ctx, data)
    if err != nil {
        telemetry.Error(ctx, "failed to authenticate user", "error", err)
        return err
    }

    telemetry.Set(ctx, "user.id", user.ID.String())
    telemetry.Info(ctx, "session created", "user.id", user.ID.String())
    return nil
}
```

Services use `Start` with slog-shaped attributes:

```go
ctx, span := telemetry.Start(ctx, "products.archive", "product.id", id)
defer span.End()

if err := p.repo.Archive(ctx, id); err != nil {
    return telemetry.Fail(ctx, err)
}
```

| Helper | Role |
| --- | --- |
| `From(etx, name)` | Span from `*echo.Context`; updates request context |
| `Start(ctx, name, key, val, ...)` | Span with attributes |
| `Set(ctx, key, val, ...)` | Attributes on the current span |
| `Fail(ctx, err)` | Record error + error status; returns `err` |
| `Debug` / `Info` / `Warn` / `Error` | Log through the telemetry pipeline and add a span event |

Prefer stable route templates and identifiers over raw URLs or unbounded user input. Never put passwords, tokens, complete email bodies, or secrets into attributes. pgx instrumentation follows `DB_OPEN_TELEMETRY` when enabled on the storage pool.

# Queues

Andurel v2 separates job insertion from job processing while sharing the PostgreSQL connection managed by `pkg/storage`.

## Generate a job

```bash
andurel generate job RebuildSearchIndex
```

Job arguments live in `queue/jobs`; workers live in `queue/workers` and are registered through `queue.Module`.

## Process boundaries

The web application includes only the queue insertion module (`storage.NewQueueInsert` as `storage.InsertQueue`), so controllers and services can enqueue work without starting River workers. `cmd/queue` is a dedicated Fx application that constructs `storage.NewQueueProcessor`, workers, periodic jobs, telemetry, and email transports.

```go
_, err := queue.Insert(ctx, jobs.RebuildSearchIndexArgs{ID: id}, nil)

err := storage.RunInTransaction(ctx, db,
    func(ctx context.Context, tx storage.Transaction) error {
        if err := products.WithTx(tx).MarkQueued(ctx, id); err != nil {
            return err
        }
        _, err := queue.InsertTx(ctx, tx, jobs.RebuildSearchIndexArgs{ID: id}, nil)
        return err
    },
)
```

Run the web and queue executables as separate processes in production. Each owns its lifecycle and graceful shutdown. See [Request Lifecycle](/docs/head/request-lifecycle).

## Configuration

`config.QueueInsert` contains settings needed to enqueue jobs (`QUEUE_MAX_ATTEMPTS`, `QUEUE_SCHEMA`). `config.QueueWorker` clones that base and adds worker count, queues, retention, polling, timeouts, advisory lock, and operational policy. Invalid settings fail during startup. See [Configuration](/docs/head/configuration).

## Reliable workers

Jobs may be retried. Make workers idempotent, load current records inside the worker when appropriate, keep arguments small, and return classified errors so River and email delivery can apply the right retry policy. Use a shared storage transaction when a durable state change and a job insert must commit together.

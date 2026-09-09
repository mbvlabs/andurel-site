# Queues

Andurel uses River to run durable PostgreSQL-backed jobs and workers.

## Generate a job

```bash
andurel generate job RebuildSearchIndex
```

A job defines its arguments and a worker defines the operation. Register workers through the queue module so the processor can dispatch them.

## Enqueue work

Controllers and services depend on the narrow queue insertion interface. Enqueue only the data required to perform the work, then load current records inside the worker when appropriate.

## Delivery guarantees

Jobs can be retried. Make workers safe to run more than once, record durable outcomes transactionally where needed, and return errors so River can apply its retry policy.

## Operations

The queue shares the application database. Apply River migrations and ensure the queue processor is running in processes expected to execute jobs.

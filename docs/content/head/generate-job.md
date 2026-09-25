# job

`andurel generate job` creates a River background job: argument types under `queue/jobs/` and worker registration in the queue package.

## When to use

- Add a new background job plus worker wiring.
- Do not confuse with `andurel inspect jobs`, which only lists existing job files.

## Usage

```bash
andurel generate job NAME [flags]
```

`NAME` is CamelCase (for example `SendWelcomeEmail` or `SendReceipt`).

## Flags

| Flag | What it does |
| --- | --- |
| `--queue` | Assign the job to a specific River queue name |
| `--dry-run` / `--diff` | Preview without writing; see [generate](/docs/head/generate) |

## Examples

```bash
andurel generate job SendWelcomeEmail --dry-run --json
andurel generate job SendReceipt
andurel generate job ProcessInvoice --queue billing
```

## Next

```bash
andurel inspect jobs --json
```

See [Queues](/docs/head/queues) for inserting jobs from the web process and running dedicated workers.

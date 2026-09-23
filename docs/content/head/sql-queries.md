# SQL Queries

Andurel v2 persists through **narsilc**: you write annotated SQL under `models/queries/`, and the tool generates Go that fills application-owned result structs.

## Generate a query file

```bash
andurel generate query UserReport
andurel generate query UserReport --table users
```

The command creates `models/queries/user_report.sql`. Without `--table`, the file contains commented annotation examples. With `--table`, it includes a starter annotated query for the existing table.

After editing the SQL, compile Go code:

```bash
andurel sync queries
```

Output is written under `models/internal/queries`. The command is a no-op when `models/queries` contains no annotated SQL files. `andurel run`, `andurel build`, scaffolding, and extension application also regenerate narsilc output when query files exist. narsilc reads configuration from `andurel.lock`.

## Keep generated queries behind models

Only model packages should import the internal generated query package. Controllers and services depend on application model APIs and projection types, not database-shaped rows or parameters.

```go
generated := queries.New(products.db)
rows, err := generated.ListProductReport(ctx)
```

Map generated rows to an application-owned projection inside the model method. This preserves the model layer as the public persistence boundary.

## Share transactions

narsilc clients accept `storage.Connection` and `storage.Transaction` directly because both implement pgx `DBTX`. During a transaction:

```go
err := storage.RunInTransaction(ctx, connection,
    func(ctx context.Context, tx storage.Transaction) error {
        return queries.New(tx).RecordProductCreated(ctx, productID)
    },
)
```

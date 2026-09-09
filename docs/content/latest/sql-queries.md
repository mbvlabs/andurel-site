# SQL Queries

Andurel v2 keeps Bun as the default for ordinary model CRUD and includes sqlc for queries that are clearer as explicit SQL.

## Generate a query file

```bash
andurel generate query UserReport
andurel generate query UserReport --table users
```

The command creates `models/queries/user_report.sql`. Without `--table`, the file contains commented annotation examples. With `--table`, it includes a starter annotated query for the existing table.

After editing the SQL, generate Go code:

```bash
andurel generate queries
```

Output is written under `models/internal/queries`. The command is a no-op when `models/queries` contains no SQL files. `andurel run`, `andurel build`, scaffolding, and extension application also regenerate sqlc output when query files exist.

## Keep sqlc behind models

Only model packages should import the internal generated query package. Controllers and services depend on application model APIs and projection types, not database-shaped sqlc rows or parameters.

```go
generated := queries.New(products.db.DB())
rows, err := generated.ListProductReport(ctx)
```

Map generated rows to an application-owned projection inside the model method. This preserves the model layer as the public persistence boundary.

## Share transactions

The sqlc configuration uses the same standard SQL pool exposed by `storage.Connection`. During a transaction, bind generated queries with `queries.WithTx(tx.SQL())` while Bun-backed model methods use `tx.Executor()`.

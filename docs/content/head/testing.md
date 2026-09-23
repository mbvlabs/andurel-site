# Testing

Andurel v2 reworks application testing around explicit constructors, `storage.Connection`, and isolated PostgreSQL databases.

## Database tests

Use `storage.NewTestCluster` to start a shared PostgreSQL 17 Alpine container (via testcontainers) and create isolated databases per test. Prefer one cluster per package over a container per case.

```go
func TestMain(m *testing.M) {
    ctx := context.Background()
    cluster, err := storage.NewTestCluster(ctx)
    if err != nil {
        panic(err)
    }
    defer cluster.Close(ctx)

    testCluster = cluster
    os.Exit(m.Run())
}

func TestProductsCreate(t *testing.T) {
    db := testCluster.NewTestDB(t, migrations.FS, "migrations")
    product, err := factories.CreateProduct(ctx, db,
        factories.WithProductName("Orbit"),
    )
    // ...
}
```

`NewTestDB` creates a database, applies Goose migrations from an `fs.FS`, and registers cleanup. Generated projects embed SQL under root `migrations/`. See [Getting Started](/docs/head/database).

## Unit and HTTP tests

Construct collaborators through ordinary Go constructors or Fx test apps. Prefer fakes at interface boundaries (email senders, clocks) over global state.

Exercise Echo handlers with the same middleware stack the web process uses when behavior depends on kiks cookies or Inertia headers. Inject a real `*kiks.Jar` (or a focused test jar) rather than stuffing session values into ambient context.

## Factories

Build entities with `models/factories` and keep them synchronized with the model Entity:

```bash
andurel sync factories --check --json
```

```go
products, err := factories.CreateProducts(ctx, db, 3,
    factories.WithProductActive(true),
)
```

See [Factories](/docs/head/factories).

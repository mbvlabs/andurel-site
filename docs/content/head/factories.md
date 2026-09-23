# Factories

Factories live under `models/factories` and build model entities for tests and seeds. Andurel-owned declarations stay synchronized with the model Entity; custom helpers belong outside generated blocks.

## Generate and sync

Model generation creates a matching factory unless you pass `--skip-factory`. Refresh factories after Entity changes:

```bash
andurel sync factory Product --check
andurel sync factory Product --sync
andurel sync factories --check --json
andurel sync factories --sync
```

`sync factory` refreshes one factory. `sync factories` covers every managed factory declaration and requires `--check` or `--sync`. Prefer `--check` in CI.

## Build and create

Generated factories expose in-memory builders and persistence helpers:

```go
product := factories.BuildProduct(
    factories.WithProductName("Orbit"),
    factories.WithProductPriceCents(1200),
)

created, err := factories.CreateProduct(ctx, db,
    factories.WithProductName("Orbit"),
    factories.WithProductActive(true),
)
if err != nil {
    return err
}

batch, err := factories.CreateProducts(ctx, db, 5,
    factories.WithProductActive(true),
)
```

`BuildProduct` leaves auto-managed fields (ID, timestamps) at zero. `CreateProduct` persists through the model API and returns the entity with database-assigned values. Option functions (`WithProductName`, …) override defaults from go-faker / sensible stubs.

## Seeds

Named seed compositions in root `seeds/` should call factories and exported model APIs, then run through `andurel db seed`. See [Migrations & Seeding](/docs/head/migrations) and [Testing](/docs/head/testing).

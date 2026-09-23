# Migrations & Seeding

SQL migrations under root `migrations/` are the schema source of truth. Named seed compositions live in root `seeds/`.

## Create and apply migrations

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel db migrate status
```

Generated projects embed migration SQL so CLI commands, tests, and production binaries share one history. narsilc output does not replace migrations.

## Seeds

```bash
andurel db seed
andurel db seed Demo
```

Seeds should call exported model and factory APIs. See [Factories](/docs/head/factories) and [db](/docs/head/db) for lifecycle commands such as `rebuild`.

# Database

Andurel uses PostgreSQL, pgx, Bun, and SQL migrations. Models and factories are generated from your migration schema.

## Create and migrate

```bash
andurel database create
andurel database migrate up
andurel database migrate status
```

Create a migration before generating a model:

```bash
andurel database migrate new create_products_table
andurel generate model Product
```

The table needs a supported primary key. UUID, serial, bigserial, and string primary keys are supported.

## Models

Generated models contain an entity, create and update data types, and persistence methods built on Bun. Use services when an operation coordinates multiple models or side effects.

## Factories

Factories follow model fields and support tests and seeds. Check for drift without writing files:

```bash
andurel generate factories --check --json
```

Use `--sync` to update generated factory declarations while preserving custom helpers with non-conflicting names.

## Rebuild locally

Recreate the development database, apply migrations, and run seeds with one command:

```bash
andurel database rebuild
```

Destructive commands prompt by default. Reserve `--force` for explicit automation.

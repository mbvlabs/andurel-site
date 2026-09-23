# generate

Andurel generates application-owned code that you can edit after creation. Framework behavior is imported from versioned packages; project policy and domain code remain local.

## Generate a resource

Start with the canonical SQL migration, then generate CRUD:

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel generate scaffold Product
```

The scaffold includes a constructed model API, entity and data types, a synchronized factory, controller, typed routes, and views. Generation follows the UI in `andurel.lock` (Inertia pages or Templ). Pass `--api` for JSON handlers. See [Generators](/docs/head/inertia-generators) for the Inertia payload, TypeScript, and route-helper flow.

## Generate individual parts

```bash
andurel generate model Product
andurel generate model Product --mode read-only
andurel sync factory Product --sync
andurel generate controller Product index show
andurel generate controller Dashboard overview --model-name User
andurel generate query SalesReport
andurel generate job SendReceipt
andurel generate email Receipt
```

Model modes are `crud`, `read-only`, and `create-only`. Model generation reads migration history, requires a supported primary key for CRUD operations, and updates `models.Module` so Fx can supply the new plural model API.

## Keep generated state current

```bash
andurel sync factories --check --json
andurel sync factories --sync --json
andurel sync views
andurel sync queries
andurel sync routes
andurel sync payloads
```

Use `--dry-run --diff --json` before broad mutations. narsilc output under `models/internal/queries` and Templ-generated Go are derived files; models, factories, controllers, SQL query files, routes, and authored views are application code.

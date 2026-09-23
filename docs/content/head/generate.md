# generate

Andurel generates application-owned code that you can edit after creation. Framework behavior is imported from versioned packages; project policy and domain code remain local.

## Generate a resource

Start with the canonical SQL migration, then generate CRUD:

```bash
andurel generate migration create_products_table
andurel db migrate up
andurel generate scaffold Product
andurel generate scaffold Product --dry-run --diff --json
```

The scaffold includes a constructed model API, entity and data types, a synchronized factory, controller, typed routes, and views. Generation follows the UI recorded in `andurel.toml` (`project.inertia`, or Templ when unset). See [Generators](/docs/head/inertia-generators) for the Inertia payload, TypeScript, and route-helper flow.

## Generate individual parts

```bash
andurel generate model Product
andurel generate model Product --mode read-only
andurel generate model Product --update --yes
andurel sync factory Product --sync
andurel generate controller Product index show
andurel generate controller Dashboard overview --model-name User
andurel generate query SalesReport --table products
andurel generate job SendReceipt
andurel generate email Receipt
```

| Generator | Creates |
| --- | --- |
| `migration` | SQL file under root `migrations/` |
| `model` | Model API, entity, `models/queries/*.sql`, optional factory |
| `query` | Annotated SQL under `models/queries/` |
| `controller` | Controller, views/pages, route wiring |
| `scaffold` | Model + controller + views + routes together |
| `job` | Job args and worker registration |
| `email` | Authored email template |

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

Use `--dry-run --diff --json` before broad mutations. narsilc output under `models/internal/queries` and Templ-generated Go are derived files; models, factories, controllers, SQL query files, routes, and authored views are application code. See [sync](/docs/head/sync).

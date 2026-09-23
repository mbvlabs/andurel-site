# Code Generation

Andurel generates application-owned code that you can edit after creation.

## Generate a resource

Start with a migration, then generate a complete CRUD resource:

```bash
andurel database migrate new create_products_table
andurel generate scaffold Product
```

The scaffold includes a model, factory, controller, typed routes, and views. Pass `--inertia` for pages using the project's configured Inertia adapter, or `--api` for JSON handlers.

## Generate individual parts

```bash
andurel generate model Product
andurel generate factory Product
andurel generate controller Product index show
andurel generate controller Dashboard overview
andurel generate job SendReceipt
andurel generate email Receipt
```

## Preview first

Use `--dry-run --diff --json` before a large generation or update. Generated resource code belongs to the application after creation, while framework-owned internal files are managed by `andurel upgrade`.

## Refresh views and routes

```bash
andurel generate view
andurel generate routes
```

The route command generates TypeScript helpers only for Inertia projects.

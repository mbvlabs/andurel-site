# controller

`andurel generate controller` creates a controller, views or Inertia pages, and typed routes. Pass the controller name in CamelCase and optional actions.

## When to use

- Add an HTTP resource when the model already exists.
- Use [scaffold](/docs/head/generate-scaffold) when you also need the model.
- Use `--api` for JSON under `controllers/api`; otherwise generation follows the project UI.

## Usage

```bash
andurel generate controller NAME [action action ...] [flags]
```

When no actions are provided, the standard CRUD set is generated. When one or more of `index`, `show`, `new`, `create`, `edit`, `update`, or `destroy` are listed, only those resource actions are created. Partial CRUD views only link to companion actions that are also present.

Non-CRUD actions become empty controller methods, matching empty Templ components or Inertia pages, and conventional `GET /<resource>/<action>` routes.

Names may include one lowercase namespace segment (`admin/Widget`). Namespaced controllers live under `controllers/admin`, use `admin.*` route names, and Admin-prefixed symbols.

## Flags

| Flag | What it does |
| --- | --- |
| `--api` | Generate a JSON API controller under `controllers/api` (no views; default actions exclude `new`/`edit`) |
| `--model-name` | Back the controller with a different existing model name |
| `--dry-run` / `--diff` | Preview without writing; see [generate](/docs/head/generate) |

With `--api`, any namespace segment in the name nests under `api`.

## Examples

```bash
andurel generate controller Product
andurel generate controller Product index show --dry-run --json
andurel generate controller Dashboard overview
andurel generate controller Dashboard --model-name User
andurel generate controller admin/Widget export
andurel generate controller Users --api
```

In Inertia projects, pages use the adapter in `andurel.lock`. Templ/Datastar projects get Templ views instead. See [Generators](/docs/head/inertia-generators) for payload structs, `FromStruct`, and TypeScript helpers.

## Next

```bash
andurel sync routes --json
andurel sync payloads --json
andurel inspect routes --json
```

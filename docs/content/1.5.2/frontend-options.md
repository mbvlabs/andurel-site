# Frontend Options

Andurel supports server-rendered hypermedia and optional Inertia frontends from the same Go application.

## Templ and Datastar

The default stack renders type-safe HTML with Templ. Datastar adds signals, server-sent events, and HTML fragment updates without a client-side application framework.

Choose this approach when server rendering and small interactive updates fit the product.

```go
return hypermedia.RenderPage(etx, views.ProductsIndex(products))
```

## Inertia

Pass `--inertia` to create a Vue, React, or Svelte frontend powered by Vite:

```bash
andurel new orbit --inertia vue
andurel new orbit --inertia react/pnpm
andurel new orbit --inertia svelte/bun
```

Supported JavaScript runtimes are npm, pnpm, bun, and yarn. The selected adapter and runtime are stored in `andurel.lock`.

## Generated resources

Resource views default to Templ, including inside an Inertia-enabled project. Pass `--inertia` when generating a controller or scaffold to create pages for the configured adapter.

```bash
andurel generate scaffold Product
andurel generate scaffold Product --inertia
```

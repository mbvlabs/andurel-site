# Frontend Options

Andurel v2 supports server-rendered Templ with Datastar and an independently versioned Inertia v3 adapter.

## Templ and Datastar

The default stack renders type-safe HTML with Templ. Datastar adds signals, server-sent events, and fragment updates without requiring a client-side application framework.

```go
return hypermedia.RenderPage(etx, views.ProductsIndex(products))
```

Import the reusable renderer from `github.com/mbvlabs/andurel/pkg/hypermedia`.

## Inertia v3

Choose Vue, React, or Svelte at project creation:

```bash
andurel new orbit --inertia vue
andurel new orbit --inertia react/pnpm
andurel new orbit --inertia svelte/yarn
```

Controllers receive `*inertia.Renderer` and render serializable props. The server adapter, Vite integration, protocol behavior, and SSR runtime come from `github.com/mbvlabs/andurel/pkg/inertia`. Your application owns `views/root.templ`, frontend entry points, shared props, and `config/inertia.go`.

Resource generation still defaults to Templ, even inside an Inertia project. Opt in per resource:

```bash
andurel generate scaffold Product --inertia
andurel generate routes
```

`generate routes` creates typed TypeScript URL helpers from `router/routes/*.go`.

## Server-side rendering

SSR is disabled by default. `INERTIA_SSR_MODE=managed` starts and monitors a Node renderer through the application lifecycle. `external` connects to an operator-managed renderer at `INERTIA_SSR_URL`. Startup timeout, request timeout, response size, runtime, bundle, and fail-fast behavior are configured independently.

The JavaScript package manager in `andurel.lock` is separate from the SSR runtime; choosing Bun for installs does not silently replace Node for SSR.

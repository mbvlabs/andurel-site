# Views

Templ is the type-safe HTML view layer for hypermedia apps and for Inertia's root document. Datastar can update fragments and consume server-sent events; Inertia v3 can render Vue, React, or Svelte pages from the same Go backend.

## Templ pages

Define components under `views/` and render them with `pkg/hypermedia`:

```templ
templ ProductsIndex(products []models.Product) {
    @base(SetTitle("Products")) {
        <main><h1>Products</h1></main>
    }
}
```

```go
return hypermedia.RenderPage(etx, views.ProductsIndex(products))
```

Regenerate Go code after editing Templ files:

```bash
andurel sync views
```

This command also compiles Tailwind utilities in email templates.

## Datastar fragments

Named fragments let a controller return only changed HTML. Keep application state on the server and use signals or SSE where partial updates fit better than a full page navigation. See [Templ & Datastar](/docs/head/hypermedia).

## Inertia pages

An Inertia application owns `views/root.templ` and its adapter-specific files under `resources/js`. Pass application-facing structs or maps with stable JSON tags; do not serialize database model structs or narsilc-generated rows directly.

```go
return c.renderer.Page(etx, "Products/Index", inertia.FromStruct(props)).Render()
```

The v3 adapter supports partial reloads, deferred and once props, merge metadata, flash messages, redirects, asset-version reloads, and optional SSR. Configure the renderer once through Fx and inject it into controllers. See [Inertia](/docs/head/inertia).

# Frontend Options

Andurel supports server-owned HTML with Templ and Datastar, or client-owned pages with Inertia and Vue, React, or Svelte. Both use the same Echo router, models, services, sessions, validation, and PostgreSQL infrastructure. The choice is where rendering and interaction state live.

## Choose by ownership

| Concern | Templ and Datastar | Inertia |
| --- | --- | --- |
| Primary renderer | Go server | Vue, React, or Svelte browser runtime |
| Initial response | Complete Templ HTML | Templ root containing page JSON; optional SSR |
| Navigation | Browser requests and targeted Datastar actions | Inertia visits swap page components |
| Interaction state | Server state plus small signals | Frontend component and application state |
| Data contract | Typed Templ parameters | JSON payload structs and TypeScript types |
| Partial work | Named fragments and SSE events | Partial reloads, lazy/deferred/once props, merge metadata |
| JavaScript surface | Small by default | Full frontend ecosystem |
| Strong fit | Forms, content, operations screens, server-led workflows | App-like navigation and rich local interaction |

Complexity alone does not decide it. Templ can stream live multi-fragment updates, while a small Inertia page still introduces a client runtime and JSON boundary.

## Working with Templ

A controller passes typed values directly into a compiled component:

```go
func (c Products) Index(etx *echo.Context) error {
    products, err := c.products.List(etx.Request().Context())
    if err != nil { return err }
    return hypermedia.RenderPage(etx, views.ProductsIndex(products))
}
```

There is no serialization layer between controller and view. Parameter or type changes fail compilation. `RenderPage` writes full HTML, `RenderComponent` is its readability alias, and `RenderFragment(s)` extracts named fragments.

Datastar adds targeted interaction without creating a separate SPA. A controller reads signals, validates and applies domain work, then emits an SSE patch:

```go
var signals ArchiveSignals
if err := hypermedia.ReadSignals(etx.Request(), &signals); err != nil {
    return err
}
if err := c.products.Archive(etx.Request().Context(), signals.ID); err != nil {
    return err
}
return hypermedia.PatchComponent(etx, views.ProductRow(signals.ID))
```

Patches can replace, remove, append, prepend, insert, or update named fragments. Signals, custom events, URL changes, scripts, and long-lived broadcasters form the server-to-browser protocol. The application owns fragment identifiers and state transitions. See [Hypermedia](/docs/latest/hypermedia).

## Working with Inertia

An Inertia controller names a frontend page and crosses an explicit JSON boundary:

```go
return c.renderer.Page(etx, "Products/Index", inertia.Props{
    "products": toProductResources(products),
})
```

The first visit renders application-owned `views/root.templ`, embeds the page object, and loads Vite. Later visits return page JSON and the official adapter swaps the component. Go retains authorization, validation, queries, and payload construction; Vue, React, or Svelte owns page rendering and local interaction.

Define page payload structs with stable JSON tags and matching TypeScript declarations. Do not serialize database models or sqlc rows directly. Partial reloads and prop policies reduce work but do not remove the need for a deliberate browser API. See [Inertia](/docs/latest/inertia).

## Forms and validation

Templ commonly re-renders the form or a named fragment with typed submitted values and field errors, using a suitable status such as 422.

Inertia normally uses redirect-after-write. Pass mapped errors with `inertia.WithValidationErrors`, respect named error bags, preserve flash through redirects, and let the client adapter expose errors to the form.

Domain validation should be shared. Only its transport and presentation differ.

## Mixing both approaches

One Echo application can have Templ and Inertia routes. Keep each route's response contract consistent: an Inertia navigation must receive an Inertia page or explicit location response, not arbitrary HTML.

Resource generation defaults to Templ, including inside an Inertia scaffold. Opt in per resource:

```bash
andurel generate scaffold Product --inertia
andurel generate routes
```

`generate routes` creates typed TypeScript URL helpers from marked route declarations. Templ and frontend layouts are not shared components: one renders in Go, the other in the client runtime.

A clean mixed boundary is public Templ marketing pages beside an authenticated Inertia application. Avoid implementing the same feature twice without a concrete reason.

## Development, assets, and SSR

Both choices compile Templ because base documents and email use it:

```bash
andurel generate view
```

Inertia additionally runs Vite during development and embeds production output in the Go binary. The package manager recorded in `andurel.lock` controls dependency installation; it does not select the managed SSR runtime.

SSR is disabled by default. Managed mode starts a Node renderer with the web process; external mode connects to an operator-owned renderer. Each page still opts into SSR, and client rendering is the fallback unless fail-fast is enabled.

## Migration cost

Templ to Inertia replaces typed component calls with a JSON contract, recreates UI in a frontend framework, and adopts client navigation and form semantics. Inertia to Templ moves render state into Go and replaces local interactions with browser requests or Datastar signals and patches.

Models and services should survive either move. If they must change, frontend responsibilities have leaked into the domain boundary.

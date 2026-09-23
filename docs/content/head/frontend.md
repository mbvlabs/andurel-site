# Frontend

Andurel supports client-owned pages with Inertia and Vue, React, or Svelte as the default, or server-owned HTML with Templ and Datastar. Both use the same Echo router, models, services, sessions, validation, and PostgreSQL infrastructure. The choice is where rendering and interaction state live.

New projects default to **Inertia React/pnpm**. Pass `--ui templ/datastar` (or another Inertia combo) at `andurel new` time. The choice is stored in `andurel.toml` (`project.inertia`, `project.javascriptPackageManager`) and drives later generators.

## Choose by ownership

| Concern | Inertia (default) | Templ and Datastar |
| --- | --- | --- |
| Primary renderer | Vue, React, or Svelte browser runtime | Go server |
| Initial response | Templ root containing page JSON; optional SSR | Complete Templ HTML |
| Navigation | Inertia visits swap page components | Browser requests and targeted Datastar actions |
| Interaction state | Frontend component and application state | Server state plus small signals |
| Data contract | JSON payload structs and TypeScript types | Typed Templ parameters |
| Partial work | Partial reloads, lazy/deferred/once props, merge metadata | Named fragments and SSE events |
| JavaScript surface | Full frontend ecosystem | Small by default |
| Strong fit | App-like navigation and rich local interaction | Forms, content, operations screens, server-led workflows |

Complexity alone does not decide it. A small Inertia page still introduces a client runtime and JSON boundary, while Templ can stream live multi-fragment updates.

## Working with Inertia

An Inertia controller names a frontend page and crosses an explicit JSON boundary:

```go
return c.renderer.Page(etx, "Products/Index", inertia.Props{
    "products": toProductResources(products),
}).Render()
```

The first visit renders application-owned `views/root.templ`, embeds the page object, and loads Vite. Later visits return page JSON and the official adapter swaps the component. Go retains authorization, validation, queries, and payload construction; Vue, React, or Svelte owns page rendering and local interaction.

Define page payload structs with stable JSON tags and matching TypeScript declarations. Do not serialize database models or narsilc rows directly. Partial reloads and prop policies reduce work but do not remove the need for a deliberate browser API. See [Inertia](/docs/head/inertia) and [Props](/docs/head/inertia-props).

## Working with Templ

A controller passes typed values directly into a compiled component:

```go
func (c Products) Index(etx *echo.Context) error {
    products, err := c.products.List(etx.Request().Context())
    if err != nil {
        return err
    }
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

Patches can replace, remove, append, prepend, insert, or update named fragments. Signals, custom events, URL changes, scripts, and long-lived broadcasters form the server-to-browser protocol. The application owns fragment identifiers and state transitions. See [Templ & Datastar](/docs/head/hypermedia).

## Forms and validation

Inertia normally uses redirect-after-write. Pass mapped errors with `Page(...).ValidationErrors(...)`, respect named error bags, preserve flash through redirects, and let the client adapter expose errors to the form. See [Shared Data and Redirects](/docs/head/inertia-shared).

Templ commonly re-renders the form or a named fragment with typed submitted values and field errors, using a suitable status such as 422.

Domain validation should be shared. Only its transport and presentation differ.

## Mixing both approaches

One Echo application can have Templ and Inertia routes. Keep each route's response contract consistent: an Inertia navigation must receive an Inertia page or explicit location response, not arbitrary HTML.

Resource generation follows the UI recorded in `andurel.toml`. In an Inertia project, scaffolds and controllers emit Inertia pages for that adapter. Use `--api` for JSON handlers instead of pages:

```bash
andurel generate scaffold Product
andurel generate scaffold Product --api
andurel sync routes
```

`sync routes` creates typed TypeScript URL helpers from marked route declarations. Templ and frontend layouts are not shared components: one renders in Go, the other in the client runtime.

A clean mixed boundary is public Templ marketing pages beside an authenticated Inertia application. Avoid implementing the same feature twice without a concrete reason. To add Templ resources inside an Inertia app, change the project UI in `andurel.toml` or hand-author Templ handlers; generators follow the manifest.

## Development, assets, and SSR

Both choices compile Templ because base documents and email use it:

```bash
andurel sync views
```

Inertia additionally runs Vite during development and embeds production output in the Go binary. The package manager recorded in `andurel.toml` controls dependency installation; it does not select the Node runtime used by `cmd/ssr`.

SSR is per page: call `.SSR()` on the page builder. In development the renderer posts to Vite; otherwise `cmd/app` posts to `INERTIA_SSR_URL` and `cmd/ssr` owns Node. Client rendering is the fallback unless fail-fast is enabled. See [SSR](/docs/head/inertia-ssr).

## Migration cost

Inertia to Templ moves render state into Go and replaces local interactions with browser requests or Datastar signals and patches. Templ to Inertia replaces typed component calls with a JSON contract, recreates UI in a frontend framework, and adopts client navigation and form semantics.

Models and services should survive either move. If they must change, frontend responsibilities have leaked into the domain boundary.

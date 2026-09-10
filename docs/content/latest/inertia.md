# Inertia

`github.com/mbvlabs/andurel/pkg/inertia` is Andurel's Echo-native implementation of the Inertia v3 server protocol. It owns request classification, page responses, partial prop evaluation, redirects, Vite integration, and optional SSR. The browser adapter is the official `@inertiajs/*` package for Vue, React, or Svelte.

The package has no Fx dependency. Generated applications construct one `Renderer` in `cmd/app`, register `renderer.Middleware()`, and inject `*inertia.Renderer` into controllers. Node process ownership belongs to `cmd/ssr`, not the web process.

## Start here

Read this page for the contract. Use the child pages when you need to implement a controller, wire SSR, or run a generator.

| Topic | Use it when |
| --- | --- |
| [Renderer](/docs/latest/inertia-renderer) | Constructing the renderer, middleware, shared providers, and flash reflash |
| [Pages and Visits](/docs/latest/inertia-pages) | Calling `Page`, understanding initial HTML vs JSON visits, and reading protocol headers |
| [Root Document and Vite](/docs/latest/inertia-vite) | Editing `views/root.templ`, development tags, production manifests, and asset versions |
| [Props](/docs/latest/inertia-props) | Building JSON payloads, lazy resolvers, and v3 evaluation policies |
| [Shared Data and Redirects](/docs/latest/inertia-shared) | Sharing actor data, validation errors, flash, and redirect-after-write |
| [SSR](/docs/latest/inertia-ssr) | Opting a page into SSR and running the separate Node runtime |
| [Diagnostics](/docs/latest/inertia-diagnostics) | Classifying protocol, prop, root, and SSR failures |
| [Generators](/docs/latest/inertia-generators) | Scaffolding an Inertia app and generating pages, types, and TypeScript routes |

## Application shape

Controllers name a frontend page and cross an explicit JSON boundary:

```go
return c.renderer.Page(etx, "Products/Index", inertia.Props{
    "products": productResources,
}).Render()
```

A normal browser request receives the full `views.Root` document, page JSON, and development or production Vite tags. A request with `X-Inertia: true` receives only the JSON page object. Both append `X-Inertia` to `Vary`.

Go retains authorization, validation, queries, and payload construction. Vue, React, or Svelte owns page rendering and local interaction. Do not serialize Bun models or sqlc rows directly; define backend-owned payload structs with stable JSON tags.

## Protocol essentials

The page always contains `component`, `props`, `url`, and `version`. Conditional metadata covers merging, deferred props, rescued failures, shared keys, once retention, scroll, history, and flash. `props.errors` is protected and always an object. Flash is a top-level page field, not an ordinary prop.

After writes, `Renderer.Redirect` issues a 3xx response. Middleware upgrades unsafe 302 responses to 303, converts empty successful responses into a redirect back, and preserves flash across redirect chains. `Renderer.Location` forces a full browser visit. Version mismatches reload GET requests without replaying unsafe methods.

SSR is per-response. Call `.SSR()` on a page builder when that document should be server-rendered. In development the renderer posts to Vite; in other environments `cmd/app` posts to `INERTIA_SSR_URL`. Failures fall back to client rendering unless fail-fast is enabled.

For unexpected behavior, inspect `X-Inertia-*` headers, the exact component name, selected nested paths, `Vary`, asset version, and redirect status. See [Diagnostics](/docs/latest/inertia-diagnostics).

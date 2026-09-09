# Inertia

`github.com/mbvlabs/andurel/pkg/inertia` is Andurel's Echo-native implementation of the Inertia v3 server protocol. It owns request classification, page responses, partial prop evaluation, redirects, Vite integration, and optional SSR. The browser adapter is the official `@inertiajs/*` package for Vue, React, or Svelte.

## Renderer construction and lifecycle

Construct one `Renderer` and inject it into controllers. Generated projects supply embedded assets, an application-owned Templ root, Vite settings, identity, shared props, and SSR policy.

```go
renderer, err := inertia.New(
    inertia.WithRoot(views.Root),
    inertia.WithAssetFS(assets.Files),
    inertia.WithProjectName(appCfg.ProjectName),
    inertia.WithEnvironment(appCfg.Environment),
    inertia.WithBuildPathURL(routes.ViteBuild.Path()),
    inertia.WithEntryPoint(cfg.EntryPoint),
    inertia.WithContainerID(cfg.ContainerID),
    inertia.WithViteDevURL(cfg.ViteDevURL),
    inertia.WithShared(inertia.Props{"appVersion": appVersion}),
    inertia.WithProtocolDebug(cfg.ProtocolDebug),
)
```

The package has no Fx dependency. The application registers `renderer.Start` and `renderer.Shutdown` hooks; they are no-ops except for managed SSR.

## Initial and client visits

Controllers use one call for both response shapes:

```go
return c.renderer.Page(etx, "Products/Index", inertia.Props{
    "products": productResources,
})
```

A normal browser request receives the full `views.Root` document, page JSON, and development or production Vite tags. A request with `X-Inertia: true` receives only the JSON page object and `X-Inertia: true`. Both append `X-Inertia` to `Vary`.

The page always contains `component`, `props`, `url`, and `version`. Conditional metadata covers merging, deferred, rescued, shared, once, scroll, history, and flash behavior. `props.errors` is protected and always an object; flash is a top-level field.

## Root document and Vite

`views/root.templ` controls the initial document shell, metadata, mount element, SSR insertion, and asset tags. Page layout still belongs in Vue, React, or Svelte components.

```templ
templ Root(data inertia.RootData) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <title>{ data.ProjectName }</title>
            @inertia.SSRHead(data.SSR)
            @templ.Raw(string(data.ViteHead))
        </head>
        <body>
            if data.SSR != nil {
                @inertia.SSRBody(data.SSR)
            } else {
                @inertia.PageScript(data.ContainerID, data.PageJSON)
                @inertia.AppMount(data.ContainerID)
            }
            @templ.Raw(string(data.ViteBody))
        </body>
    </html>
}
```

Development assets come from `INERTIA_VITE_DEV_URL`. Production assets resolve through the embedded Vite manifest. `WithBuildPathURL` also supplies the default asset version; override it with `WithVersion` or a request-scoped `WithVersionProvider`.

## Props and evaluation policies

Use `Props` or `FromStruct` with backend-owned payload structs. JSON tags define browser names. Do not pass Bun models or sqlc rows directly.

Plain values resolve immediately. `Resolver` and `PropResolver` defer computation. Policies compose around values or resolvers:

- `Always` survives partial filtering.
- `Optional` resolves only when explicitly selected.
- `Deferred` announces a reload group; `InGroup` names it and `Rescue` makes failures recoverable.
- `Merge`, `Prepend`, `DeepMerge`, and `MatchOn` describe client merging.
- `Once` supports a stable key, absolute or relative expiry, and forced refresh.
- `Scroll` and `ScrollAt` add infinite-scroll metadata and merge behavior.

```go
props := inertia.Props{
    "account": inertia.Always(account),
    "filters": inertia.Optional(loadFilters),
    "stats": inertia.Deferred(loadStats, inertia.InGroup("dashboard"), inertia.Rescue()),
    "products": inertia.Scroll(page, page),
    "countries": inertia.Once(loadCountries, inertia.OnceFor(24*time.Hour)),
}
```

Nested dot paths participate in `only` and `except` partial reloads. `only` is applied first. Filtering is ignored when the request's partial component differs from the rendered component. Resolvers run only when their policy and the active request require them.

## Shared props, errors, flash, and redirects

`WithShared` copies static props; `WithSharedProvider` adds request-scoped values such as the current actor. Page props replace shared props except for protected `errors`. `WithValidationErrors` populates errors and respects named error bags.

Use `WithFlash` or `WithFlashProvider` for the v3 page flash field. Generated session integration installs reflash handling so consumed messages survive redirect chains.

After writes, `Renderer.Redirect` is normalized from 302 to 303 for POST, PUT, PATCH, and DELETE Inertia requests. `Renderer.Location` produces a 409 external visit. Empty successful responses redirect back. The renderer middleware also performs version reloads for GET mismatches without replaying unsafe methods.

## SSR modes

SSR configuration selects a backend; `inertia.WithSSR()` opts an individual page response into it.

| Mode | Ownership | Behavior |
| --- | --- | --- |
| `disabled` | None | Client-mount document; the default |
| `managed` | Web process | Starts, health-checks, monitors, and stops a Node renderer |
| `external` | Operator | Connects to an already managed HTTP renderer |

Managed mode defaults to Node.js 22 or newer and `assets/dist/ssr/ssr.js`. Startup timeout, render timeout, URL, and maximum response size are independent. SSR failures fall back to client rendering unless fail-fast is enabled, which is useful in release verification.

## Diagnostics and errors

`INERTIA_PROTOCOL_DEBUG=true` logs request classification and metadata keys without logging prop or session values. Renderer errors are typed with a kind, operation, component, method, URL, and wrapped cause, allowing prop, protocol, root-render, and SSR failures to be distinguished.

For unexpected behavior, inspect the incoming `X-Inertia-*` headers, exact component name, selected nested paths, response `Vary`, asset version, and redirect status. For SSR, separately verify mode, runtime version, bundle path, health URL, timeout, and response-size bound.

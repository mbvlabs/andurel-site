# Renderer

Construct one `inertia.Renderer` at process start and inject it into controllers. The package is an ordinary Go constructor: it does not register Fx hooks, read environment variables, or own a Node process.

Generated applications do this in `cmd/app` through `newInertia`, then pass the renderer into router middleware and every Inertia controller.

## Construction

`NewRenderer` takes required protocol settings as positional arguments. Options configure application-owned behavior.

```go
renderer, err := inertia.NewRenderer(
    cfg.ContainerID,
    routes.ViteBuild.Path(),
    cfg.EntryPoint,
    cfg.ViteDevURL,
    cfg.SSRURL,
    cfg.SSRRequestTimeout,
    cfg.SSRMaxResponseBytes,
    inertia.WithRoot(views.Root),
    inertia.WithAssetFS(assets.Files),
    inertia.WithProjectName(appCfg.ProjectName),
    inertia.WithEnvironment(appCfg.Environment),
    inertia.WithProtocolDebug(cfg.ProtocolDebug),
    inertia.WithShared(inertia.Props{"appVersion": appVersion}),
    inertia.WithSSRFailFast(cfg.SSRFailFast),
)
```

| Argument | Typical source | Role |
| --- | --- | --- |
| `containerID` | `INERTIA_CONTAINER_ID`, default `app` | DOM id for `AppMount` and the `data-page` script |
| `buildPathURL` | `routes.ViteBuild.Path()` | Production asset URL prefix and default asset version |
| `entryPoint` | `INERTIA_ENTRY_POINT` | Vite manifest key and development module URL |
| `viteDevURL` | `INERTIA_VITE_DEV_URL` | Development asset origin; also derives `/__inertia_ssr` |
| `ssrURL` | `INERTIA_SSR_URL` | Where `cmd/app` POSTs `/render` outside development |
| `ssrTimeout` | `INERTIA_SSR_REQUEST_TIMEOUT` | Render, health, and shutdown deadline |
| `ssrMaxResponseBytes` | `INERTIA_SSR_MAX_RESPONSE_BYTES` | SSR body size bound |

`WithRoot` is required. Production also requires `WithAssetFS` so the renderer can read `dist/vite/manifest.json`. Empty `containerID`, a nil root, or a missing production manifest fail construction before the process serves traffic.

`WithEnvironment` selects development Vite tags versus the production manifest. The string `"production"` is the production check; any other environment uses the Vite dev server.

## Options

| Option | Purpose |
| --- | --- |
| `WithRoot` | Application-owned Templ document, usually `views.Root` |
| `WithAssetFS` | Embedded filesystem used for the production Vite manifest |
| `WithProjectName` | Passed to the root as `RootData.ProjectName` |
| `WithEnvironment` | Selects development or production asset tags |
| `WithVersion` | Fixed asset version; defaults to `buildPathURL` |
| `WithVersionProvider` | Request-scoped version that overrides the fixed value |
| `WithShared` | Copied static shared props |
| `WithSharedProvider` | Request-scoped shared props such as the current actor |
| `WithFlashProvider` | Extra flash source; the renderer already reads `FlashFromContext` |
| `WithSSRRenderer` | Replace the default HTTP or Vite SSR client |
| `WithSSRFailFast` | Return SSR errors instead of falling back to client rendering |
| `WithReflash` | Preserve flash across redirects |
| `WithProtocolDebug` | Log request classification and metadata keys |

Shared props are copied at construction. Later mutation of the caller's map does not change the renderer. Providers run on every page response; they must be deterministic and cheap enough for full and partial visits.

`SetReflashHandler` is the runtime equivalent of `WithReflash`. Generated routers call it while assembling middleware so session code can stay in the application:

```go
if err := renderer.SetReflashHandler(func(etx *echo.Context) error {
    flashes := appctx.Flashes(etx.Request().Context())
    return cookieSession.Reflash(etx, flashes)
}); err != nil {
    return nil, err
}
```

## Middleware

`Renderer.Middleware()` is the protocol gate. Register it after session middleware and before handlers that call `Page`, `Redirect`, or `Location`.

Generated order is:

```text
session store
session validation
renderer.Middleware()
request metadata (flashes onto context)
CORS
CSRF
```

The middleware parses every Inertia header into a `Request` and stores it on the Echo request context. Non-Inertia requests continue unchanged except that later page responses still append `X-Inertia` to `Vary`.

For Inertia GET requests it compares `X-Inertia-Version` with the current asset version. A mismatch returns `409` with `X-Inertia-Location` set to the current URL and `X-Inertia` removed. Unsafe methods are not replayed as a version reload.

For Inertia visits it captures the handler response and then:

1. Appends `X-Inertia` to `Vary`.
2. Turns an empty `200` into a redirect to `Referer`, or `/` when Referer is missing.
3. Upgrades `302` after POST, PUT, PATCH, or DELETE to `303`.
4. Converts fragment redirects (`Location` containing `#`) into a `409` control response with `X-Inertia-Redirect`, except for prefetch requests.
5. Runs the reflash handler when the handler issued a 3xx response.

Because the writer is captured, `Page` must finish before the response is committed. Do not write to the Echo response yourself in an Inertia controller; return `Page(...).Render()`, `Redirect`, `Location`, or `FreshRedirect`.

## Lifecycle

`Renderer` has no `Start` or `Shutdown`. Vite is an external development process. Node SSR is `NewSSRRuntime` in `cmd/ssr`. The web process only constructs an HTTP or Vite SSR client.

Fx belongs to the application:

```go
var inertiaModule = fx.Module(
    "inertia",
    fx.Provide(newInertia),
)
```

Tests can construct a renderer with a stub root and in-memory asset filesystem without starting Echo or Node.

## Application usage

Inject `*inertia.Renderer`, not a wider application container:

```go
type Products struct {
    products models.Products
    renderer *inertia.Renderer
}

func NewProducts(products models.Products, renderer *inertia.Renderer) Products {
    return Products{products: products, renderer: renderer}
}
```

Keep renderer construction in `cmd/app`. Controllers should not build a second renderer per request. If a test needs different shared props, construct a dedicated renderer or pass page props that replace the shared keys.

See [Pages and Visits](/docs/latest/inertia-pages) for `Page` and [Shared Data and Redirects](/docs/latest/inertia-shared) for providers, flash, and redirects.

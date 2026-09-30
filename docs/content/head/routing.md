# Routing

Andurel keeps route declarations separate from controller registration and provides typed URL helpers for application links via `github.com/mbvlabs/andurel/pkg/routing`.

## Declare a route

Generated files in `router/routes` declare route values. The router binds them to controllers and middleware.

```go
import "github.com/mbvlabs/andurel/pkg/routing"

var ProductShow = routing.NewRouteWithUUIDID(
    "/products/:id",
    "products.show",
    "",
    routing.InertiaRoute(),
)
```

`Path()` returns the Echo registration path, `Name()` returns the name, `URL(...)` substitutes typed parameters, and `FullURL(...)` prepends the host origin. `Host()` returns the `routing.HostName` (primary unless you pass `routing.Host(...)`). Use `InertiaRoute()` when the route should appear in generated TypeScript helpers.

Use the helper instead of hard-coding an application URL:

```go
routes.ProductShow.URL(product.ID)
```

## Typed parameter routes

| Constructor | Placeholder | Go argument |
| --- | --- | --- |
| `NewSimpleRoute` | none | none |
| `NewRouteWithUUIDID` | `:id` | `uuid.UUID` |
| `NewRouteWithSerialID` | `:id` | `int32` |
| `NewRouteWithBigSerialID` | `:id` | `int64` |
| `NewRouteWithStringID` | `:id` | `string` |
| `NewRouteWithSlug` | `:slug` | `string` |
| `NewRouteWithToken` | `:token` | `string` |
| `NewRouteWithParams[T]` | tagged fields | typed struct |

Multiple parameters use a struct so call sites cannot swap same-typed values:

```go
type ProjectParams struct {
    Team    string `param:"team"`
    Project string `param:"project"`
}
var ProjectShow = routing.NewRouteWithParams[ProjectParams](
    "/:team/projects/:project", "show", "projects",
)
url := ProjectShow.URL(ProjectParams{Team: "platform", Project: "orbit"})
```

## Register a handler

Controllers register Echo routes with the application router. Pass the hosted route first; `AddRoute` fills empty `Path` and `Name` from it so controllers do not copy them:

```go
func (p Products) RegisterRoutes(r *router.Router) error {
    _, err := r.AddRoute(routes.ProductShow, echo.Route{
        Method:  http.MethodGet,
        Handler: p.Show,
    })
    return err
}
```

Fx stays `RegisterRoutes(r)` with no host argument. The host is declared only on the route. A missing `HOST_*` for that `HostName` fails at boot, naming the route.

The router and controllers are Fx modules. Dependencies are constructor parameters; do not store them in an Echo or standard-library context.

## Multiple hosts

Single-host apps are the default: set `HOST_PRIMARY` and keep path namespaces such as `/admin/...` on the primary Echo.

To serve a second domain from the same process:

1. Declare a host name in `config/hosts.go`, for example `const HostAdmin routing.HostName = "admin"`.
2. Load `HOST_ADMIN` (or another `HOST_*`) into `App.Hosts` at boot. Missing hostnames fail startup.
3. Register aliases (`www`, extra names) as additional hostnames for that `HostName`.
4. Declare the host on the route with `routing.Host(config.HostAdmin)`.

```go
var WidgetIndex = routing.NewSimpleRoute(
    "/widgets",
    "admin.widgets.index",
    "",
    routing.InertiaRoute(),
    routing.Host(config.HostAdmin),
)
```

`--host` and `--prefix` on `generate controller` / `generate scaffold` are independent:

```bash
andurel generate controller Widget --host=admin --prefix=admin   # HostAdmin, /admin/widgets
andurel generate controller Widget --host=admin                  # HostAdmin, /widgets
andurel generate controller Widget --prefix=admin                # HostPrimary, /admin/widgets
```

Do not invent `--host` values: declare `routing.HostName` consts in `config/hosts.go` and load them into `App.Hosts` first. `primary` always works.

The reverse proxy must forward the public `Host` header unchanged. When a configured hostname has no port, Andurel also registers `hostname:80` and `hostname:443` aliases. Unknown hosts are served by the primary Echo. Duplicate hostnames or aliases across `HostName`s fail at boot.

Layout CSS and JS stay on the primary host unless those routes set `routing.Host(...)`. Layout `<link>` / `<script>` tags should use `FullURL()` so they always point at the asset host.

`/api` stays on the primary host unless those routes set `routing.Host(...)`.

Inertia Vite files stay on the primary host. Production tags from `pkg/inertia` use `ViteBuild.FullURL()` plus `crossorigin` so a secondary-host page loads the entry from primary. Vite production `base` is `'./'` so hashed chunks and fonts resolve against that script URL, not the page origin. Development keeps `base: '/assets/dist/'` and tags that hit `INERTIA_VITE_DEV_URL`.

Inertia TypeScript helpers return full URLs after boot calls `configureRouteHosts` from the shared `hosts` prop. Use `routes.name.path()` for a relative path. Not-found handlers are registered on every configured host (`AddRouteNotFoundEachHost`).

See [Configuration](/docs/head/configuration) and [TypeScript Sync](/docs/head/inertia-typescript-sync).

## Inspect and export routes

```bash
andurel inspect routes --json
andurel sync routes
```

The manifest reports route names, paths, parameters, and source locations. In Inertia projects, `sync routes` writes `resources/js/routes.ts`. See [TypeScript Sync](/docs/head/inertia-typescript-sync).

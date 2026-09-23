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

`Path()` returns the Echo registration path, `Name()` returns the name, `URL(...)` substitutes typed parameters, and `FullURL(base, ...)` prepends the application base URL. Use `InertiaRoute()` when the route should appear in generated TypeScript helpers.

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

Controllers register Echo routes with the application router:

```go
_, err := r.AddRoute(echo.Route{
    Method:  http.MethodGet,
    Path:    routes.ProductShow.Path(),
    Name:    routes.ProductShow.Name(),
    Handler: products.Show,
})
```

The router and controllers are Fx modules. Dependencies are constructor parameters; do not store them in an Echo or standard-library context.

## Inspect and export routes

```bash
andurel inspect routes --json
andurel sync routes
```

The manifest reports route names, paths, parameters, and source locations. In Inertia projects, `sync routes` writes `resources/js/routes.ts`. See [TypeScript Sync](/docs/head/inertia-typescript-sync).

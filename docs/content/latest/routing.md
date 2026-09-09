# Routing

Andurel keeps route declarations separate from controller registration and provides typed URL helpers for application links.

## Declare a route

Routes use the independently versioned routing package:

```go
import "github.com/mbvlabs/andurel/pkg/routing"

var ProductShow = routing.NewRouteWithUUIDID(
    "/products/:id",
    "products.show",
    "",
)
```

Use the helper instead of hard-coding an application URL:

```go
routes.ProductShow.URL(product.ID)
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
andurel routes --json
andurel generate routes
```

The manifest reports route names, paths, parameters, and source locations. In Inertia projects, the generator uses that same manifest to write `resources/js/routes.ts`.

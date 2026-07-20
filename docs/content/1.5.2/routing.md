# Routing

Andurel keeps route declarations separate from controller registration and provides typed URL helpers for application links.

## Declare a route

Routes live under `router/routes`:

```go
var ProductShow = routing.NewRouteWithUUIDID(
    "/products/:id",
    "products.show",
    "",
)
```

Use the generated helper instead of hard-coding application URLs:

```go
routes.ProductShow.URL(product.ID)
```

## Register a route

Controllers register Echo routes with the shared router:

```go
_, err := r.AddRoute(echo.Route{
    Method:  http.MethodGet,
    Path:    routes.ProductShow.Path(),
    Name:    routes.ProductShow.Name(),
    Handler: p.Show,
})
```

## Inspect routes

Use the CLI to inspect route names, paths, parameters, and source files:

```bash
andurel routes
andurel routes --json
```

Inertia projects can generate TypeScript URL helpers with `andurel generate routes`.

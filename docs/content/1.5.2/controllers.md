# Controllers

Controllers receive Echo requests, coordinate application work, and return HTML, Inertia, JSON, redirects, or errors.

## Register handlers

A controller exposes `RegisterRoutes` and is connected through the controllers Fx module. Keep request parsing and response decisions at this boundary.

## Render HTML

Use the hypermedia helpers for Templ pages and fragments:

```go
func (p Products) Index(etx *echo.Context) error {
    products, err := p.models.All(etx.Request().Context())
    if err != nil {
        return err
    }
    return hypermedia.RenderPage(etx, views.ProductsIndex(products))
}
```

## Generate controllers

Generate a complete resource controller or selected actions:

```bash
andurel generate controller Product
andurel generate controller Product index show
andurel generate controller Dashboard overview
```

Use `--api` for JSON controllers and `--inertia` for pages using the configured Inertia adapter.

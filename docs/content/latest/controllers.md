# Controllers

Controllers are the HTTP boundary. They receive Echo contexts, parse and validate input, coordinate injected models or services, and choose an HTML, Inertia, JSON, redirect, or error response.

## Explicit dependencies

Model APIs are values such as `models.Products`, not package globals. Fx supplies them to the controller constructor:

```go
type Products struct {
    products models.Products
}

func NewProducts(products models.Products) Products {
    return Products{products: products}
}
```

Use `etx.Request().Context()` when calling models and services so cancellation crosses the HTTP boundary.

## Render a Templ page

```go
func (p Products) Index(etx *echo.Context) error {
    products, err := p.products.All(etx.Request().Context())
    if err != nil {
        return err
    }
    return hypermedia.RenderPage(etx, views.ProductsIndex(products))
}
```

For Inertia, inject `*inertia.Renderer` and call its page, redirect, or location helpers. See [Inertia](/docs/latest/inertia). For JSON controllers, return Echo JSON responses and keep application types at the boundary.

## Generate controllers

```bash
andurel generate controller Product
andurel generate controller Product index show
andurel generate controller Dashboard overview
andurel generate controller v1/User --api
```

Use `--model-name` when a controller is backed by a differently named model, `--api` for JSON handlers, or `--inertia` for pages using the configured frontend adapter.

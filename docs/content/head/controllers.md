# Controllers

Controllers are the HTTP boundary. They receive Echo contexts, parse and validate input, coordinate injected models or services, and choose an Inertia, Templ, JSON, redirect, or error response.

## Explicit dependencies

Model APIs are values such as `models.Products`, not package globals. Fx supplies them (and the Inertia renderer when needed) to the controller constructor:

```go
type Products struct {
    products models.Products
    renderer *inertia.Renderer
}

func NewProducts(products models.Products, renderer *inertia.Renderer) Products {
    return Products{products: products, renderer: renderer}
}
```

Use `etx.Request().Context()` when calling models and services so cancellation crosses the HTTP boundary. Do not stash collaborators on Echo context keys.

## Inertia (default UI)

Generated Inertia apps inject `*inertia.Renderer`. One `Page` call covers the first HTML document and later JSON visits:

```go
func (p Products) Index(etx *echo.Context) error {
    products, err := p.products.All(etx.Request().Context())
    if err != nil {
        return err
    }

    return p.renderer.Page(
        etx,
        "Products/Index",
        inertia.FromStruct(ProductIndexProps{Items: toProductData(products)}),
    ).Render()
}
```

Validation failures stay on the same visit with protected errors:

```go
if validationErrors, ok := validation.As(err); ok {
    return p.renderer.Page(etx, "Products/Edit", props).
        ValidationErrors(validationErrors.ToMap()).
        Render()
}
```

Redirect helpers:

| Method | Use |
| --- | --- |
| `Redirect(etx, location, status)` | Ordinary Inertia-aware redirect (POST upgrades 302 to 303) |
| `Location(etx, location)` | Hard visit / full document load (login success often uses this) |

```go
return p.renderer.Redirect(etx, routes.ProductIndex.URL(), http.StatusSeeOther)
```

Map domain rows to payload structs before `Page`. Do not pass `models.Product` or narsilc types into props. See [Props](/docs/head/inertia-props) and [Shared Data and Redirects](/docs/head/inertia-shared).

## Templ and JSON

Templ controllers skip the renderer and call hypermedia helpers:

```go
func (p Products) Index(etx *echo.Context) error {
    products, err := p.products.All(etx.Request().Context())
    if err != nil {
        return err
    }
    return hypermedia.RenderPage(etx, views.ProductsIndex(products))
}
```

API / `--api` controllers return Echo JSON and keep application DTOs at the boundary:

```go
return etx.JSON(http.StatusOK, ProductIndexProps{Items: items})
```

| Response | Typical call | Contract |
| --- | --- | --- |
| Inertia | `renderer.Page(...).Render()` | Page name + JSON props / redirects |
| Templ | `hypermedia.RenderPage` / fragments | Typed Templ components |
| JSON | `etx.JSON` | Application DTOs only |

## Generate controllers

```bash
andurel generate controller Product
andurel generate controller Product index show
andurel generate controller Dashboard overview
andurel generate controller v1/User --api
```

Use `--model-name` when a controller is backed by a differently named model, or `--api` for JSON handlers. Otherwise generation follows the UI recorded in `andurel.toml` (Inertia pages by default, or Templ).

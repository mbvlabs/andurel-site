# Pages and Visits

Every Inertia screen is one `Renderer.Page` call. The renderer decides whether to return the Templ root document or a JSON page object from the incoming headers. Controllers do not branch on `X-Inertia` themselves.

## Page builder

```go
func (c Products) Index(etx *echo.Context) error {
    products, err := c.products.List(etx.Request().Context())
    if err != nil {
        return err
    }

    return c.renderer.Page(etx, "Products/Index", inertia.Props{
        "products": toProductResources(products),
    }).Render()
}
```

`Page` returns a `*PageBuilder`. Chain options, then call `Render`:

| Method | Effect |
| --- | --- |
| `SSR()` | Opt this initial document into server-side rendering |
| `Status(code)` | HTTP status, default `200` |
| `ValidationErrors(map[string]string)` | Protected `props.errors`, with named bags when requested |
| `HistoryEncryption(bool)` | Sets `encryptHistory` metadata; the client encrypts history |
| `HistoryClear()` | Sets `clearHistory` |
| `PreserveFragment()` | Sets `preserveFragment` across a redirect |
| `Flash(value)` | Sets the top-level `flash` field for this response |

`Render` resolves shared props, evaluates page props against the current request, encodes the v3 page object, and writes either JSON or HTML.

Component names are adapter paths, not Go types. Generated React pages live at `resources/js/Pages/Products/Index.tsx` and resolve from:

```ts
import.meta.glob('./Pages/**/*.tsx', { eager: true })
```

The string `"Products/Index"` must match that file. A missing module fails in the browser, not at Go compile time. Treat the component name as part of the HTTP contract.

## Initial and client visits

`ParseRequest` recognizes an Inertia visit only when `X-Inertia` is `true`.

An initial visit has no Inertia header. `Render` writes `views.Root` with:

- validated page JSON in `PageScript`
- an empty mount node from `AppMount`
- Vite tags in `RootData.ViteHead` and `ViteBody`
- SSR markup instead of the script and mount when `.SSR()` succeeds

A client visit receives `application/json` with `X-Inertia: true`. The official adapter swaps the page component. Both responses append `X-Inertia` to `Vary` without replacing existing values.

The JSON page always includes:

```json
{
  "component": "Products/Index",
  "props": { "errors": {} },
  "url": "/products",
  "version": "/assets/dist/vite/*"
}
```

`url` is the request URI, including query string. Empty conditional fields are omitted: merge lists, deferred groups, once metadata, scroll metadata, history flags, and flash.

Error pages are ordinary pages with a non-200 status:

```go
return c.renderer.Page(etx, "Errors/NotFound", inertia.Props{}).
    Status(http.StatusNotFound).
    Render()
```

## Request headers

The middleware parses this state once per request:

| Header | Field | Meaning |
| --- | --- | --- |
| `X-Inertia` | `Inertia` | Client visit when `true` |
| `X-Inertia-Version` | `Version` | Asset version for GET mismatch reloads |
| `X-Inertia-Partial-Component` | `PartialComponent` | Component the `only`/`except` lists apply to |
| `X-Inertia-Partial-Data` | `Only` | Comma-separated paths to include |
| `X-Inertia-Partial-Except` | `Except` | Comma-separated paths to omit |
| `X-Inertia-Reset` | `Reset` | Merge paths to replace instead of merge |
| `X-Inertia-Error-Bag` | `ErrorBag` | Named bag for `props.errors` |
| `X-Inertia-Infinite-Scroll-Merge-Intent` | `MergeIntent` | `append` or `prepend` |
| `X-Inertia-Except-Once-Props` | `ExceptOnceProps` | Once keys the client still holds |
| `Purpose` | `Purpose` | `prefetch` disables fragment-redirect rewriting |

Lists are trimmed, empty items discarded, and duplicates removed. An invalid merge intent is a protocol error.

Partial filtering applies only when `PartialComponent` equals the rendered component. A visit to `Products/Show` with partial data for `Products/Index` is a full evaluation. Nested paths match in both directions: selecting `product` keeps `product.name`; selecting `product.name` keeps the parent object structure. `only` is applied before `except`.

## Prefetch and history

`Request.IsPrefetch()` is true when `Purpose` is `prefetch`. Prefetch visits still receive page JSON, but fragment redirects stay ordinary 3xx responses so the client can follow them without a control-response dance.

`HistoryEncryption`, `HistoryClear`, and `PreserveFragment` are metadata. The Go adapter does not encrypt browser history. Use them when the official client should encrypt, clear, or keep a URL fragment.

## Application usage

Use one `Page` call for GET screens, validation redisplay, and error pages. After a successful write, prefer `Redirect` or `Location` instead of rendering the next page from the mutating request. See [Shared Data and Redirects](/docs/latest/inertia-shared).

Keep payload construction next to the controller or in a small resource mapper. Generated scaffolds use `inertia.FromStruct` with `ProductIndexProps` / `ProductItemProps`. Hand-written pages can use `inertia.Props` or `FromStruct`; both JSON-encode the same way.

```go
return c.renderer.Page(
    etx,
    "Products/Show",
    inertia.FromStruct(ProductItemProps{Item: newProductData(product)}),
).Render()
```

The first visit and later visits stay the same Go code. If a page needs SSR, add `.SSR()` to that builder only. See [SSR](/docs/latest/inertia-ssr).

# Props

Props are the JSON contract between a Go controller and a Vue, React, or Svelte page. `pkg/inertia` evaluates them per request: some values resolve immediately, some wait for a partial reload, and some only announce metadata.

Treat this as a browser API. JSON tags are the public field names. Changing a tag, omitting a field, or serializing a database row is a breaking change for the page.

## Maps and structs

`Props` is `map[string]any` with a fluent `Set`:

```go
props := inertia.Props{
    "title": "Catalog",
}.Set("productCount", len(products))
```

`FromStruct` reflects exported fields and uses JSON tags as browser names. `json:"-"` omits a field. Embedded structs flatten when the tag is empty. Invalid input yields empty props; `FromStructChecked` returns the error instead.

```go
type ProductData struct {
    ID    uuid.UUID `json:"id"`
    Name  string    `json:"name"`
    Price int       `json:"price"`
}

type ProductIndexProps struct {
    Items []ProductData `json:"items"`
}

return c.renderer.Page(
    etx,
    "Products/Index",
    inertia.FromStruct(ProductIndexProps{Items: items}),
).Render()
```

Generated scaffolds emit these payload types next to the controller and matching TypeScript declarations under `resources/js/types`. Do not pass `models.Product` or sqlc rows into `Page`. Map to an application-owned struct first so nullability, money, and hidden columns stay explicit.

Dotted keys unpack into nested objects:

```go
inertia.Props{"product.name": "Orbit", "product.price": 12}
```

becomes `{"product": {"name": "Orbit", "price": 12}}`. A dotted path that collides with a non-object value is a props error.

Nested maps and structs are walked during evaluation so partial paths such as `product.name` work. Values that implement `MarshalJSON`, including `time.Time`, stay opaque and are not walked.

## Lazy resolvers

Plain values encode immediately. These resolver forms run only when the prop is included:

- `inertia.Resolver` (`func(*echo.Context) (any, error)`)
- `inertia.PropResolver` (`ResolveInertiaProp(*echo.Context) (any, error)`)
- `func(*echo.Context) (any, error)`
- `func(*echo.Context) any`
- `func() (any, error)`
- `func() any`

Use them for work you do not want on every full visit: extra queries, permission-gated panels, or expensive aggregations.

```go
"stats": inertia.Deferred(func() (any, error) {
    return c.products.DashboardStats(ctx)
}),
```

A resolver may return another `Prop` so policies can be decided after loading data. Resolver errors fail the response unless the prop is deferred and marked `Rescue`.

## Inclusion policies

Policies wrap a value or resolver. They compose: `Always(Optional(...))` is valid because constructors start from the inner `Prop`.

### Always

`Always` survives `only` / `except` filtering. Use it for data every partial reload still needs, such as the current account on a layout-heavy page.

### Optional

`Optional` is omitted on full visits and resolved only when a matching partial reload selects it. Use it for tabs or drawers that the client asks for explicitly.

### Deferred

`Deferred` is announced on full visits and resolved on a later partial reload. The client sees the path in `deferredProps` grouped by name. `InGroup("dashboard")` names the group; the default group is `"default"`. `Rescue()` converts a resolver failure into `rescuedProps` and omits the value instead of failing the page.

```go
"stats": inertia.Deferred(
    loadStats,
    inertia.InGroup("dashboard"),
    inertia.Rescue(),
),
```

Deferred and optional props do not run on a full visit. Plan the first paint without them.

## Merge policies

These tell the official client how to combine the new value with data it already has. Metadata is emitted only for selected paths.

| Constructor | Client behavior |
| --- | --- |
| `Merge` | Append to the retained value |
| `Prepend` | Prepend |
| `DeepMerge` | Recursively merge objects |
| `MatchOn(value, "id")` | Match items on a nested path before merging |

`X-Inertia-Reset` containing the prop path drops merge metadata so the client replaces the value.

## Once

`Once` keeps a resolved prop on the client across later pages. The client sends retained keys in `X-Inertia-Except-Once-Props`. The server then omits the value from `props` but still emits `onceProps` metadata.

```go
"countries": inertia.Once(
    loadCountries,
    inertia.OnceKey("geo.countries"),
    inertia.OnceFor(24*time.Hour),
),
```

| Option | Meaning |
| --- | --- |
| `OnceKey` | Stable retention key; the prop path is the default |
| `OnceExpiresAt` | Absolute expiry |
| `OnceFor` | Expiry relative to evaluation time |
| `ForceFresh` | Resolve even if the client still holds the key |

Partial selection or `ForceFresh` forces a new value. Once props are for slowly changing reference data, not per-user secrets that must not live in client history.

## Scroll

`Scroll(value, metadata)` adds infinite-scroll metadata and merges the nested `data` collection. `ScrollAt(value, "results", metadata)` chooses a different wrapper path.

Metadata may be `ScrollMetadata`, `ProvidesScrollMetadata`, or a function of the resolved value. `X-Inertia-Infinite-Scroll-Merge-Intent` selects append versus prepend. Reset paths mark `scrollProps[path].reset`.

```go
"products": inertia.Scroll(pageResult, pageResult),
```

Deferred scroll props announce merge metadata before they resolve so the follow-up reload applies consistently, but they do not emit scroll metadata until the value exists.

## Evaluation order

For each response the renderer:

1. Rejects shared or page props named `errors` or `errors.*`.
2. Copies shared props, then overwrites with page props.
3. Unpacks dotted keys.
4. Walks keys in sorted order.
5. Skips optional/deferred values on full visits, while still announcing deferred and once metadata.
6. Applies `only`, then `except`, unless `Always` inherited.
7. Runs resolvers only for included paths.
8. Records merge, scroll, deferred, rescued, shared, and once metadata.
9. Sets `props.errors` to an object, wrapping it in the named error bag when requested.

`sharedProps` lists top-level shared keys plus `errors`. Page props replace shared props of the same name.

## Application usage

Build a mapper next to the controller and keep it boring:

```go
func productResources(products []models.Product) []ProductData {
    out := make([]ProductData, 0, len(products))
    for _, product := range products {
        out = append(out, ProductData{
            ID:    product.ID,
            Name:  product.Name,
            Price: product.Price,
        })
    }
    return out
}
```

Put Always props that layouts need on every visit, defer secondary widgets, and keep the index query on the first response. If a page feels slow, inspect which resolvers run on full visits before adding more partial-reload machinery.

Generated `--inertia` controllers already use `FromStruct` and TypeScript types. Hand-written pages should follow the same pattern so `generate` output and application code stay interchangeable. See [Generators](/docs/latest/inertia-generators).

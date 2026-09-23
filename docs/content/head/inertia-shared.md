# Shared Data and Redirects

Shared props, validation errors, flash, and redirects are how Andurel keeps Inertia forms and navigation aligned with Echo. The adapter implements the v3 protocol; the application owns sessions, cookies, and domain errors.

## Shared props

`WithShared` copies static values at renderer construction. Generated apps share the build version:

```go
inertia.WithShared(inertia.Props{"appVersion": appVersion}),
```

`WithSharedProvider` adds request-scoped values. Providers run on every `Page` call, including partial reloads:

```go
inertia.WithSharedProvider(func(etx *echo.Context) (inertia.Props, error) {
    user := currentUser(etx)
    if user == nil {
        return inertia.Props{"account": nil}, nil
    }
    return inertia.Props{
        "account": AccountData{ID: user.ID, Email: user.Email},
    }, nil
}),
```

Page props replace shared props of the same key. `errors` cannot be supplied as a shared or page prop; it is reserved and always overwritten by the renderer.

Use shared props for layout data that almost every page needs. Do not put large collections there. A partial reload that does not select those keys still has to consider `Always` and shared replacement rules; expensive shared providers run even when the page then drops the value.

## Validation errors

`PageBuilder.ValidationErrors` sets the protected `props.errors` object. When the request sends `X-Inertia-Error-Bag`, the map is nested under that bag name.

Generated auth controllers redisplay the form on the same visit:

```go
if validationErrors, ok := validation.As(err); ok {
    return s.renderer.Page(etx, "Auth/Login", inertia.Props{}).
        ValidationErrors(validationErrors.ToMap()).
        Render()
}
```

`ValidationErrors.ToMap()` keeps the first message per field. The React adapter exposes `errors` as page props; generated login pages read `errors.email` next to `useForm`.

Named bags matter when two forms share a layout. The client sets `X-Inertia-Error-Bag`; the server must not invent a bag the adapter did not request.

Inertia can also redirect-after-write and rely on flash for non-field messages. Field errors belong in `props.errors`. Flash is the wrong channel for input validation.

## Flash

Flash is a top-level page field, not a prop. The renderer already installs a provider that reads `inertia.FlashFromContext`. Generated apps usually add a provider that reads kiks flashes so the page field stays in sync with the session bag:

```go
inertia.WithFlashProvider(func(etx *echo.Context) any {
    return kiks.Flashes(etx.Request().Context())
}),
```

You can also copy messages onto the request for a single response:

```go
ctx := inertia.ContextWithFlash(
    etx.Request().Context(),
    kiks.Flashes(etx.Request().Context()),
)
etx.SetRequest(etx.Request().WithContext(ctx))
```

`PageBuilder.Flash` overrides providers for one response. Extra `WithFlashProvider` callbacks run in order until one returns non-nil.

Generated `resources/js` trees read `page.flash` for toasts. The site adapter expects an array of `{type, message}` (and related) objects from `kiks.FlashMessage`. Your page can use a different shape as long as the client and the flash provider agree.

kiks middleware already persists the flash cookie on HTTP redirects (3xx), status 409, and `X-Andurel-Client-Redirect` responses. If you keep flashes only in `inertia.ContextWithFlash` without kiks, configure `SetReflashHandler` / `WithReflash` so redirect responses do not drop them. See [Cookies & Sessions](/docs/head/cookies-sessions).

Without a flash source on the next GET, a `Redirect` after `AddFlash` shows an empty flash on the following page.

## Redirects

After a successful write, do not render the next Inertia page from the POST. Redirect, then let the GET render.

```go
kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Product created")
return c.renderer.Redirect(etx, routes.ProductShow.URL(product.ID), http.StatusSeeOther)
```

`Redirect` requires a 3xx status. Middleware then:

- upgrades `302` after POST, PUT, PATCH, or DELETE to `303`
- rewrites fragment locations (`#...`) to a `409` response with `X-Inertia-Redirect`
- reflashes consumed messages

`Location` forces a full browser visit. Use it when the next document is not an Inertia page, or when session cookies must apply to a fresh document load. Generated login uses it to leave the auth screen:

```go
return s.renderer.Location(etx, routes.HomePage.URL())
```

For an Inertia client that is already on an Inertia page, `Location` returns `409` with `X-Inertia-Location` and removes `X-Inertia`. Ordinary browsers receive a `302`.

`FreshRedirect` is the same control-response pattern with `X-Inertia-Redirect`: the client performs a fresh Inertia GET instead of swapping from the current POST response.

Empty successful Inertia responses become a redirect to `Referer` or `/`. If a controller returns `nil` without writing a body, the user goes back rather than seeing a blank JSON page.

## Application usage

A typical mutating action:

1. Bind JSON from `useForm().post` / `.put` / `.delete`.
2. Validate and map to a model command.
3. On field errors, `Page` the form with `ValidationErrors`.
4. On domain failure, add an error flash and `Redirect` back to the form.
5. On success, add a success flash and `Redirect` or `Location`.

Generated Inertia resource controllers follow this pattern for create, update, and destroy. Auth controllers mix `ValidationErrors` redisplay with `Location` after login so the root document picks up the session cookie.

See [Pages and Visits](/docs/head/inertia-pages) for `Page` and [Generators](/docs/head/inertia-generators) for the scaffolded form flow.

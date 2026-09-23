# Cookies & Sessions

Andurel v2 uses `github.com/mbvlabs/andurel/pkg/kiks` for cookie jars, bagged session cookies, flash messages, and Echo middleware that loads and persists them on the HTTP path. Generated apps wire this under `router/cookies`.

## Jar and definitions

Construct a jar with keys, a session `Store`, and exactly one `NewSession` definition. Mark additional Plain/Signed/Encrypted cookies with `Bagged` so middleware can load them into the request bag. Unmarked named cookies stay native-only (`Read` / `Write` / `Clear`).

```go
keys := kiks.Keys{
    Authentication: sessionCfg.AuthenticationKey,
    Encryption:     sessionCfg.EncryptionKey,
}
store, err := kiks.NewCookieStore(keys)
if err != nil {
    return nil, err
}

return kiks.NewJar(keys, store,
    kiks.NewSession[*cookies.App](
        sessionCfg.Name,
        kiks.HTTPOnly(),
        kiks.Secure(appCfg.IsProduction()),
        kiks.MaxAge(sessionCfg.MaxAge),
    ),
    // kiks.Bagged(ConsentCookie(secure)), // bag / every request
    // CartCookie(secure),                  // native / Read-Write
)
```

Each Go type may be registered once per jar. Two same-shape cookies need distinct types. Payload types implement `kiks.Cookie` (`MarshalCookie` / `UnmarshalCookie`) and return raw bytes, never a Set-Cookie string.

Session payloads persist through `Store`. Generated apps use `kiks.NewCookieStore` (encrypted cookie body). Flashes never go through `Store`; the jar keeps them in a dedicated flash cookie named `<session>_flash`.

## Middleware

Register `jar.EchoMiddleware` early in the Echo stack so controllers and shared props can read the bag from `context.Context`. Skip static/asset prefixes when those paths must not touch cookies:

```go
e.Use(jar.EchoMiddleware(
    kiks.SkipPrefixes(routes.AssetsPrefix, routes.APIPrefix),
))
```

Middleware eager-loads the session and flash cookie. Other bagged cookies load lazily on first `Get` / `Exists`. Persist runs once before the response is committed. Persist failures fail the request instead of silently dropping cookie updates.

## Get, Set, Destroy

Bagged cookies and the session use typed helpers on the request context. Use the same pointer type you registered (for example `*cookies.App`):

```go
app, err := kiks.Get[*cookies.App](c.Request().Context())
if err != nil {
    return err
}
app.UserID = user.ID
if err := kiks.Set(c.Request().Context(), app); err != nil {
    return err
}
```

| Helper | Behavior |
| --- | --- |
| `Get[T]` | Bag value or zero when missing; errors on no bag / unknown / native-only type |
| `Exists[T]` | Present and not destroyed |
| `Set[T]` | Marks dirty; session also sets `sessionDirty` |
| `Destroy[T]` | Clears value and expires on persist; destroying the session clears flashes |

Native-only cookies use `kiks.Read` / `Write` / `Clear` with the injected `*kiks.Jar`. Mixing bag and native APIs for the same type returns `ErrNotBag` or `ErrBagCookie`.

## Flash messages

Flashes are the shared contract for Templ and Inertia:

```go
kiks.AddFlash(ctx, kiks.FlashSuccess, "Product saved.")
messages := kiks.Flashes(ctx)
```

| Type | Constant |
| --- | --- |
| success | `kiks.FlashSuccess` |
| error | `kiks.FlashError` |
| warning | `kiks.FlashWarning` |
| info | `kiks.FlashInfo` |

Added flashes are visible to whatever this response renders. Middleware persists the flash cookie on HTTP redirects (3xx), status 409, and responses that set `X-Andurel-Client-Redirect` (Datastar/SSE redirects). Other statuses clear pending flashes so they do not leak into unrelated pages.

See [Shared Data and Redirects](/docs/head/inertia-shared) for Inertia flash wiring and [Authentication](/docs/head/authentication) for identity flows built on these cookies.

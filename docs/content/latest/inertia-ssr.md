# SSR

SSR is optional and per page. The Inertia renderer always knows how to return a client-mount document. Calling `.SSR()` on a page builder asks a JavaScript renderer to produce HTML for that initial visit only. Client visits still receive JSON.

The web process never starts Node. `cmd/app` is an HTTP client. `cmd/ssr` owns the JavaScript process.

## Opt a page in

```go
return p.renderer.Page(etx, "Home", inertia.Props{}).SSR().Render()
```

Without `.SSR()`, the root document includes `PageScript` and `AppMount`. With `.SSR()`, the renderer posts the finished page object and, on success, inserts `SSRHead` and `SSRBody`.

The SSR body must contain `data-server-rendered="true"` and a `data-page=` script. Missing markers are treated as render failures.

`.SSR()` on a JSON client visit still builds the page object for the adapter; the HTML document path is what consumes the SSR response.

## Where the request goes

`NewRenderer` installs an SSR client unless `WithSSRRenderer` replaces it.

| Environment | Client | Endpoint |
| --- | --- | --- |
| `development` with a Vite dev URL | `ViteSSRRenderer` | `{ViteDevURL origin}/__inertia_ssr` |
| any other environment | `HTTPRenderer` | `{INERTIA_SSR_URL}/render` |

Development talks to the Vite plugin, not `cmd/ssr`. Production and other environments post JSON to `/render` and expect `{ "head": [...], "body": "..." }`.

`INERTIA_SSR_URL` is a client URL. It may be `localhost`, a container DNS name, or any `http`/`https` host. `INERTIA_SSR_LISTEN` is the bind address for Node and must be an IP or `localhost`, including `0.0.0.0` / `::`. Health probes loop back when the bind host is unspecified; they do not call `INERTIA_SSR_URL`.

## cmd/ssr

Generated `cmd/ssr` constructs `NewSSRRuntime` from `config.Inertia` and `assets.Files`:

```go
runtime, err := inertia.NewSSRRuntime(
    cfg.SSRRuntime,
    cfg.SSRBundle,
    cfg.SSRListen,
    cfg.SSRStartupTimeout,
    cfg.SSRMinimumMajor,
    cfg.SSRRequestTimeout,
    assets.Files,
)
if err != nil {
    return err
}
if err := runtime.Start(ctx); err != nil {
    return err
}
defer runtime.Stop(shutdownCtx)
```

`Start` looks up the executable, checks `node --version` against `INERTIA_SSR_MINIMUM_MAJOR` (default 22), spawns `node <bundle>` with `INERTIA_SSR_HOST` and `INERTIA_SSR_PORT`, and waits until `/health` returns `{ "status": "ok" }`. If the bundle path is missing on disk, the runtime extracts it from the optional `fs.FS` (generated apps pass `assets.Files`) into a temp file.

`Stop` POSTs `/shutdown` and kills the process when the context ends. `Errors()` reports unexpected child exits so a process manager can restart.

`andurel run` supervises `cmd/ssr` in development when SSR is used. In production, run `cmd/ssr` (or an equivalent Node host serving the same HTTP contract) next to `cmd/app`. The JavaScript package manager in `andurel.lock` does not choose this Node binary.

## Failures and fail-fast

When SSR fails, the renderer logs `inertia SSR fallback` and serves the client-mount document unless `WithSSRFailFast(true)` is set. Fail-fast is useful in release verification so a broken bundle cannot ship as a silent CSR page.

Generated scaffolds default `INERTIA_SSR_FAIL_FAST` to true. Set it false in production if you prefer availability over a hard error when Node is down.

Transport failures are `SSRTransportError` with kind `transport`, `status`, `encode`, `decode`, or `response`. Oversized bodies return `ErrResponseTooLarge`. Protocol-level wrapping uses `inertia.Error` with kind `ssr`. See [Diagnostics](/docs/latest/inertia-diagnostics).

Timeout and size bounds are independent: startup wait, per-request timeout, and `INERTIA_SSR_MAX_RESPONSE_BYTES`.

## Frontend server entry

Generated `resources/js/ssr.tsx` uses `@inertiajs/react/server` (or the Vue/Svelte equivalent), binds `INERTIA_SSR_HOST` / `INERTIA_SSR_PORT`, and renders with the same page resolver as the browser entry. Vite's Inertia plugin records that file as the SSR entry.

Production build is two Vite outputs: client assets under `assets/dist/vite` and `assets/dist/ssr/ssr.js`. Embed both. `cmd/app` reads the client manifest; `cmd/ssr` executes the SSR bundle.

## Application usage

Opt in page by page. Marketing and documentation benefit from SSR; authenticated CRUD often does not. A mixed app can SSR `Home` and client-render `Products/Index`.

Do not call `.SSR()` unless a renderer is reachable. In development that means Vite is running. In production that means `cmd/ssr` or another `/render` server is up, or fail-fast is off and fallback is acceptable.

Verify SSR separately from protocol tests: mode is no longer a renderer enum. Confirm listen URL, client URL, runtime version, bundle path, health, timeout, and response-size bound.

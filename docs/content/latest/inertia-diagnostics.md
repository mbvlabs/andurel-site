# Diagnostics

Inertia failures are typed so logs can say *what kind of protocol step failed* without dumping session data or prop values. Turn on protocol debug only when you need request classification.

## Protocol debug

`WithProtocolDebug(true)` or `INERTIA_PROTOCOL_DEBUG=true` logs:

- request method, URL, whether the visit is Inertia, partial component, `only` / `except` / `reset`, and purpose
- page component, URL, prop *keys*, merge paths, deferred groups, and once-prop counts

It does not log prop values, flash payloads, or cookies. Generated scaffolds default this to true for local work; keep it off in production unless you are reproducing a protocol bug.

## Error kinds

`inertia.Error` wraps the cause with:

| Field | Meaning |
| --- | --- |
| `Kind` | `protocol`, `props`, `root`, or `ssr` |
| `Operation` | Short step name such as `resolve` or `encode page` |
| `Component` | Page component when known |
| `Method`, `URL` | Request that failed |
| `PropPath` | Dotted prop path for evaluation failures |

Unwrap with `errors.As`. Handler logs should print `Kind`, `Operation`, `Component`, and `PropPath` before the cause.

Typical sources:

- **protocol** — invalid merge intent, version provider failure
- **props** — protected `errors` key, dotted-path conflict, resolver error, unsupported scroll metadata, JSON encode failure
- **root** — `WithRoot` returned nil, Templ render error
- **ssr** — missing SSR client, empty body, missing `data-server-rendered` / `data-page` markers, or a wrapped transport error when fail-fast is on

SSR HTTP problems also surface as `SSRTransportError` (`transport`, `status`, `encode`, `decode`, `response`) and `ErrResponseTooLarge`.

## What to inspect

When a visit misbehaves, compare the request you think you sent with the request the middleware parsed:

1. `X-Inertia` — missing means an initial document, not JSON.
2. Exact `component` string versus the file under `resources/js/Pages`.
3. `X-Inertia-Partial-Component` versus the rendered component. Filters are ignored on mismatch.
4. Nested `only` / `except` paths after dotted unpacking.
5. Response `Vary` containing `X-Inertia`.
6. `X-Inertia-Version` versus `buildPathURL` or `WithVersion`.
7. Redirect status: `303` after unsafe methods, `409` for location / fragment / version reloads.
8. Whether `props.errors` is an object and whether a named bag was requested.

For SSR, separately verify:

- development posts to Vite `/__inertia_ssr`, not `INERTIA_SSR_URL`
- production `cmd/app` posts to `INERTIA_SSR_URL/render`
- `cmd/ssr` binds `INERTIA_SSR_LISTEN` and the Node major meets `INERTIA_SSR_MINIMUM_MAJOR`
- bundle path or embedded `assets/dist/ssr/ssr.js`
- `/health` returns `{"status":"ok"}` on the listen loopback
- timeout and `INERTIA_SSR_MAX_RESPONSE_BYTES`
- fail-fast versus client fallback

## Application usage

Do not log `Props` maps in controllers. If you need to see which keys were evaluated, enable protocol debug. If a resolver fails, return the error and let `inertia.Error` carry the prop path.

For tests, construct a renderer with a stub root and assert status, `X-Inertia`, and decoded page keys without starting Vite. SSR tests belong against `HTTPRenderer` or `NewSSRRuntime` with a fixture bundle, not against the full web process.

# Root Document and Vite

The Inertia adapter renders Vue, React, or Svelte pages, but the first HTML document is application-owned Templ. `pkg/inertia` supplies the page object, Vite tags, and SSR insertion points. Layout inside the SPA still belongs in frontend components.

## Root document

Generated projects write `views/root.templ` and pass it to `WithRoot(views.Root)`:

```templ
templ Root(data inertia.RootData) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <title>{ data.ProjectName }</title>
            @inertia.SSRHead(data.SSR)
            @templ.Raw(string(data.ViteHead))
        </head>
        <body>
            if data.SSR != nil {
                @inertia.SSRBody(data.SSR)
            } else {
                @inertia.PageScript(data.ContainerID, data.PageJSON)
                @inertia.AppMount(data.ContainerID)
            }
            @templ.Raw(string(data.ViteBody))
        </body>
    </html>
}
```

`RootData` contains the resolved `Page`, raw `PageJSON`, `ContainerID`, `ProjectName`, `Environment`, trusted Vite HTML, and an optional `SSR` response.

Customize the root for fonts, analytics, CSP nonces, or extra meta tags. Do not move page layout here. Authenticated chrome, flash toasts, and nav belong in `resources/js/Layouts` so client visits can reuse them without reloading the document.

`PageScript` validates the JSON, escapes slashes, and writes:

```html
<script data-page="app" type="application/json">...</script>
```

`AppMount` writes `<div id="app"></div>`. `SSRHead` and `SSRBody` insert markup from the SSR runtime and replace both the script and the empty mount. Treat SSR HTML as trusted local output.

After editing Templ, regenerate Go:

```bash
andurel generate view
```

## Development assets

When `WithEnvironment` is not `"production"`, the renderer emits Vite development tags from `INERTIA_VITE_DEV_URL`:

- `<script type="module" src="{viteDevURL}/@vite/client">`
- the entry module, for example `{viteDevURL}/resources/js/app.tsx`
- React Refresh preamble when the entry point ends in `.tsx`

`andurel run` starts Vite beside the Go process. The JavaScript package manager in `andurel.lock` installs dependencies; it does not select the SSR Node binary.

If development tags 404, the Vite origin or entry point does not match `vite.config`. Keep `INERTIA_VITE_DEV_URL` and `INERTIA_ENTRY_POINT` aligned with the Vite root and `build.outDir`.

## Production assets

In production the renderer reads `dist/vite/manifest.json` from `WithAssetFS`. It looks up `entryPoint`, emits CSS `<link>` tags, then a module `<script>` for the hashed file. URLs are prefixed with `buildPathURL` after a trailing `*` is stripped, matching the Echo asset route.

`andurel build` installs frontend dependencies, runs the Vite client build, and embeds `assets/` in the Go binary. A missing manifest, unknown entry point, or nil asset filesystem fails `NewRenderer` at startup rather than serving a document without JavaScript.

The SSR bundle is a second Vite build, typically `assets/dist/ssr/ssr.js`. `cmd/app` does not execute it. See [SSR](/docs/latest/inertia-ssr).

## Asset version

Inertia uses `X-Inertia-Version` so a deployed client reloads when assets change.

The default version is `buildPathURL`. Override it with `WithVersion` for a release id, or `WithVersionProvider` when the version depends on the request:

```go
inertia.WithVersion(appVersion),
```

On Inertia GET requests, middleware compares the header with `currentVersion`. A mismatch returns `409` and `X-Inertia-Location` of the current URL, without the `X-Inertia` page-response header. The adapter then performs a full document visit and picks up new Vite tags.

POST, PUT, PATCH, and DELETE are not converted into version reloads. Finish the write with `Redirect` or `Location`, then let the following GET reload if needed.

## Frontend entry

Generated `resources/js/app.tsx` (or `app.ts`) calls `createInertiaApp`, resolves `./Pages/${name}.tsx`, and hydrates when `data-server-rendered="true"`. Flash toasts typically wrap the page tree so both initial and client visits can read the page `flash` field.

Page components receive JSON props as function arguments. Import generated TypeScript declarations instead of re-declaring the shape. See [Generators](/docs/latest/inertia-generators) and [Props](/docs/latest/inertia-props).

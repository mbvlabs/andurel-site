# Generators

Andurel generates an Inertia application you can edit. Framework behavior stays in `pkg/inertia`; routes, controllers, payloads, and page components are yours after creation.

Resource generation still defaults to Templ, even inside an Inertia scaffold. Pass `--inertia` when a resource should be a frontend page.

## Create an Inertia application

```bash
andurel new orbit --inertia vue
andurel new orbit --inertia react/pnpm
andurel new orbit --inertia svelte/bun
```

Adapters are `vue`, `react`, and `svelte`. Optionally append `/npm`, `/pnpm`, `/bun`, or `/yarn` (default npm). That package manager is stored in `andurel.lock` and used to install and build frontend assets. It does not select the Node binary for `cmd/ssr`.

The scaffold writes, among other things:

| Area | Files |
| --- | --- |
| Renderer | `cmd/app` `newInertia`, `config/inertia.go` |
| SSR process | `cmd/ssr`, `resources/js/ssr.tsx` (or `.ts`) |
| Root document | `views/root.templ` |
| Client entry | `resources/js/app.tsx` or `app.ts`, Vite config, `package.json` |
| Layout and flash | `resources/js/Layouts`, flash toast helpers |
| Auth and errors | Inertia pages under `resources/js/Pages/Auth` and `Pages/Errors` |
| TypeScript routes | `resources/js/routes.ts` |

The welcome page remains Templ. Auth, confirmation, password reset, and default error pages use Inertia. Controllers import `github.com/mbvlabs/andurel/pkg/inertia` directly.

After scaffolding, install JavaScript dependencies with the recorded package manager, copy `.env`, and run `andurel run`. That process starts Go, Vite, and `cmd/ssr` as needed.

## Wire-up the generators already emitted

Generated `newInertia` calls `NewRenderer` with `views.Root`, embedded `assets.Files`, project name, environment, protocol debug, `appVersion` as a shared prop, and fail-fast from config. The router registers `renderer.Middleware()`, copies flashes with `inertia.ContextWithFlash`, and sets `SetReflashHandler`.

You change application behavior by editing those files, not by forking the package. Add `WithSharedProvider` next to `newInertia` for the current account. Opt individual pages into SSR with `.SSR()`. See [Renderer](/docs/latest/inertia-renderer) and [SSR](/docs/latest/inertia-ssr).

## Generate Inertia resources

```bash
andurel generate scaffold Product --inertia
andurel generate controller Product --inertia
andurel generate controller Product index show --inertia
andurel generate controller Dashboard overview --inertia
```

`--inertia` reads the adapter from `andurel.lock`. It cannot combine with `--api`. Without `--inertia`, views are Templ even in an Inertia project.

A generated Inertia resource includes:

- controller methods that inject `*inertia.Renderer`
- `ProductData`, `ProductIndexProps`, and `ProductItemProps` with JSON tags
- `newProductData` mapping from `models.Product` (never the Bun model itself)
- `Page` calls with `inertia.FromStruct(...)`
- create/update/destroy that flash and `Redirect`
- page components under `resources/js/Pages/Product/`
- TypeScript declarations under `resources/js/types/`
- route values marked with `routing.InertiaRoute()` when they belong in `routes.ts`

Index example:

```go
return c.renderer.Page(
    etx,
    "Product/Index",
    inertia.FromStruct(ProductIndexProps{
        Items: newProductDataList(list.Products),
    }),
).Render()
```

The React page imports the generated type and `routes`:

```ts
import type { ProductData, ProductIndexProps } from '@/types/product'
import { routes } from '@/routes'

export default function Index({ items }: ProductIndexProps) {
  return <Link href={routes.productShow(routeID(item))}>View</Link>
}
```

Create and edit pages use `useForm` (or the Vue/Svelte equivalent) and POST/PUT JSON to the generated route helpers. Field errors belong in `props.errors`; generated auth shows that pattern with `ValidationErrors`. Resource create currently redirects with flash on failure rather than redisplaying field errors—add `ValidationErrors` yourself when the form needs them.

Custom actions (`generate controller Dashboard overview --inertia`) add an empty controller method, a page component, and a GET route. Fill in props the same way as a hand-written page.

## TypeScript route helpers

```bash
andurel generate routes
```

This reads `router/routes/*.go` and writes `resources/js/routes.ts`. Only routes constructed with `routing.InertiaRoute()` are exported. That keeps asset and internal endpoints out of the browser bundle.

```go
var ProductShow = routing.NewRouteWithUUIDID(
    "/:id",
    "show",
    "/products",
    routing.InertiaRoute(),
)
```

After adding or changing an Inertia route, regenerate helpers and compile the frontend. `andurel doctor` fails when `routes.ts` is missing or stale. `andurel routes --json` is the same manifest if you need it from other tools.

Import `routes` instead of hard-coding paths. A parameter change is a cross-boundary API change.

## Mixing Templ and Inertia

One Echo app can serve both. Keep each route's contract consistent: an Inertia visit must receive an Inertia page or a location/redirect control response, not arbitrary HTML.

A clean split is public Templ marketing pages beside an authenticated Inertia application. Generating a Templ resource inside an Inertia app is the default and is useful for admin HTML that should not load the SPA.

Do not share Templ components with Vue/React/Svelte layouts. One renders in Go; the other renders in the client runtime.

## Frontend files you own

Edit page components, layouts, and CSS freely. Re-running a generator will not reconstruct a file you already customized unless you remove it first. Payload structs and TypeScript declarations should stay in sync: if you add a JSON field on `ProductData`, add it to the `.ts` declaration or change the generator inputs and regenerate.

`andurel generate view` compiles Templ, including `views/root.templ`. It does not compile Vite. `andurel build` installs JS dependencies and builds client and SSR bundles.

## Checklist after generating

1. Map domain entities to payload structs; never pass models into `Page`.
2. Mark browser-facing routes with `InertiaRoute()` and run `generate routes`.
3. Use `Redirect` or `Location` after writes; use `ValidationErrors` for field errors.
4. Keep component names aligned with `resources/js/Pages/...`.
5. Add `.SSR()` only on pages that should hit Vite or `cmd/ssr`.
6. Run `andurel doctor` before committing when `routes.ts` or generated views changed.

See [Frontend Options](/docs/latest/frontend-options) for choosing Templ versus Inertia, and [Code Generation](/docs/latest/code-generation) for non-Inertia generators.

# Routing Package

`github.com/mbvlabs/andurel/pkg/routing` provides typed route values for constructing paths and URLs. It does not register Echo handlers or select HTTP methods; generated files in `router/routes` declare route values, while `router.Router` binds them to controllers and middleware.

## Route declarations

A route stores its local path, name, optional prefix, and whether it should appear in generated Inertia helpers.

```go
var ProductsIndex = routing.NewSimpleRoute("", "index", "/products")
var ProductShow = routing.NewRouteWithUUIDID(
    "/:id", "show", "/products", routing.IncludeInInertia(),
)
```

`Path()` returns the Echo registration path, `Name()` returns the prefixed name, `URL(...)` substitutes typed parameters, and `FullURL(base, ...)` prepends the application base URL. Central declarations prevent links, redirects, and tests from duplicating strings.

## Typed parameter routes

| Constructor | Placeholder | Go argument |
| --- | --- | --- |
| `NewSimpleRoute` | none | none |
| `NewRouteWithUUIDID` | `:id` | `uuid.UUID` |
| `NewRouteWithSerialID` | `:id` | `int32` |
| `NewRouteWithBigSerialID` | `:id` | `int64` |
| `NewRouteWithStringID` | `:id` | `string` |
| `NewRouteWithSlug` | `:slug` | `string` |
| `NewRouteWithToken` | `:token` | `string` |
| `NewRouteWithFile` | `:file` | `string` |
| `NewRouteWithParams[T]` | tagged fields | typed struct |

Multiple parameters use a struct so call sites cannot swap same-typed values:

```go
type ProjectParams struct {
    Team    string `param:"team"`
    Project string `param:"project"`
}
var ProjectShow = routing.NewRouteWithParams[ProjectParams](
    "/:team/projects/:project", "show", "projects",
)
url := ProjectShow.URL(ProjectParams{Team: "platform", Project: "orbit"})
```

The generic route substitutes tagged string fields. A non-struct input returns an empty URL. Prefer dedicated constructors for numeric and UUID IDs.

## Prefixes and normalization

Prefixes affect path and name. Joining tolerates a missing or doubled slash and logs a warning before normalization. Paths are trimmed, backslashes converted, embedded query strings and fragments removed, double slashes collapsed, a leading slash added, and trailing slashes removed except for `/`.

Normalization is a guardrail. Declare canonical paths so warnings reveal mistakes rather than becoming routine output.

## Query parameters and full URLs

`QueryParam` appends query values; `FullURL` combines a route with `config.App.BaseURL`:

```go
next := ProductsIndex.URL(routing.QueryParam("page", "2"))
absolute := ProductShow.FullURL(app.BaseURL, product.ID)
```

Values are appended as supplied, so encode arbitrary input before passing it. Route declarations should not include query strings or fragments.

## Inertia TypeScript generation

`IncludeInInertia()` marks a route for `andurel generate routes`, which writes typed helpers to `resources/js/routes.ts`. Opt-in avoids exporting internal endpoints or assets into the browser bundle. `JsExpr` marks a JavaScript expression for generator interpolation; it is not a general escaping function.

After changing a marked route, regenerate helpers and compile the frontend. A parameter or path change is a cross-boundary API change.

## Registration boundary

The package does not verify registration, attach middleware, reverse-match requests, or parse Echo parameters. Register `route.Path()` with the intended method and parse request parameters in controllers. Important route tests should cover both URL construction and router registration.

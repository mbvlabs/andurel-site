# Hypermedia

`github.com/mbvlabs/andurel/pkg/hypermedia` joins compiled Templ components to Echo responses and provides Datastar-compatible server-sent events. It does not prescribe component structure or store UI state; the application owns its HTML, fragment keys, signals, and event flow.

## Rendering Templ

`RenderPage` and its readability alias `RenderComponent` buffer a component and write an HTML response. `WithStatus` changes the response status. `RenderHTML` renders without an HTTP response, while `RenderFragment` and `RenderFragments` extract named Templ fragments.

```go
return hypermedia.RenderPage(
    etx,
    views.ProductsIndex(products),
    hypermedia.WithStatus(http.StatusOK),
)
```

Rendering errors are returned before a partial response is written. Keep data loading outside components so failures remain ordinary controller errors.

## Datastar actions and signals

`DataAction` constructs a Datastar action expression for an HTTP method and route. Options add request headers, filter signals, or keep the connection open. On the server, `ReadSignals` decodes the `datastar` query value for GET requests and the JSON body for other methods.

```go
type RenameSignals struct { Name string `json:"name"` }

var signals RenameSignals
if err := hypermedia.ReadSignals(etx.Request(), &signals); err != nil {
    return err
}
```

Treat signals as request input: validate and authorize them before changing state. They are not a trusted server-side store.

## Element patches

`PatchComponent` renders Templ and emits a `datastar-patch-elements` SSE event. Use `PatchHTML` for already rendered content, or `PatchFragment` and `PatchFragments` for named component fragments.

Patch modes cover outer and inner replacement, removal, prepend, append, insertion before or after, and replacement. `WithSelector` targets CSS selectors; `WithSelectorID` targets an ID. View transitions, event IDs, and retry durations are explicit options.

```go
return hypermedia.PatchComponent(
    etx,
    views.ProductRow(product),
    hypermedia.WithSelectorID("product-"+product.ID.String()),
    hypermedia.WithModeOuter(),
    hypermedia.WithViewTransitions(),
)
```

`RemoveElement`, `RemoveElementByID`, and `RemoveElementf` emit removal patches without requiring empty placeholder HTML.

## Signal, script, and browser events

`PatchSignal` and `PatchSignals` send JSON signal changes. `WithOnlyIfMissing` initializes state without overwriting an existing browser value. `MergeSignals` supports the older merge-signals event when interoperability requires it.

`ExecuteScript` sends a transient script element. `Redirect`, `ReplaceURL`, `ReplaceURLQuery`, and `Prefetch` are higher-level helpers built on that channel. Prefer element and signal patches for normal UI updates; scripts are an escape hatch for behavior without a declarative Datastar event.

`DispatchCustomEvent` serializes a detail value and dispatches a browser `CustomEvent`. Options control selector, bubbling, cancelability, composed behavior, event ID, and retry duration.

## Long-lived streams

Single patch functions open an SSE response, send an event, and flush it. For multiple events, use `NewBroadcaster`, which opens the stream once and serializes concurrent writes with a mutex.

```go
b, err := hypermedia.NewBroadcaster(etx)
if err != nil { return err }

for update := range updates {
    if b.IsClosed() { return nil }
    if err := b.PatchComponent(views.JobStatus(update)); err != nil {
        return err
    }
}
```

The broadcaster mirrors patch, signal, script, event, redirect, URL, and prefetch operations. Always observe the request context or `IsClosed`; a disconnected browser must stop upstream work. Event IDs support browser reconnection policy, but replay storage and resumption semantics remain application responsibilities.

## Response and proxy behavior

SSE helpers set `Content-Type: text/event-stream` and `Cache-Control: no-cache`, flush headers, write complete event frames, and flush each event. HTTP/1 responses also request keep-alive for broadcasters.

Reverse proxies must allow streaming and avoid buffering these routes. The server write timeout must accommodate the longest intended stream. Heartbeats, authorization expiry, fan-out, backpressure, and replay are application concerns rather than hidden package behavior.

## Failure behavior

Errors are wrapped with `hypermedia:` operation context for rendering, JSON encoding, body reads, writes, and flushes. Once streaming begins, the status and headers are committed; return errors for logging and cancellation, not for rendering an HTML error page on the same response.

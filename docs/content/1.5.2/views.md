# Views

Templ is Andurel's default type-safe view layer. Datastar can update named fragments and consume server-sent events for interactive pages.

## Render a page

Define components in `views/*.templ` and render them from a controller:

```templ
templ ProductsIndex(products []models.ProductEntity) {
    @base(SetTitle("Products")) {
        <main>
            <h1>Products</h1>
        </main>
    }
}
```

Generate Go code after changing Templ files:

```bash
andurel generate view
```

## Render a fragment

Named Templ fragments let a controller return only the HTML that changed. Use this for Datastar-driven interactions instead of rebuilding client state.

## Inertia pages

Inertia-enabled projects store pages under the selected frontend source directory. Controllers call the Inertia renderer with serializable props while Go continues to own routing and application work.

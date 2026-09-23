# TypeScript Sync

Inertia projects keep Go routes and controller payloads as the source of truth. `andurel sync` writes matching TypeScript helpers without creating application-owned Go files.

## Routes

```bash
andurel inspect routes --json
andurel sync routes
```

`sync routes` writes `resources/js/routes.ts` from the route manifest. Import helpers via `@/routes` in pages and layouts. Mark routes that should appear in the TypeScript helpers with `routing.InertiaRoute()` when declaring them.

## Payloads

```bash
andurel sync payloads
```

`sync payloads` writes TypeScript types for controller payload / bind structs under `resources/js/types/`. Keep JSON tags stable and avoid serializing database rows directly — see [Props](/docs/head/inertia-props).

## Related commands

Scaffolding pages and controllers remains under [Generators](/docs/head/inertia-generators) and [generate](/docs/head/generate). Use [sync](/docs/head/sync) for the full derived-artifact surface.

# sync

`andurel sync` refreshes **derived** files from a source of truth. It does not create application-owned domain files — that is [generate](/docs/head/generate).

## Commands

```bash
andurel sync views
andurel sync queries
andurel sync routes --json
andurel sync payloads --json
andurel sync email
andurel sync factory User --check
andurel sync factories --sync
```

| Command | Source → output |
| --- | --- |
| `views` | `.templ` → `*_templ.go` (and email compile when inputs exist) |
| `queries` | `models/queries/*.sql` → generated Go clients |
| `routes` | route declarations → `resources/js/routes.ts` |
| `payloads` | controller payload structs → `resources/js/types/` |
| `email` | authored email templates → inlined renderers |
| `factory` / `factories` | model Entity → factory declarations |

`andurel tool sync` downloads pinned binaries. It is not this group.

See [TypeScript Sync](/docs/head/inertia-typescript-sync) for the Inertia-focused route and payload workflow.

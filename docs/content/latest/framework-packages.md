# Framework Packages

Andurel v2 extracts reusable infrastructure into independently versioned Go modules. New applications pin a set of package versions verified by the framework release.

## Available packages

```text
github.com/mbvlabs/andurel/pkg/email
github.com/mbvlabs/andurel/pkg/hypermedia
github.com/mbvlabs/andurel/pkg/inertia
github.com/mbvlabs/andurel/pkg/routing
github.com/mbvlabs/andurel/pkg/server
github.com/mbvlabs/andurel/pkg/storage
github.com/mbvlabs/andurel/pkg/validation
```

Each package has its own module, semantic version, and changelog. A framework v2 project may therefore use package versions whose numbers differ from one another and from the CLI.

## Ownership boundary

Packages expose ordinary constructors and configuration primitives. The application owns environment variable names, policy, Fx composition, routes, controllers, request metadata, models, and views. Packages do not use Fx as a service locator and do not hide database handles in contexts or globals.

`go.mod` is the source of truth for package dependencies. `andurel.lock` records the framework version, scaffold choices, extensions, and tools used to produce and maintain the project.

## Updating packages

Review each package changelog and update the dependency normally with Go tooling. Do not assume a CLI update automatically changes independently pinned package versions in existing application code. Test the web process, queue process, migrations, and any Inertia SSR mode together after an infrastructure update.

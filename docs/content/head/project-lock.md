# Project Lock

Andurel records scaffold choices and tool pins in two project files: `andurel.toml` (human-edited manifest) and `andurel.lock` (CLI-owned digests). Keep both under version control.

## andurel.toml

The manifest is the intent generators and upgrades respect:

```toml
schemaVersion = 1
version = 'latest'

[project]
  name = 'orbit'
  inertia = 'vue'
  javascriptPackageManager = 'pnpm'
  javascriptSSRRuntime = 'node'

[database]
  engine = 'postgresql'
  nullType = 'pgtype.Null'

[tools]
  narsilc = 'v0.4.4'
  templ = 'v0.3.1020'
  goose = 'v3.27.1'
  # ...
```

| Field | Meaning |
| --- | --- |
| `schemaVersion` | Manifest schema; currently `1` |
| `version` | Framework scaffold / lock lineage marker (often `latest` on new apps) |
| `project.name` | Application identity used by generators and cookie naming |
| `project.inertia` | Inertia adapter (`react`, `vue`, `svelte`) or empty for Templ |
| `project.javascriptPackageManager` | `pnpm`, `bun`, or `npm` for Vite builds |
| `project.javascriptSSRRuntime` | SSR executable expectation (typically `node`) |
| `database.engine` | SQL engine; only `postgresql` (and legacy `postgres`) |
| `database.nullType` | Nullable column strategy: `pgtype.Null` or `pointer` |
| `tools.*` | Pinned versions for managed binaries (`narsilc`, `templ`, `goose`, …) |

Edit the manifest deliberately. Changing `nullType` or Inertia adapter affects generation; bump tool versions with `andurel tool set-version` when you intend to pin a new binary.

## andurel.lock

The lock stores per-platform SHA-256 digests for managed tool downloads. It is TOML shaped like:

```toml
[[hashes]]
  tool = 'narsilc'
  version = 'v0.4.4'
  platform = 'linux/amd64'
  sha256 = '0b57f97ecf4813bc60e361cd5cd5dc806b41bb434defd7cb399754f11d7c1c66'
```

`andurel tool sync` downloads binaries listed in the manifest and verifies digests from the lock into `bin/`. Do not hand-edit digests unless you are restoring a known-good lock from the CLI.

## How tooling reads them

`layout.ReadLockFile` assembles both files into one in-memory lock: scaffold and database settings plus tool metadata. Error messages that mention `andurel.lock` usually mean that combined project pad is missing or invalid.

`go.mod` remains the source of truth for Go module versions. Updating the CLI does not rewrite application imports. Use `andurel packages list` / `andurel packages update` deliberately. See [packages](/docs/head/packages).

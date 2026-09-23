# Testing

Andurel v2 reworks application testing around explicit constructors, `storage.Connection`, and isolated PostgreSQL databases.

## Database tests

Use `storage.NewTestCluster` (or the generated test helpers) to start a shared PostgreSQL 17 Alpine cluster and create isolated databases per package or test. Prefer one cluster per package over a container per case.

Apply embedded migrations before exercising models. See [Getting Started](/docs/head/database).

## Unit and HTTP tests

Construct collaborators through ordinary Go constructors or Fx test apps. Prefer fakes at interface boundaries (email senders, clocks) over global state. Exercise Echo handlers with the same middleware stack the web process uses when behavior depends on kiks cookies or Inertia headers.

## Factories

Build entities with `models/factories` and keep them synchronized with `andurel sync factories --check`. See [Factories](/docs/head/factories).

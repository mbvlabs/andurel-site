# Introduction

Andurel is a Rails-like web framework for Go. The version documented here follows the `master` branch and is the development line for Andurel v2.

## What v2 emphasizes

Andurel v2 keeps generated application code explicit while moving reusable infrastructure into independently versioned packages. A new application uses Fx for dependency injection and lifecycle management, Bun for ordinary PostgreSQL persistence, sqlc for complex queries, River for jobs, and either Templ with Datastar or Inertia v3 for its UI.

The framework still favors one-time generation: a scaffold creates models, factories, controllers, routes, and views that belong to your application and can be edited normally.

## Development status

The `latest` documentation describes current `master`, not the stable v1 release. APIs and generated layouts may change before `v2.0.0`. Use the [1.5.2 documentation](/docs/1.5.2/introduction) for an existing v1 application.

There is no automated upgrade path from v1 to v2. Create new v2 projects with the v2 CLI and migrate existing applications deliberately.

## Platform and requirements

The v2 branch requires Go 1.27 and targets Linux and macOS on amd64 and arm64. PostgreSQL is the supported database.

Continue to [Installation](/docs/latest/installation) to install the development CLI and create a project.

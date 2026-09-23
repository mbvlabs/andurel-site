# Introduction

Andurel is a Rails-like web framework for Go. It favors useful conventions, complete generated resources, and a fast path from a new project to working product code.

## Why Andurel

Andurel removes repetitive setup while leaving generated application code in your repository. A scaffold creates the model, factory, controller, routes, and views needed for a complete resource. You can then change that code like any other part of your application.

The default stack includes Echo, PostgreSQL, Bun, Templ, Datastar, River, OpenTelemetry, and Fx. Inertia with Vue, React, or Svelte is optional.

## When to use it

Choose Andurel when you are building a full-stack Go application and value development speed, conventional project structure, and server-rendered or Inertia-based user interfaces.

Andurel v1 supports Linux and macOS on amd64 and arm64. Windows is not currently supported.

This is the current stable release. It matches `go install github.com/mbvlabs/andurel@latest`. The [head documentation](/docs/head/introduction) covers the v2 development line on `master`.

## Next steps

Continue to [Installation](/docs/1.5.5/installation) to install the CLI and create a project.

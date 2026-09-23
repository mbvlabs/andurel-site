# Build and Deploy

Andurel compiles templates, styles, optional Vite assets, and the Go application into a production release.

## Build

```bash
andurel build --version v1.2.3
```

The build generates Templ code, minifies Tailwind CSS, builds Vite assets when Inertia is enabled, downloads Go dependencies, and compiles a static Linux binary.

## Environment

Provide production database, session, token, email, CORS, and telemetry settings through the deployment environment. Set the application environment to production so secure cookie and server behavior is enabled.

## Database changes

Apply migrations as a deliberate release step before traffic reaches code that depends on the new schema. Back up important data before destructive changes.

## Verify releases

Official Andurel CLI releases publish checksums, SBOMs, signatures, and GitHub attestations. Verify downloaded archives before installing them in trusted build environments.

## Health and shutdown

The generated server handles termination signals and performs a bounded graceful shutdown. Ensure the deployment platform allows enough time for in-flight requests and queue workers to stop.

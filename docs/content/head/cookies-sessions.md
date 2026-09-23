# Cookies & Sessions

Andurel v2 uses `github.com/mbvlabs/andurel/pkg/kiks` for cookie jars, bagged session cookies, flash messages, and middleware that loads and persists them on the HTTP path.

## Jar and definitions

Construct a jar with signed/encrypted cookie definitions. Mark bagged types with `Bagged` (or use `NewSession`) so middleware can load them into a request bag. Native-only cookies stay available through `Read` / `Write` / `Clear` without entering the bag.

## Middleware

Register kiks middleware early in the Echo stack so controllers and Inertia shared props can read the current session and flashes from context. Persist failures fail the request instead of silently dropping cookie updates.

## Flash messages

Flashes live in a dedicated flash cookie beside the session cookie. Prefer the package flash helpers so Inertia and Datastar redirects see the same messages. See [Shared Data and Redirects](/docs/head/inertia-shared) for Inertia flash wiring.

## Authentication

Generated authentication flows build on these cookies. See [Authentication](/docs/head/authentication) for identity, login, and account routes.

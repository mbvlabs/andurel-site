# Authentication

A generated v2 application includes registration, sessions, email confirmation, password reset, and identity services.

## Generated flow

Account routes live under `/users`. Controllers own HTTP decisions, services coordinate workflows, plural model APIs such as `models.Users` and `models.Tokens` own persistence, and the email package delivers confirmation and reset messages.

These dependencies are supplied explicitly through Fx constructors. Request metadata and flash messages use typed helpers in the application-owned `router/appctx` package; database handles and services are never stored there.

## Sessions

Sessions use authenticated and encrypted cookies with `HttpOnly`, `SameSite=Lax`, and `Path=/`. Production enables secure cookies. The generated `config.Session` validates hexadecimal authentication and encryption keys and sets the lifetime from `SESSION_MAX_AGE`.

## CSRF and CORS

Unsafe cookie-authenticated requests remain CSRF protected. The default `header_only` strategy requires Fetch Metadata. Use `header_or_legacy_token` only when legacy forms must submit `_csrf` or `X-CSRF-Token`.

Credentialed CORS trusts the application base URL and explicitly configured origins. Wildcard origins are rejected.

## Authorization

Authentication establishes identity only. Add authorization checks for every protected resource and action before calling a model or service.

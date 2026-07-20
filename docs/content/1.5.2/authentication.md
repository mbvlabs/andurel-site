# Authentication

A generated Andurel application includes registration, sessions, email confirmation, and password reset flows.

## Generated flow

Account routes live under `/users`. Controllers handle browser requests, services own authentication workflows, models persist users and tokens, and email components deliver confirmations and reset links.

## Sessions

Sessions use authenticated and encrypted cookies with `HttpOnly` and `SameSite=Lax`. Production enables secure cookies. Rotate session and token keys through environment configuration.

## CSRF protection

Browser mutations require a valid CSRF token. Keep CSRF middleware enabled and use explicit trusted origins. API requests using application cookies remain subject to CSRF protection.

## Authorization

Authentication establishes identity. Add authorization checks for each protected resource and action instead of treating a signed-in session as universal access.

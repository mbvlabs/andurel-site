# Authentication

A generated v2 application includes registration, sessions, email confirmation, password reset, and an `services.Identity` workflow. Controllers own HTTP. Identity owns peppered passwords, tokens, and queue inserts for mail. Plural model APIs such as `models.Users` and `models.Tokens` own persistence. Session state and flashes use **kiks**.

Account routes live under `/users` (`router/routes/users.go`). Controllers register as Fx constructors with injected `services.Identity` and, for Inertia apps, `*inertia.Renderer`.

## Fx wiring

Generated roots compose these pieces:

| Piece | Role |
| --- | --- |
| `config.NewAuth` / `config.NewSession` | Pepper, token signing key, session cookie keys |
| `cookies.Module` (`NewJar`) | kiks jar with `NewSession[*cookies.App]` |
| `jar.EchoMiddleware` | Loads bag and flashes onto `context.Context` |
| `services.NewIdentity` | Users, tokens, queue insert, auth/mail config |
| Controllers (`Sessions`, `Registrations`, `Confirmations`, `ResetPasswords`) | HTTP adapters |

```go
func NewIdentity(
    db storage.Connection,
    users models.Users,
    tokens models.Tokens,
    insertOnly storage.InsertQueue,
    appCfg config.App,
    authCfg config.Auth,
    mailCfg config.Mail,
) Identity
```

Identity never reads environment variables. Controllers never store models or sessions on Echo context keys. See [Cookies & Sessions](/docs/head/cookies-sessions) and [Configuration](/docs/head/configuration).

## Sessions: login and logout

Inertia session controller (generated shape):

```go
type Sessions struct {
    identity services.Identity
    renderer *inertia.Renderer
}

func NewSessions(identity services.Identity, renderer *inertia.Renderer) Sessions {
    return Sessions{identity: identity, renderer: renderer}
}

func (s Sessions) New(etx *echo.Context) error {
    return s.renderer.Page(etx, "Auth/Login", inertia.Props{}).Render()
}

func (s Sessions) Create(etx *echo.Context) error {
    ctx, span := telemetry.From(etx, "sessions.create")
    defer span.End()

    var payload struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := etx.Bind(&payload); err != nil {
        return s.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
    }

    user, err := s.identity.AuthenticateUser(ctx, services.LoginData{
        Email:    payload.Email,
        Password: payload.Password,
    })
    if err != nil {
        if validationErrors, ok := validation.As(err); ok {
            return s.renderer.Page(etx, "Auth/Login", inertia.Props{}).
                ValidationErrors(validationErrors.ToMap()).
                Render()
        }
        return s.renderer.Redirect(etx, routes.SessionNew.URL(), http.StatusSeeOther)
    }

    if err := kiks.Set(etx.Request().Context(), &cookies.App{
        UserID:          user.ID.String(),
        IsAdmin:         user.IsAdmin,
        IsAuthenticated: true,
    }); err != nil {
        return err
    }

    return s.renderer.Location(etx, routes.HomePage.URL())
}

func (s Sessions) Destroy(etx *echo.Context) error {
    ctx, span := telemetry.From(etx, "sessions.destroy")
    defer span.End()

    if err := kiks.Destroy[*cookies.App](ctx); err != nil {
        return err
    }
    return s.renderer.Redirect(etx, routes.SessionNew.URL(), http.StatusSeeOther)
}
```

Read the session later with the same pointer type:

```go
app, err := kiks.Get[*cookies.App](c.Request().Context())
if err != nil {
    return err
}
```

`middleware.AuthOnly` redirects guests to `routes.SessionNew` when `app == nil` or `!app.IsAuthenticated`.

Flash after a successful write (for example a password update) with kiks, not a custom flash API:

```go
kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Password updated.")
```

Inertia surfaces those messages through `WithFlashProvider` reading `kiks.Flashes`. See [Shared Data and Redirects](/docs/head/inertia-shared).

## Registration, confirmation, and reset

| Flow | Routes (prefix `/users`) | Service entry | Outcome |
| --- | --- | --- | --- |
| Register | `GET /sign-up`, `POST /users` | `RegisterUser` | Redirect to confirmation; email job queued |
| Confirm | `GET /confirmation/new`, `POST /confirmation` | `VerifyEmail` | Sets `cookies.App`, `Location` home |
| Reset request | `GET /password/new`, `POST /password` | request token + email | Redirect back with flash |
| Reset apply | `GET /password/:token/edit`, `PUT/PATCH /password` | update password | Session optional; flash success |

Registration Create binds email/password/confirmPassword, maps `validation.As` errors onto `Auth/Registration`, and on success redirects to confirmation without signing the user in. Confirmation Create verifies the code, then `kiks.Set`s the session the same way login does.

Password reset issues a signed token (`TOKEN_SIGNING_KEY`), queues transactional mail with a reset URL, and consumes the token when the user submits a new password. Keep pepper rotation (`PEPPER` / `PREVIOUS_PEPPERS`) in deployment secrets. See [Configuration](/docs/head/configuration).

## CSRF and CORS

Unsafe cookie-authenticated requests stay CSRF protected.

| Setting | Behavior |
| --- | --- |
| `CSRF_STRATEGY=header_only` (default) | Requires Fetch Metadata (`Sec-Fetch-Site` / related) for unsafe methods |
| `CSRF_STRATEGY=header_or_legacy_token` | Also accepts `_csrf` form fields or `X-CSRF-Token` |
| `CSRF_TRUSTED_ORIGINS` | Extra trusted origins beyond the app base URL |
| `CORS_ALLOWED_ORIGINS` | Explicit credentialed origins; wildcards rejected |

Inertia and fetch clients send credentials with same-origin requests. Legacy HTML forms that post without Fetch Metadata need the compatibility strategy or a token field.

Credentialed CORS trusts the application base URL (`PROTOCOL` + `DOMAIN`) plus `CORS_ALLOWED_ORIGINS`. Do not enable credentialed `*`.

## Authorization

Authentication only proves identity. Gate every protected resource with middleware (`AuthOnly`) and explicit checks (ownership, admin) before model or service calls.

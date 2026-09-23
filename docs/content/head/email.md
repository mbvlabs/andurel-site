# Email

Andurel v2 provides `github.com/mbvlabs/andurel/pkg/email` for transport-neutral messages, application-owned Templ templates, Tailwind-to-inline-style compilation, and separate transactional and marketing sender interfaces.

## Sender interfaces

Consumers depend on capability-specific interfaces. The composition root selects `EMAIL_PROVIDER` and publishes these interfaces:

```go
type TransactionalSender interface {
    SendTransactional(context.Context, TransactionalPayload) error
}
type MarketingSender interface {
    SendMarketing(context.Context, MarketingPayload) error
}
```

Generated apps default `EMAIL_PROVIDER` to `mailpit` for local capture. Wire a production transport in the composition root that implements the same interfaces.

## Generate and compile email

```bash
andurel generate email WelcomeEmail
andurel sync email
```

That creates `email/welcome_email.templ` with a `Transformer`:

```go
type WelcomeEmail struct {
}

var _ Transformer = (*WelcomeEmail)(nil)

func (e WelcomeEmail) ToHTML() (string, error) {
    var buf bytes.Buffer
    if err := e.render().Render(context.Background(), &buf); err != nil {
        return "", err
    }
    return buf.String(), nil
}

func (e WelcomeEmail) ToText() (string, error) {
    html, err := e.ToHTML()
    if err != nil {
        return "", err
    }
    return HTMLToText(html)
}

templ (e WelcomeEmail) render() {
    @baseLayout("Subject", "Pre-header text") {
        // authored Tailwind markup
    }
}
```

Write utilities in `css/email.css`. `andurel run`, `andurel sync views`, `andurel sync email`, and `andurel build` compile templates automatically without changing the authored `.templ` files.

## Send typed messages

Build `email.TransactionalData`, then call `SendTransactional` with an injected sender:

```go
html, err := welcome.ToHTML()
if err != nil {
    return err
}
text, err := welcome.ToText()
if err != nil {
    return err
}

return email.SendTransactional(ctx, email.TransactionalData{
    To:       user.Email,
    From:     mailCfg.DefaultSenderSignature,
    Subject:  "Welcome",
    HTMLBody: html,
    TextBody: text,
}, sender)
```

Marketing messages require an unsubscribe URL. Validation, temporary, and permanent errors are distinguished so queue workers can choose retry behavior with `email.IsRetryable` / `email.IsValidationError`.

## Mailpit locally

```bash
andurel tool mailpit
```

Open the Mailpit UI (default `http://127.0.0.1:8025`) and confirm the message appears after the send call. Keep `EMAIL_PROVIDER=mailpit` in development `.env` so the Fx module publishes the Mailpit client as `TransactionalSender` / `MarketingSender`.

## Delivery from the queue

Use the queue for delivery that should not hold open an HTTP request. Commit domain state and a River job together, then send from `cmd/queue`. See [Queues](/docs/head/queues).

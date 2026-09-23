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

## Generate and compile email

```bash
andurel generate email Welcome
andurel sync email
```

Write email templates under `email/` and utilities in `css/email.css`. `andurel run`, `andurel sync views`, `andurel sync email`, and `andurel build` compile the templates automatically without changing the authored `.templ` files.

## Send typed messages

Build `email.TransactionalData` or `email.MarketingData`, then call the corresponding send helper. Marketing messages require an unsubscribe URL. Validation, temporary, and permanent errors are distinguished so queue workers can choose retry behavior correctly.

Local development defaults to Mailpit:

```bash
andurel tool mailpit
```

## Production providers

The `aws-ses` extension adds an SES client and provider configuration. Use the queue for delivery that should not hold open an HTTP request — commit domain state and a River job together, then send from `cmd/queue`. See [Queues](/docs/head/queues).

# Email

Andurel v2 provides a standalone email package, application-owned Templ messages, Tailwind-to-inline-style compilation, and separate transactional and marketing sender interfaces.

## Generate and compile email

```bash
andurel generate email Welcome
andurel email compile
```

Write email templates under `email/` and utilities in `css/email.css`. `andurel run`, `andurel generate view`, and `andurel build` compile the templates automatically without changing the authored `.templ` files.

## Send typed messages

The application selects a transport in `config.MailTransport` and Fx supplies `email.TransactionalSender` and `email.MarketingSender`. Local development defaults to Mailpit:

```bash
andurel tool mailpit
```

Build `email.TransactionalData` or `email.MarketingData`, then call the corresponding send helper. Marketing messages require an unsubscribe URL. Validation, temporary, and permanent errors are distinguished so queue workers can choose retry behavior correctly.

## Production providers

The `aws-ses` extension adds an SES client and provider configuration:

```bash
andurel extension add aws-ses --dry-run --diff --json
andurel extension add aws-ses
```

Use the queue for delivery that should not hold open an HTTP request. The web process inserts the job; the dedicated queue process resolves the configured sender and performs delivery.

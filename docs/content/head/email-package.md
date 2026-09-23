# Email Package

`github.com/mbvlabs/andurel/pkg/email` defines transport-neutral transactional and marketing messages, validates delivery contracts, classifies failures for retry policy, converts HTML to text, and provides a development Mailpit client. Templates and production adapters belong to the application.

## Sender interfaces

Consumers depend on capability-specific interfaces:

```go
type TransactionalSender interface {
    SendTransactional(context.Context, TransactionalPayload) error
}
type MarketingSender interface {
    SendMarketing(context.Context, MarketingPayload) error
}
```

A provider may implement one or both. The composition root selects `EMAIL_PROVIDER` and publishes these interfaces, keeping services and workers provider-neutral.

## Data and payload layers

`TransactionalData` and `MarketingData` are caller-facing values. `SendTransactional` and `SendMarketing` validate them, copy supported fields into payloads, and invoke the sender.

```go
err := email.SendTransactional(ctx, email.TransactionalData{
    To:       user.Email,
    From:     cfg.DefaultSenderSignature,
    Subject:  "Verify your account",
    HTMLBody: htmlBody,
    TextBody: textBody,
    Metadata: map[string]string{"user_id": user.ID.String()},
}, sender)
```

Transactional mail requires at least one To, Cc, or Bcc recipient, sender, subject, and HTML body. Marketing mail requires an unsubscribe URL and HTML body. Provider adapters may enforce more constraints.

## Templates and text alternatives

Generated applications compile email views with Templ, render them in application code, and pass strings to this package. `HTMLToText` extracts a fallback, omits head, script, and style content, preserves useful links, and normalizes whitespace.

Automatic conversion is a baseline. Review important text mail or author it explicitly. The package does not inline CSS, rewrite assets, or localize templates.

## Failure classification and retries

- `ValidationError` is permanent caller-input failure.
- `PermanentError` marks a provider rejection that must not retry.
- `TemporaryError` marks a known retryable failure.

`IsValidationError` recognizes typed and sentinel validation failures. `IsRetryable` returns true for temporary errors, false for permanent and validation errors, and true for unknown failures so mail is not silently lost.

Adapters should translate provider errors into these types. Do not retry malformed payloads indefinitely; retry timeouts, throttling, and transient service failures with bounded backoff.

## Mailpit transport

`NewMailpit(MailpitConfig{Host, Port})` validates its address and returns both sender capabilities. It uses SMTP and builds multipart alternative messages for HTML and optional text.

Mailpit is for local capture, not authenticated production delivery. Its context argument satisfies the interface, but the standard SMTP operation is not context-cancellable.

## Production adapters

Extensions such as `aws-ses` add application-owned clients and provider configuration. The adapter maps common payloads, classifies native errors, and is selected in the composition root. Keep provider-only features behind an application interface instead of leaking them through all domain services.

## Transactions and queues

Email is an external side effect and cannot join a PostgreSQL transaction. For important mail, commit domain state and a River job together, then render and send in an idempotent worker. Avoid sending before commit, because the recipient could learn about state that later rolls back. Keep large attachments out of job arguments and record provider identifiers where delivery auditing matters.

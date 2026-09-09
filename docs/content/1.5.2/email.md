# Email

Andurel includes Templ email components and interfaces for transactional and marketing delivery.

## Local delivery

Generated development projects use Mailpit. Start the local mail server with:

```bash
andurel tool mailpit
```

Mailpit accepts SMTP messages and provides a browser inbox for inspecting rendered email.

## Build an email

Generate a Templ email component:

```bash
andurel generate email Welcome
```

Keep delivery behind the generated sender interfaces so development and production providers can differ without changing application workflows.

## Production delivery

The optional `aws-ses` extension adds AWS SES configuration and a production sender:

```bash
andurel extension add aws-ses
```

Use background jobs for delivery that should not hold open an HTTP request.

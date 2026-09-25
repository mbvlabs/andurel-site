# email

`andurel generate email` authors a new email template under `email/`. Compiling Tailwind in emails and refreshing renderers is `andurel sync email`.

## When to use

- Create `email/NAME.templ` for transactional or marketing mail.
- Do not use this to recompile styles after editing a template; run [sync](/docs/head/sync) `email`.

## Usage

```bash
andurel generate email NAME [flags]
```

`NAME` is CamelCase (for example `WelcomeEmail` or `Receipt`).

## Flags

| Flag | What it does |
| --- | --- |
| `--dry-run` / `--diff` | Preview without writing; see [generate](/docs/head/generate) |

Email has no other command-local flags.

## Examples

```bash
andurel generate email WelcomeEmail --dry-run --json
andurel generate email WelcomeEmail
andurel generate email Receipt
```

This creates a template such as `email/welcome_email.templ`. Keep delivery behind the generated sender interfaces so development and production providers can differ without changing application workflows.

## Next

```bash
andurel sync email
```

See [Email](/docs/head/email) for compile, send, and provider boundaries.

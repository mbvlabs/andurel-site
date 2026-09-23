# Validation

Validate untrusted request data in controllers before passing it into services or models.

## Define rules

The generated `internal/validation` package provides common rules and helpers for structured validation errors.

Keep validation close to the HTTP boundary so every service receives values that satisfy the request contract.

## Return errors

For server-rendered forms, return field errors with the submitted values so the view can explain what needs correction. JSON APIs should use a stable error shape and an appropriate 4xx status.

## Protect trust boundaries

Validation does not replace authorization. Check both that input is well-formed and that the current actor is allowed to perform the operation.

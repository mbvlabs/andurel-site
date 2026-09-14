# Validation

`github.com/mbvlabs/andurel/pkg/validation` builds structured, JSON-ready field errors and records rule metadata during the same pass. Configuration and model inputs use it without coupling validation to Echo, Fx, Bun, Templ, or Inertia.

## Builder lifecycle

Create a builder, add rules, and return `Err()`. Every rule records its code, message, and parameters even when the value passes. A failed rule also appends a field error.

```go
func (d CreateProduct) Validate() error {
    b := validation.NewBuilder()
    b.Required("Name", d.Name)
    b.LenBetween("Name", d.Name, 2, 120)
    b.MinInt("PriceCents", d.PriceCents, 0)
    b.OneOf("Status", d.Status, "draft", "active", "archived")
    return b.Err()
}
```

`Err()` returns nil or `ValidationErrors`. `Errors()` exposes failures; `Rules()` returns a defensive copy of client-usable rule metadata. `ValidationErrors.ToMap()` keeps the first message for each field, matching common form contracts.

## Built-in rules

| Rule | Behavior |
| --- | --- |
| `Required`, `RequiredWhen` | Reject missing or zero-like values, conditionally when needed |
| `MinLen`, `MaxLen`, `LenBetween` | Count Unicode code points in present strings |
| `RecommendedLenBetween` | Publish recommendation metadata without creating an error |
| `MinInt`, `MaxInt` | Bound present integers |
| `OneOf`, `OneOfWithMessage` | Restrict a present string to allowed values |
| `URL`, `RequiredURL` | Accept HTTP or HTTPS URLs, optionally requiring presence |
| `MinItems`, `MaxItems`, `NoBlankItems` | Validate collection size and blank string elements |
| `TimeBeforeOrEqual` | Compare two present time values |
| `True` | Require agreement or confirmation |

Most format and range rules ignore absent optional values. Combine them with `Required` when absence is invalid. Messages may be customized; stable codes and parameters let clients localize or present richer guidance.

## Supported values and zero semantics

Presence helpers understand strings, integers, booleans, slices, `time.Time`, pointers, common `database/sql` nullable types, and defined aliases whose underlying kind is supported.

Be explicit about zero semantics. `Required` treats integer zero and zero time as missing, which is correct for many identifiers but wrong when zero is legitimate input. Model optionality with a pointer or nullable type and apply the domain's range rule. `NoBlankItems` reports indexed fields such as `Tags[2]`; an empty field name becomes `unknown`.

## Structured errors

Each error carries `field`, `code`, `message`, and optional `params`:

```json
{
  "field": "Name",
  "code": "min",
  "message": "must be at least 2 characters",
  "params": {"min": 2}
}
```

`validation.As(err)` extracts `ValidationErrors` through wrapped errors. `validation.Validate(value)` calls `Validate()` only when a value implements `Validatable`, which is useful at generic service boundaries.

## Domain-specific rules

Use `AddField` or `AddFieldWithParams` for failures and publish matching rule metadata when a client needs it:

```go
b.AddRuleWithParams("Slug", "reserved", "is reserved", map[string]any{
    "values": reservedSlugs,
})
if slices.Contains(reservedSlugs, data.Slug) {
    b.AddField("Slug", "reserved", "is reserved")
}
```

Keep rules deterministic and side-effect free. Database uniqueness, authorization, and race-sensitive invariants belong in services or persistence backed by constraints. Controllers may translate those failures into the same field-error shape.

## Templ and Inertia presentation

Templ controllers can pass `ValidationErrors` to a typed form component and return a 422 full response or fragment patch. Inertia controllers normally call `ToMap()` and pass it through `Page(...).ValidationErrors(...)`; named error bags follow the request protocol.

Domain validation stays identical across frontend choices. Only transport and presentation differ.

## Configuration validation

Generated constructors combine parsing errors with builder errors and semantic checks. The builder covers required values and bounds; constructors cover relationships such as valid ports, supported modes, or timeout ordering. The joined error prevents Fx from starting dependents.

Validation is not authorization. Authenticate the actor, authorize the operation and resource, validate input, then execute it. Never use field validation or hidden request context as a permission system.

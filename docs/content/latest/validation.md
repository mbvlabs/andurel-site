# Validation

Andurel v2 uses the standalone `github.com/mbvlabs/andurel/pkg/validation` package for composable rules and structured errors.

## Validate at each boundary

Generated configuration constructors validate parsed values before Fx starts dependent components. Generated model create and update types validate domain data before persistence. Controllers remain responsible for transport concerns and for presenting useful field errors.

```go
b := validation.NewBuilder()
b.Required("Name", data.Name)
b.MinInt("Price", data.Price, 0)
if err := b.Err(); err != nil {
    return err
}
```

Use exact rule names and signatures from the version of `pkg/validation` pinned in your application's `go.mod`.

## Keep authorization separate

Well-formed input is not necessarily allowed input. Authenticate the actor, authorize the resource and operation, then validate and execute the request. Services and models should receive `context.Context` and explicit dependencies, never a service locator hidden in request context.

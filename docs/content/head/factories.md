# Factories

Factories live under `models/factories` and build model entities for tests and seeds.

## Generate and sync

```bash
andurel sync factory User --check
andurel sync factory User --sync
andurel sync factories --check --json
andurel sync factories --sync
```

`sync factory` refreshes one Andurel-owned factory from the model Entity. `sync factories` covers every managed factory declaration. Prefer `--check` in CI.

## Seeds

Named seed compositions in root `seeds/` should call factories and exported model APIs, then run through `andurel db seed`. See [Migrations & Seeding](/docs/head/migrations).

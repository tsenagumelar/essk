# ESSK Backend

Go modular monolith backend for ESSK.

## Run

```text
go run ./cmd/server
```

## Migrations

Run from `services/backend`:

```text
go run ./cmd/essk migrate up
```

## Seed Admin

Run after migrations:

```text
go run ./cmd/essk seed admin
```

Default local credentials:

```text
admin@essk.local
Admin123!
```

## Policy Check

```text
go run ./cmd/essk policy check --all
```

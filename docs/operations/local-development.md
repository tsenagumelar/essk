# Local Development

## Requirements

- Go 1.25+.
- Node.js 24+.
- pnpm through Corepack.
- Docker.

## Backend

```text
cd services/backend
go run ./cmd/server
```

## Frontend

```text
pnpm --filter @essk/web dev
```

## Docker Compose

```text
pnpm dev:up
```

This command runs PostgreSQL, Redis, migrations, admin seed, backend, web, and Nginx.

Default admin:

```text
admin@essk.local
Admin123!
```

## Policy Check

```text
go -C services/backend run ./cmd/essk policy check --all
```

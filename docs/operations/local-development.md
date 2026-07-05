# Local Development

## Requirements

- Go 1.25+.
- Node.js 24+.
- pnpm through Corepack.
- Docker.

## Backend

```text
cd services/backend
go run ./services/api-gateway
```

## Frontend

```text
pnpm --filter @essk/web dev
```

## Backend in Docker, Web Local

Use this mode when changing the web UI and keeping PostgreSQL, Redis, migrations, seed, and API Gateway in Docker.

Terminal 1:

```text
pnpm dev:backend
```

Terminal 2:

```text
pnpm dev:web
```

This keeps the web service out of Docker, so port `3000` is owned only by the local Next.js dev server.

## Docker Compose

```text
pnpm dev:up
```

This command runs PostgreSQL, Redis, migrations, admin seed, backend, and web.

Default local URLs:

```text
http://localhost:3000
http://localhost:18080/health
```

Nginx is optional for local reverse-proxy testing and is disabled by default to avoid blocking startup when Docker cannot pull `nginx:1.27-alpine`:

```text
pnpm dev:edge
```

Default admin:

```text
admin@essk.local
Admin123!
```

## Policy Check

```text
go -C services/backend run ./cmd/essk policy check --all
```

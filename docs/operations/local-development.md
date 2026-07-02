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
docker compose -f infra/compose/docker-compose.yml up --build
```

## Policy Check

```text
go -C services/backend run ./cmd/essk policy check --all
```

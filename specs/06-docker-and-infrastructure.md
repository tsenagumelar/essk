# ESSK Specs - Docker And Infrastructure

## Local Infrastructure

Use Docker Compose as the default local runtime.

Required services:

- `backend`.
- `web`.
- `postgres`.
- `redis`.
- `nginx`.

Optional services:

- `admin`.
- `mailpit` for email testing.
- `minio` for file storage later.

## Compose Files

```text
infra/compose/
  docker-compose.yml
  docker-compose.override.yml
```

`docker-compose.yml` should be CI-safe and deterministic.

`docker-compose.override.yml` can contain local development conveniences:

- Bind mounts.
- Hot reload.
- Debug ports.

## Ports

Recommended local ports:

- Web: `3000`.
- Admin: `3001`.
- Backend: `8080`.
- PostgreSQL: `5432`.
- Redis: `6379`.
- Nginx: `80`.

## Backend Dockerfile

Use multi-stage build:

1. Builder stage with Go toolchain.
2. Runtime stage with distroless or Alpine.

Requirements:

- Non-root user.
- Static or minimal runtime binary.
- Healthcheck endpoint.
- Build metadata injected through ldflags.

## Frontend Dockerfile

Use Next.js standalone output.

Requirements:

- Multi-stage build.
- `pnpm` through Corepack.
- Non-root runtime.
- Environment variables handled correctly.

## Nginx

Nginx responsibilities:

- Route `/` to web.
- Route `/api` to backend.
- Route `/swagger` to backend.
- Apply basic security headers.
- Enforce upload/body limits.
- Future TLS termination.

## Environment Files

Root-level examples:

```text
.env.example
services/backend/.env.example
apps/web/.env.example
```

Rules:

- Commit example files.
- Never commit real secrets.
- Use strong generated secrets for local setup where possible.

## Healthchecks

Backend:

- `/health`: process is alive.
- `/ready`: database and Redis are reachable.

PostgreSQL:

- `pg_isready`.

Redis:

- `redis-cli ping`.

Web:

- HTTP GET `/`.

## Kubernetes Future

Reserve structure:

```text
infra/k8s/
  base/
  overlays/
    local/
    staging/
    production/
```

Initial Kubernetes specs can be introduced after Docker Compose milestone is stable.

Future required objects:

- Deployment.
- Service.
- Ingress.
- ConfigMap.
- Secret.
- Job for migrations.
- HorizontalPodAutoscaler.

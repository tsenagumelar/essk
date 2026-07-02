# Enterprise SaaS Starter Kit

ESSK is an opinionated, production-oriented starter kit for enterprise SaaS products.

The implementation follows the technical specifications in [`specs/`](specs/) and the task tracker in [`TASKS.md`](TASKS.md).

## Current Focus

Milestone 4: hardening and module scaffolding preparation.

## One-Command Local Startup

Use either command:

```text
pnpm dev:up
```

or:

```text
make dev
```

This starts:

- PostgreSQL.
- Redis.
- Database migrations.
- Admin seed.
- Backend API.
- Web app.
- Nginx reverse proxy.

Local URLs:

- Web through Nginx: `http://localhost`
- Web direct: `http://localhost:3000`
- Backend direct: `http://localhost:8080`
- Backend health: `http://localhost:8080/health`

Default local admin:

```text
admin@essk.local
Admin123!
```

Stop services:

```text
pnpm dev:down
```

Reset local database volume:

```text
pnpm dev:reset
```

## Verification

```text
make test
```

## Repository Layout

```text
apps/
services/
packages/
infra/
docs/
specs/
```

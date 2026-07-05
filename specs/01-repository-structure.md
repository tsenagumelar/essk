# ESSK Specs - Repository Structure

## Target Structure

```text
apps/
  web/
  admin/

services/
  backend/

packages/
  sdk-go/
  sdk-ts/
  shared/

infra/
  docker/
  compose/
  nginx/
  k8s/

docs/
specs/
```

## Directory Responsibilities

### `apps/web`

Primary product web app. In the starter phase this app contains:

- Login page.
- Authenticated app shell.
- Dashboard placeholder.
- User profile page.
- API health/status screen.

### `apps/admin`

Admin console. It can be bootstrapped later than `apps/web`, but repository structure reserves it for:

- Tenant administration.
- User management.
- Role and permission management.
- Audit log viewer.
- System configuration.

### `services/backend`

Main Go backend workspace. It starts as a modular monolith and evolves into separated services. It owns:

- Public HTTP API.
- Auth flows.
- Core platform modules.
- Database access.
- Redis access.
- Migrations.
- OpenAPI docs.
- Service-specific entrypoints under `services/backend/services/*`.
- Shared protobuf contracts under `services/backend/shared/protobuf`.
- Shared service schema migrations under `services/backend/shared/migrations`.

### `packages/sdk-ts`

Generated or hand-maintained TypeScript SDK for frontend apps.

Preferred source:

- Generate from OpenAPI.
- Use `openapi-typescript` for types.
- Use a thin fetch/axios wrapper for request execution.

### `packages/sdk-go`

Optional Go SDK for service-to-service usage, CLI tooling, or future integrations.

### `packages/shared`

Shared non-runtime assets:

- API constants.
- Documentation snippets.
- JSON schema.
- Shared TypeScript types only when generated types are not enough.

Avoid sharing business logic between frontend and backend through this package.

### `infra`

Infrastructure and deployment assets:

- Dockerfiles.
- Docker Compose files.
- Nginx config.
- Kubernetes manifests or Helm chart later.

### `docs`

Human-facing documentation:

- Getting started.
- Architecture decisions.
- Module creation guide.
- Local development guide.
- Deployment guide.

### `specs`

Planning and technical contracts before implementation. Specs are allowed to be more detailed and prescriptive than user documentation.

## Package Manager

Use `pnpm` for JavaScript/TypeScript workspaces.

Root files:

```text
pnpm-workspace.yaml
package.json
```

Recommended root scripts:

```json
{
  "scripts": {
    "dev": "docker compose -f infra/compose/docker-compose.yml up",
    "lint": "pnpm -r lint",
    "test": "pnpm -r test",
    "typecheck": "pnpm -r typecheck",
    "codegen": "pnpm --filter @essk/sdk-ts codegen"
  }
}
```

## Go Workspace

Backend can start as a standalone Go module:

```text
services/backend/go.mod
```

Use a root `go.work` only when multiple Go modules are introduced, such as `packages/sdk-go`.

## Naming Conventions

- Repository modules: kebab-case.
- Go packages: short lowercase names.
- Database tables: snake_case plural names.
- API paths: kebab-case or snake_case must be consistent; recommended kebab-case.
- Environment variables: uppercase snake_case with `ESSK_` prefix where useful.

## Versioning

- App version is injected at build time.
- Backend exposes version in `/health` and `/version`.
- Docker images use semantic version tags and commit SHA tags.

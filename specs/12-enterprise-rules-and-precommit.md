# ESSK Specs - Enterprise Rules And Pre-Commit

## Purpose

This spec defines enterprise-grade guardrails that must be enforced during development. These rules are not optional documentation; they become checks in pre-commit and CI.

## Enterprise Architecture Rules

### CQRS

All generated modules must separate commands and queries.

Command responsibilities:

- Create.
- Update.
- Delete or soft delete.
- Activate/deactivate.
- Assign/revoke relationships.
- Validate business invariants.
- Write audit logs.
- Invalidate affected caches.

Query responsibilities:

- List.
- Detail.
- Search.
- Filter.
- Sort.
- Paginate.
- Project response DTOs.

Rules:

- Query handlers must not mutate state.
- Command handlers must not return large list projections.
- Repository methods must make tenant and soft-delete filtering explicit.
- Generated module templates must include command and query files.

### Enterprise Data Rules

Every database table must include:

```sql
is_active boolean not null default true,
created_by uuid null,
created_date timestamptz not null default now(),
updated_by uuid null,
updated_date timestamptz not null default now(),
is_deleted boolean not null default false
```

Rules:

- No `created_at`, `updated_at`, or `deleted_at` columns.
- No default hard delete for business tables.
- All list/detail queries must filter `is_deleted = false`.
- Active resources should filter `is_active = true` unless inactive records are explicitly requested.
- Actor context must populate `created_by` and `updated_by`.

### Enterprise API Rules

Required:

- Standard response envelope.
- Standard error envelope.
- Request ID.
- Correlation ID.
- Auth middleware for protected routes.
- Permission middleware for business routes.
- Rate limiter middleware for public and sensitive routes.
- OpenAPI update for endpoint changes.

### Enterprise Frontend Rules

Required:

- Feature-based modular folder structure.
- Generated CRUD modules use route-level pages.
- Forms use React Hook Form and Zod.
- Server state uses TanStack Query.
- Loading, error, empty, and success states are implemented.
- Protected pages use authenticated layout.

## OWASP Top 10 Gate

Changed code must not weaken the OWASP control baseline documented in `09-security-and-observability.md`.

Pre-commit and CI must pay special attention to changes under:

```text
services/backend/internal/modules/auth/
services/backend/internal/modules/tenant/
services/backend/internal/modules/role/
services/backend/internal/modules/permission/
services/backend/internal/modules/audit/
services/backend/internal/middleware/
services/backend/internal/authn/
infra/
apps/web/lib/auth/
apps/web/lib/api/
```

Security-sensitive changes require tests or a documented exception.

## Pre-Commit Setup

Use `pre-commit`.

Config file:

```text
.pre-commit-config.yaml
```

Developer setup:

```text
pre-commit install
```

Manual full run:

```text
pre-commit run --all-files
```

## Required Pre-Commit Hooks

Generic:

- Check added large files.
- Check merge conflict markers.
- Check trailing whitespace.
- Check end-of-file.
- Validate YAML.
- Validate JSON.
- Validate TOML.
- Lint Markdown.

Security:

- Detect committed secrets.
- Scan dependency manifests.
- Block `.env` files except `.env.example`.
- Block hardcoded JWT secrets, private keys, and database passwords.

Backend:

- `gofmt`.
- `go vet`.
- `golangci-lint`.
- `govulncheck` when Go dependency files change.
- Validate migration names.
- Validate mandatory lifecycle columns.
- Validate no `created_at`, `updated_at`, or `deleted_at`.
- Validate no raw SQL string concatenation with request input.

Frontend:

- ESLint.
- TypeScript typecheck.
- Prevent secrets in `NEXT_PUBLIC_*`.
- Validate generated feature modules include schema, hooks, API wrapper, loading page, and error page.

Architecture:

- `essk policy check --staged`.

## `essk policy check`

The policy checker is a project-owned guardrail command.

Command:

```text
essk policy check --staged
essk policy check --all
```

Responsibilities:

- Inspect changed file paths.
- Apply rule groups based on changed area.
- Print actionable failure messages.
- Exit non-zero when rules are violated.

Rule groups:

- `database`.
- `backend-module`.
- `frontend-module`.
- `security`.
- `openapi`.
- `scaffold-template`.

Example failures:

```text
database: migration 000010_create_orders.up.sql creates table orders without updated_by
backend-module: services/backend/internal/modules/order/route.go registers POST /orders without permission middleware
frontend-module: apps/web/app/(app)/orders/page.tsx is missing error.tsx
security: auth middleware changed but no related tests changed
openapi: handler route changed but swagger docs were not updated
```

## CI Enforcement

CI must run the same checks as pre-commit:

```text
pre-commit run --all-files
essk policy check --all
```

CI is the source of truth. Pre-commit is a fast local guardrail.

## Exceptions

Exceptions are allowed only with explicit documentation.

Required exception format in PR summary:

```text
Policy exception:
- Rule:
- Reason:
- Risk:
- Compensating control:
- Follow-up:
```

Examples:

- Temporary migration exception during bootstrap.
- Public unauthenticated route for healthcheck.
- System table without foreign key to `users` because of bootstrap ordering.

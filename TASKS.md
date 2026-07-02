# ESSK Implementation Tasks

## Milestone 1 - Foundation

Status legend:

- `[ ]` Not started.
- `[~]` In progress.
- `[x]` Done.

### Task 1 - Monorepo Bootstrap

- `[x]` Create root workspace metadata.
- `[x]` Create repository directory structure.
- `[x]` Add root development scripts.
- `[x]` Add root environment example.
- `[x]` Add ignore rules.

### Task 2 - Backend Foundation

- `[x]` Create Go module under `services/backend`.
- `[x]` Add Fiber server bootstrap.
- `[x]` Add config loader.
- `[x]` Add structured logger.
- `[x]` Add standard API response helper.
- `[x]` Add global error handler.
- `[x]` Add validation helper.
- `[x]` Add health, readiness, and version endpoints.
- `[x]` Add PostgreSQL and Redis connection wiring.

### Task 3 - Frontend Foundation

- `[x]` Create Next.js app under `apps/web`.
- `[x]` Add Tailwind setup.
- `[x]` Add authenticated layout shell.
- `[x]` Add login page placeholder.
- `[x]` Add dashboard page.
- `[x]` Add API health status page.
- `[x]` Add API client foundation.

### Task 4 - Infrastructure Foundation

- `[x]` Add backend Dockerfile.
- `[x]` Add web Dockerfile.
- `[x]` Add Docker Compose.
- `[x]` Add PostgreSQL and Redis services.
- `[x]` Add Nginx reverse proxy.
- `[x]` Add healthchecks.

### Task 5 - Quality And Guardrails

- `[x]` Add pre-commit config.
- `[x]` Add `essk policy check` skeleton.
- `[x]` Add basic GitHub Actions CI.
- `[x]` Add formatting and lint commands.

### Task 6 - Verification

- `[x]` Run backend tests.
- `[x]` Run frontend typecheck/lint when dependencies are installed.
- `[x]` Validate Docker Compose config.
- `[x]` Commit and push completed foundation slice.

## Milestone 2 - Auth Core

- `[x]` Add core auth migrations.
- `[x]` Add password hashing.
- `[x]` Add login endpoint.
- `[x]` Add refresh token rotation.
- `[x]` Add logout endpoint.
- `[x]` Add auth middleware.
- `[x]` Add `/auth/me`.
- `[x]` Add seeded admin user.
- `[x]` Connect frontend login flow.

## Milestone 3 - Tenant And RBAC

- `[x]` Add tenant module.
- `[x]` Add role module.
- `[x]` Add permission module.
- `[x]` Add user-role assignment.
- `[x]` Add role-permission assignment.
- `[x]` Add permission middleware.
- `[x]` Add tenant boundary tests.

## Milestone 4 - Audit And Hardening

- `[x]` Add audit log module.
- `[ ]` Add Redis-backed rate limiter.
- `[ ]` Add OWASP-sensitive tests.
- `[ ]` Add security headers.
- `[ ]` Add observability hooks.

## Milestone 5 - Module Scaffolding

- `[ ]` Add `essk add-module` command.
- `[ ]` Add backend scaffolding templates.
- `[ ]` Add frontend scaffolding templates.
- `[ ]` Add migration/query scaffolding.
- `[ ]` Add scaffold policy tests.

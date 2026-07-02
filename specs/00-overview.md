# ESSK Specs - Overview

## Purpose

Enterprise SaaS Starter Kit (ESSK) is a reusable code foundation for building enterprise SaaS products consistently. The primary goal is not to build a finished product, but to provide an opinionated starter kit that can be reused for Approval Workflow, HR System, ERP Modules, Inventory, Audit Platform, Customer Portal, and other multi-tenant products.

## Product Type

- Type: opinionated SaaS starter kit.
- Initial mode: modular monolith.
- Evolution path: microservices-ready when domain boundaries and traffic justify service separation.
- Initial deployment: Docker Compose.
- Future deployment: Kubernetes.

## Core Principles

- Go first for backend.
- Clean Architecture.
- API first.
- Security by default.
- Docker first.
- Testable by design.
- Simple before complex.
- Consistent module template.
- Documentation updated with every feature.

## Phase Target

### Phase 1 - Foundation

Goal: a new developer can clone the repository and run the system with one command.

Deliverables:

- Monorepo bootstrap.
- Backend service.
- Web frontend.
- Docker Compose.
- PostgreSQL.
- Redis.
- Environment configuration.
- Health endpoint.
- Structured logging.
- Standard API response.
- Global error handler.
- Request validation.
- Swagger/OpenAPI.
- Migration system.

### Phase 2 - Core Platform

Goal: provide the minimum platform modules required for enterprise SaaS products.

Deliverables:

- Authentication.
- User management.
- Tenant management.
- RBAC.
- Audit log.
- Protected API.
- Frontend login and authenticated shell.

## Architecture Summary

Request flow:

```text
HTTP Request
  -> Router
  -> Middleware
  -> Handler
  -> Service
  -> Repository
  -> Database/Cache
```

Module shape:

```text
module/
  entity.go
  dto.go
  repository.go
  service.go
  handler.go
  route.go
  validation.go
  errors.go
  *_test.go
```

## Main Technology Decisions

Backend:

- Language: Go 1.23+.
- HTTP framework: Fiber v2.
- Database: PostgreSQL 16+.
- Cache/session store: Redis 7+.
- Query layer: sqlc + pgx.
- Migration: golang-migrate.
- Architecture guardrail: CQRS for module application layer.
- Security baseline: OWASP Top 10 controls.
- Local policy gate: pre-commit hooks.
- Validation: go-playground/validator.
- Logging: zerolog.
- Configuration: koanf.
- OpenAPI: swaggo/swag for early phase, OpenAPI file generation can be improved later.
- JWT: golang-jwt/jwt/v5.
- Password hashing: argon2id.

Frontend:

- Framework: Next.js 15+ App Router.
- Runtime: React 19+.
- Language: TypeScript.
- Styling: Tailwind CSS.
- UI primitives: Radix UI.
- Component baseline: shadcn/ui style structure.
- Server/client state: TanStack Query.
- Lightweight client state: Zustand.
- Forms: React Hook Form.
- Validation: Zod.
- HTTP client: generated TypeScript SDK where possible, fallback axios.

Infrastructure:

- Local orchestration: Docker Compose.
- Reverse proxy: Nginx.
- CI/CD: GitHub Actions.
- Container build: Docker BuildKit.
- Future orchestration: Kubernetes.

## Non-Goals For Foundation

- No domain-specific product module.
- No premature microservice split.
- No complex workflow engine in phase 1.
- No multi-cloud abstraction.
- No event streaming unless a module requires it later.
- No custom UI design system beyond a stable starter structure.

## Definition of Done

Milestone 1 is done when a developer can:

- Clone repository.
- Run `docker compose up`.
- Open frontend.
- Call backend API.
- Login with seeded admin account.
- Access protected endpoint.
- View Swagger/OpenAPI docs.
- Create and run a migration.
- Create a new backend module from the documented template.
- Run unit and integration tests.

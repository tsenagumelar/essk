# ESSK Specs - Backend

## Stack

- Go 1.23+.
- Fiber v2 for HTTP routing.
- PostgreSQL 16+.
- Redis 7+.
- sqlc for type-safe SQL.
- pgx/v5 as PostgreSQL driver.
- golang-migrate for schema migrations.
- go-playground/validator for request validation.
- zerolog for structured logging.
- koanf for configuration loading.
- swaggo/swag for Swagger generation.
- golang-jwt/jwt/v5 for JWT.
- argon2id for password hashing.
- testify for tests.
- testcontainers-go for integration tests.

## Backend Structure

```text
services/backend/
  cmd/
    server/
      main.go

  internal/
    app/
    config/
    database/
    cache/
    middleware/
    logger/
    response/
    validator/
    errors/
    authn/

    modules/
      auth/
      user/
      tenant/
      role/
      permission/
      audit/

  migrations/
  queries/
  docs/
  tests/
```

## Layering Rules

Handler:

- Accepts HTTP request.
- Parses params/body.
- Calls validator.
- Calls service.
- Converts result to standard API response.
- Must not contain business logic.

Service:

- Owns business rules.
- Controls transaction boundary when use case spans multiple repositories.
- Calls repositories and external clients.
- Emits audit events.
- Must not depend on Fiber.

Repository:

- Owns database queries.
- Uses sqlc generated query methods.
- Accepts `context.Context`.
- Must not know HTTP or request DTO.

Entity:

- Represents core domain data.
- Must be independent from transport concerns.

DTO:

- Request and response shapes.
- Versioned by API contract.

## CQRS Rules

ESSK uses pragmatic CQRS inside each module.

Rules:

- Command operations mutate state: create, update, delete, activate, deactivate, assign, revoke.
- Query operations read state and must not mutate data.
- Command handlers/services and query handlers/services are separated for non-trivial modules.
- Simple foundation modules may keep command and query methods in one service file early, but generated modules must already follow CQRS file naming so the pattern is stable.
- Commands run validation, authorization, transaction handling, domain rules, audit logging, and cache invalidation.
- Queries handle filtering, pagination, sorting, projection, and tenant scoping.
- Query DTOs must not expose internal database models directly.
- Read-side optimization can use dedicated SQL queries or views later without changing command APIs.

Recommended module shape with CQRS:

```text
internal/modules/{module}/
  entity.go
  dto.go
  commands.go
  queries.go
  command_handler.go
  query_handler.go
  repository.go
  service.go
  handler.go
  route.go
  validation.go
  errors.go
  audit.go
  command_handler_test.go
  query_handler_test.go
  handler_test.go
```

## Module Template

```text
internal/modules/{module}/
  entity.go
  dto.go
  commands.go
  queries.go
  command_handler.go
  query_handler.go
  repository.go
  service.go
  handler.go
  route.go
  validation.go
  errors.go
  audit.go
  service_test.go
  handler_test.go
```

Required for every feature:

- Migration.
- Query file.
- Repository.
- Service.
- Handler.
- Route.
- Validation.
- Tests where applicable.
- API documentation update.

## Module Scaffolding

Backend must provide a module scaffolding command so new product modules follow the same structure.

Command:

```text
essk add-module user-management
```

Backend output for `user-management`:

```text
services/backend/internal/modules/user_management/
  entity.go
  dto.go
  repository.go
  service.go
  handler.go
  route.go
  validation.go
  errors.go
  audit.go
  service_test.go
  handler_test.go

services/backend/queries/user_management.sql
services/backend/migrations/{version}_create_user_management.up.sql
services/backend/migrations/{version}_create_user_management.down.sql
```

Generated backend CRUD endpoints:

```text
GET    /api/v1/user-management
POST   /api/v1/user-management
GET    /api/v1/user-management/:id
PATCH  /api/v1/user-management/:id
DELETE /api/v1/user-management/:id
```

Generated table requirements:

- `id uuid primary key`.
- Optional `tenant_id uuid null` when `--tenant-scoped=true`.
- Default display field such as `name varchar(160) not null`.
- Mandatory lifecycle columns: `is_active`, `created_by`, `created_date`, `updated_by`, `updated_date`, `is_deleted`.

Recommended implementation:

- CLI package lives in `services/backend/cmd/essk`.
- Templates live in `services/backend/internal/scaffold/templates`.
- Use Go `text/template`.
- Convert module name into:
  - kebab-case route path: `user-management`.
  - snake_case package/table name: `user_management`.
  - PascalCase type name: `UserManagement`.
  - camelCase variable name: `userManagement`.

Safety rules:

- Command must fail if target files already exist unless `--force` is passed.
- Command must not modify unrelated modules.
- Command must register route in a central module registry or generate a clear TODO if automatic registry update is not implemented yet.
- Command must create compilable code, even if generated CRUD business fields are minimal.
- Command must generate CQRS command/query files by default.

## Configuration

Config sources, in priority order:

1. Environment variables.
2. `.env.local`.
3. `.env`.
4. Defaults in code.

Recommended library: `github.com/knadh/koanf`.

Required config groups:

- App: name, env, version, port.
- HTTP: read timeout, write timeout, body limit.
- Database: URL, max open connections, max idle connections.
- Redis: address, password, database.
- Auth: access token TTL, refresh token TTL, JWT issuer, signing key.
- CORS: allowed origins, methods, headers.
- Rate limit: enabled flag, store, per-route limits, auth limits.
- Logging: level, pretty mode.

## HTTP API

Base path:

```text
/api/v1
```

System endpoints:

```text
GET /health
GET /ready
GET /version
GET /swagger/*
```

## Middleware

Default middleware order:

1. Recover.
2. Request ID.
3. Logger.
4. CORS.
5. Security headers.
6. Body limit.
7. Rate limiting when enabled.
8. Authentication for protected routes.
9. Authorization for permission-protected routes.

Recommended libraries:

- Fiber recover middleware.
- Fiber requestid middleware or custom request ID.
- Fiber cors middleware.
- Fiber helmet middleware.
- Fiber limiter middleware with Redis-backed store for distributed rate limiting.

## Standard Response

Success:

```json
{
  "success": true,
  "message": "OK",
  "data": {},
  "meta": {}
}
```

Error:

```json
{
  "success": false,
  "message": "Validation Error",
  "errors": []
}
```

Meta for paginated lists:

```json
{
  "page": 1,
  "page_size": 20,
  "total_items": 100,
  "total_pages": 5
}
```

## Error Handling

Use centralized application errors:

- `BAD_REQUEST`.
- `UNAUTHORIZED`.
- `FORBIDDEN`.
- `NOT_FOUND`.
- `CONFLICT`.
- `VALIDATION_ERROR`.
- `INTERNAL_ERROR`.

Each error response must include:

- Stable error code.
- Human-readable message.
- Validation details if applicable.
- Request ID through response header.

## Logging

Use JSON logs by default.

Required log fields:

- `timestamp`.
- `level`.
- `message`.
- `request_id`.
- `correlation_id`.
- `method`.
- `path`.
- `status`.
- `latency_ms`.
- `user_id` when authenticated.
- `tenant_id` when available.

Pretty logs are allowed only in local development.

## Background Jobs

Do not add a queue in phase 1.

Future options:

- Asynq with Redis for simple background jobs.
- Temporal for complex workflows.

Decision rule:

- Use Asynq for emails, notifications, and retries.
- Use Temporal only for long-running business workflows.

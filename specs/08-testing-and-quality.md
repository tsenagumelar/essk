# ESSK Specs - Testing And Quality

## Testing Pyramid

Backend:

- Unit tests for services and validators.
- Handler tests for API behavior.
- Repository integration tests with PostgreSQL.
- End-to-end API smoke tests.

Frontend:

- Unit tests for utilities and API wrappers.
- Component tests for forms and UI states.
- E2E tests for critical flows.

Infrastructure:

- Docker Compose validation.
- Container startup smoke test.

Operational readiness:

- Backup and restore test.
- Migration rollback test.
- Basic load test for critical API paths.

## Backend Test Tools

- `testing` standard package.
- `testify`.
- `httptest` or Fiber app test helpers.
- `testcontainers-go`.
- `go-sqlmock` only for narrow cases; prefer real PostgreSQL for repositories.

Required coverage:

- Auth login success/failure.
- Password hashing verification.
- Refresh token rotation.
- Tenant boundary enforcement.
- RBAC permission checks.
- Standard API error format.
- Audit log creation for critical actions.
- Rate limiter behavior.
- CQRS query methods do not mutate state.
- OWASP-sensitive route tests.

## Frontend Test Tools

- Vitest.
- React Testing Library.
- MSW for API mocking.
- Playwright for E2E.

Required coverage:

- Login form validation.
- Login success redirect.
- Login failure error display.
- Protected route behavior.
- API error normalization.
- Logout flow.
- CRUD scaffolding output renders loading, error, empty, and success states.

## Code Quality

Backend:

- `gofmt`.
- `go vet`.
- `golangci-lint`.
- No unchecked errors unless intentionally documented.

Frontend:

- ESLint.
- TypeScript strict mode.
- Prettier.
- Tailwind class sorting optional.

## Pre-Commit Quality Gate

Use `pre-commit` to enforce local rules before code enters commits.

Config file:

```text
.pre-commit-config.yaml
```

Install command:

```text
pre-commit install
```

Required checks:

- Markdown formatting/lint for `specs` and `docs`.
- Secret scanning.
- Large file prevention.
- End-of-file and trailing whitespace checks.
- YAML/JSON/TOML validation.
- Go formatting through `gofmt`.
- Go lint through `golangci-lint` when backend files change.
- Go vulnerability check through `govulncheck` for dependency-sensitive changes.
- TypeScript/JavaScript lint when frontend files change.
- TypeScript typecheck when frontend source or generated SDK changes.
- Migration naming validation.
- Mandatory database lifecycle columns validation.
- Module architecture validation.
- OWASP/security rule validation for changed auth, middleware, tenant, RBAC, and audit files.

Recommended tools:

- `pre-commit`.
- `detect-secrets` or `gitleaks`.
- `golangci-lint`.
- `govulncheck`.
- `pnpm lint`.
- `pnpm typecheck`.
- Custom `essk policy check` command.

## Architecture Policy Checks

The starter kit must include a policy checker that can run locally and in CI.

Command:

```text
essk policy check --staged
```

The checker validates staged changes before commit.

Required policy checks:

- New migration files must include mandatory lifecycle columns for every created table.
- New soft-delete capable tables must use `is_deleted`, not physical delete by default.
- New backend module must include route, handler, command/query files, repository, validation, and tests.
- New module routes must use permission middleware unless explicitly marked public.
- New tenant-scoped repositories must filter by tenant ID and `is_deleted = false`.
- New read/list SQL must default to `is_deleted = false`.
- New write operations must set `updated_by` and `updated_date`.
- New create operations must set `created_by`, `created_date`, `updated_by`, and `updated_date`.
- Auth, RBAC, tenant, and audit changes must include tests.
- Frontend module pages must include loading and error states.
- Frontend forms must use Zod and React Hook Form.
- Public API changes must update OpenAPI documentation.
- Scaffolding templates must stay aligned with current module rules.
- Security-sensitive changes must include tests or a documented policy exception.
- Migration changes must include rollback notes.

Policy violations should fail pre-commit with actionable messages.

## Commit Quality

Recommended convention:

```text
feat: add auth login endpoint
fix: enforce tenant boundary on users
docs: add module creation guide
chore: update ci cache
```

## Pull Request Expectations

Each PR should include:

- Summary.
- Architectural decisions when relevant.
- Test evidence.
- Migration notes if schema changes.
- API changes if endpoint contract changes.

## Definition of Ready For Feature Work

A new module can be started when:

- Backend bootstrap runs.
- Database migration command works.
- Standard response helper exists.
- Global error handler exists.
- Validator helper exists.
- Module template is documented.
- CI runs basic checks.

## Definition of Done For A Feature

Feature is complete when:

- Migration exists when needed.
- Queries generated.
- Repository implemented.
- Service implemented.
- Handler and route registered.
- Request validation implemented.
- Tests added.
- OpenAPI updated.
- Docs updated.
- CI passes.

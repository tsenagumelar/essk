# Phase 01 - Foundation Stabilization

## Status

Current status: In Progress

## Goal

Make the current starter kit clean, consistent, testable, and easy to extend before adding heavier enterprise infrastructure.

## Current State

- Backend already has modular packages for auth, tenants, RBAC, users, products, and audit.
- Frontend has started moving toward `app`, `features`, and `shared`.
- Local Docker Compose exists for backend, database, Redis, migration, and seed.
- Seed data exists for tenants, roles, and users.
- Product CRUD exists as a sample module.
- Some frontend admin views still share temporary implementation instead of complete feature isolation.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| Frontend structure | Some feature logic is still shared through temporary admin components | Harder to maintain and scaffold new modules |
| Backend structure | Modules are modular but command/query boundaries are not consistently enforced | CQRS remains mostly documented rather than implemented |
| API contract | OpenAPI generation and drift validation are not fully enforced | Frontend/backend contract can drift |
| Testing | Integration coverage is still limited | Tenant and RBAC regressions can slip through |
| Error handling | API response shape needs strict consistency | Frontend handling becomes inconsistent |
| Module documentation | Module-level README template is missing | New contributors may implement modules inconsistently |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 01-01 | Finalize frontend feature isolation | In Progress | Move tenant, user, role, product, dashboard, profile, and auth logic into their own feature folders |
| 01-02 | Create shared atomic design structure | Not Started | Add `shared/atoms`, `shared/molecules`, `shared/organisms`, `shared/templates`, `shared/hooks`, `shared/utils` |
| 01-03 | Remove temporary cross-feature admin implementation | Not Started | Replace temporary shared admin workspaces with feature-specific implementations and shared reusable components |
| 01-04 | Standardize backend module shape | In Progress | Enforce handler, command service, query service, repository, DTO, validation, route, and tests |
| 01-05 | Implement real command/query split in backend modules | Not Started | Separate write use cases from read use cases at application layer |
| 01-06 | Standardize API error response | Not Started | Define and enforce one error response format across all handlers |
| 01-07 | Standardize list API behavior | In Progress | Pagination, search, filter, sort, and export conventions must be consistent |
| 01-08 | Add OpenAPI generation | Not Started | Generate OpenAPI from backend routes or maintained contract file |
| 01-09 | Add OpenAPI drift check in CI | Not Started | Fail CI when API changes without contract update |
| 01-10 | Add integration tests for auth and RBAC | Not Started | Cover login, me, refresh, logout, permission checks |
| 01-11 | Add integration tests for tenant isolation | Not Started | Verify admin/user cannot access other tenant data |
| 01-12 | Add module README template | Not Started | Document module purpose, routes, permissions, events, tables, and tests |

## Exit Criteria

- Existing frontend routes are thin route files only.
- Feature-specific logic lives under `features/{feature}`.
- Shared UI and utilities are reusable without feature coupling.
- Backend modules follow one consistent structure.
- Command and query boundaries are visible in backend application code.
- API errors and list responses are standardized.
- CI validates lint, typecheck, tests, build, and API contract.

## Dependencies

- Existing backend and frontend modules.
- Existing specs for frontend, backend, testing, and scaffolding.

## Notes

This phase should be completed before the platform adds event bus, workers, plugin registry, or Kubernetes complexity.

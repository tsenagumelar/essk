# Phase 07 - Module Scaffolding And Plugin Model

## Status

Current status: Not Started

## Goal

Make new modules fast, consistent, secure, and easy to maintain.

## Current State

- Specs describe module scaffolding.
- Product module exists as a sample CRUD module.
- No working `essk add-module` command exists yet.
- Plugin registry does not exist yet.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| CLI scaffolder | Missing | New modules require manual boilerplate |
| Backend templates | Missing | Module quality depends on developer memory |
| Frontend templates | Missing | UI consistency can drift |
| Permission generation | Missing | RBAC can be forgotten for new modules |
| Audit event generation | Missing | Compliance logging can be incomplete |
| Plugin model | Missing | Modules cannot be enabled per tenant |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 07-01 | Define scaffolder CLI contract | Not Started | `essk add-module {module-name}` with options |
| 07-02 | Add backend module templates | Not Started | Entity, DTO, command, query, repository, handler, routes, tests |
| 07-03 | Add migration template | Not Started | Include required audit fields and soft delete |
| 07-04 | Add SQL/query template | Not Started | Standard list, get, create, update, delete queries |
| 07-05 | Add frontend feature template | Not Started | `index.tsx`, `types.ts`, `hooks.ts`, route file |
| 07-06 | Add CRUD UI template | Not Started | Search, filter, export, table, pagination, popup form |
| 07-07 | Add permission template | Not Started | `{module}:read`, `{module}:create`, `{module}:update`, `{module}:delete` |
| 07-08 | Add audit event template | Not Started | Created, updated, deleted, exported |
| 07-09 | Add menu registration template | Not Started | Optional sidebar menu registration |
| 07-10 | Add generated test template | Not Started | Backend unit/integration and frontend basic render tests |
| 07-11 | Add plugin manifest template | Not Started | Name, code, version, permissions, routes, events |
| 07-12 | Validate generated module compiles | Not Started | Scaffolder must run lint, typecheck, and tests |

## Exit Criteria

- Running `essk add-module user-management` generates backend, frontend, database, permission, audit, menu, and test artifacts.
- Generated code compiles without manual boilerplate fixes.
- Generated modules follow the same enterprise rules as hand-written modules.

## Dependencies

- Phase 01 stable module structure.
- Phase 02 permission and audit policy.
- Phase 06 plugin service for full tenant-level plugin enablement.

## Notes

This phase should start after module structure is stable. Otherwise the scaffolder will encode unstable patterns.

# Phase 02 - Enterprise Multi-Tenancy And Security

## Status

Current status: Partially Done

## Goal

Enforce tenant isolation and enterprise security at application, database, gateway, and operational layers.

## Current State

- Tenant, user, and role management exist.
- Default roles exist: `super_admin`, `admin`, and `user`.
- Super admin can manage global tenant data.
- Tenant admin and tenant user behavior exists at application level for several screens and APIs.
- Redis-backed rate limiting exists as a foundation.
- Audit fields and soft delete conventions exist.
- OWASP and enterprise rule specs exist.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| Tenant enforcement | Application-level scoping exists but database-level RLS is missing | Bugs can bypass tenant isolation if query code is wrong |
| Tenant context | Tenant resolution needs a formal middleware and propagation contract | Services may implement tenant checks inconsistently |
| API keys | No tenant-aware API key model | Third-party integrations cannot be safely scoped |
| SSO/MFA | Not implemented | Enterprise identity requirements are not met |
| Rate limiting | Current foundation needs per-tenant and per-user policy | No differentiated limits by tenant, user, endpoint, or API key |
| Audit security events | Audit exists but event coverage needs policy | Sensitive actions may not be fully traceable |
| Secrets policy | Production secrets governance needs implementation | Risk of insecure config or leaked secrets |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 02-01 | Define tenant context contract | Not Started | Standardize tenant id, actor id, role, permissions, correlation id |
| 02-02 | Add tenant resolver middleware | Not Started | Resolve tenant from JWT, API key, route, or trusted gateway header |
| 02-03 | Enforce tenant context in repositories | In Progress | Ensure tenant-scoped queries always use tenant id |
| 02-04 | Add PostgreSQL RLS templates | Not Started | Create policy pattern for all tenant-scoped tables |
| 02-05 | Enable RLS for tenant-scoped core tables | Not Started | Apply to users, roles, products, and future tenant data |
| 02-06 | Document permission matrix | Not Started | Define role permissions for super admin, admin, and user |
| 02-07 | Add API key database model | Not Started | Tenant-scoped API keys with hash storage, expiry, status, and audit fields |
| 02-08 | Add API key auth middleware | Not Started | Authenticate third-party systems without user JWT |
| 02-09 | Add per-tenant rate limit policy | Not Started | Configure limits by tenant and endpoint group |
| 02-10 | Add MFA-ready auth model | Not Started | Add tables and flow design for TOTP or email challenge |
| 02-11 | Add SSO-ready auth model | Not Started | Add identity provider mapping and SAML/OIDC design |
| 02-12 | Expand security audit events | Not Started | Login, logout, failed login, role assignment, permission change, API key change |
| 02-13 | Add production secrets policy | Not Started | Document and enforce required secrets for production |
| 02-14 | Add OWASP verification checklist to CI/pre-commit where practical | Not Started | Validate dependency, secret, lint, security, and config checks |

## Exit Criteria

- Tenant isolation is enforced by middleware, service/repository code, and PostgreSQL RLS.
- Super admin can access all tenants; tenant admin and user cannot cross tenant boundaries.
- API keys are tenant-scoped and auditable.
- Security-sensitive actions create audit records.
- Rate limits can be applied by tenant, user, endpoint, and API key.
- SSO and MFA have clear implementation paths even if providers are added later.

## Dependencies

- Phase 01 API and module consistency.
- PostgreSQL migration discipline.
- Audit module foundation.

## Notes

This is a security-critical phase. It should be prioritized before external integrations and public API exposure.
